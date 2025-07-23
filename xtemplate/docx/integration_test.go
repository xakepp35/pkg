package docx

import (
	"archive/zip"
	"bytes"
	"io"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestData представляет тестовые данные для шаблона
type TestData struct {
	Title     string
	Footer    string
	Items     []TestItem
	ImageData []byte
}

type TestItem struct {
	ID     int64
	Name   string
	Email  string
	Number float64
}

// createMinimalDocx создает минимальный DOCX для тестирования
func createMinimalDocx() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	// Content_Types.xml
	contentTypes := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
    <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
    <Default Extension="xml" ContentType="application/xml"/>
    <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`
	writeZipFile(zw, "[Content_Types].xml", contentTypes)

	// _rels/.rels
	rels := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
    <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`
	writeZipFile(zw, "_rels/.rels", rels)

	// word/document.xml с тестовыми плейсхолдерами
	document := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
    <w:body>
        <w:p>
            <w:r><w:t>Title: {{.Ti</w:t></w:r>
            <w:r><w:t>tle}}</w:t></w:r>
        </w:p>
        <w:p>
            <w:r><w:t>{{addImage .ImageData 200 200}}</w:t></w:r>
        </w:p>
        <w:p>
            <w:r><w:t>Footer: {{.Footer}}</w:t></w:r>
        </w:p>
    </w:body>
</w:document>`
	writeZipFile(zw, "word/document.xml", document)

	// word/_rels/document.xml.rels
	docRels := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
</Relationships>`
	writeZipFile(zw, "word/_rels/document.xml.rels", docRels)

	zw.Close()
	return buf.Bytes()
}

func writeZipFile(zw *zip.Writer, name, content string) {
	w, _ := zw.Create(name)
	w.Write([]byte(content))
}

func TestTemplate_Integration_BasicTemplate(t *testing.T) {
	// Создаем минимальный DOCX
	docxData := createMinimalDocx()

	// Создаем шаблон
	tmpl := New("test")
	tmpl.ParseDocxFileData(docxData)
	require.NoError(t, tmpl.err)

	// Подготавливаем данные
	data := TestData{
		Title:  "Test Title",
		Footer: "Test Footer",
	}

	// Выполняем шаблон
	var result bytes.Buffer
	err := tmpl.Execute(&result, data)
	require.NoError(t, err)

	// Проверяем результат
	resultData := result.Bytes()
	assert.True(t, len(resultData) > 0)

	// Извлекаем и проверяем document.xml
	docXML := extractDocumentXML(t, resultData)
	assert.Contains(t, docXML, "Title: Test Title")
	assert.Contains(t, docXML, "Footer: Test Footer")
}

func TestTemplate_Integration_BrokenPlaceholders(t *testing.T) {
	docxData := createMinimalDocx()

	tmpl := New("test")
	tmpl.ParseDocxFileData(docxData)
	require.NoError(t, tmpl.err)

	data := TestData{
		Title:  "Fixed Title",
		Footer: "Fixed Footer",
	}

	var result bytes.Buffer
	err := tmpl.Execute(&result, data)
	require.NoError(t, err)

	// Проверяем, что разбитый плейсхолдер {{.Ti + tle}} был склеен
	docXML := extractDocumentXML(t, result.Bytes())
	assert.Contains(t, docXML, "Fixed Title")
	assert.NotContains(t, docXML, "{{.Ti")
	assert.NotContains(t, docXML, "tle}}")
}

func TestTemplate_Integration_WithImages(t *testing.T) {
	docxData := createMinimalDocx()

	tmpl := New("test")
	tmpl.ParseDocxFileData(docxData)
	require.NoError(t, tmpl.err)

	// Тестовое изображение (PNG заголовок)
	imageData := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00}

	data := TestData{
		Title:     "With Image",
		Footer:    "Image Test",
		ImageData: imageData,
	}

	var result bytes.Buffer
	err := tmpl.Execute(&result, data)
	require.NoError(t, err)

	// Проверяем структуру результата
	resultData := result.Bytes()
	zipReader, err := zip.NewReader(bytes.NewReader(resultData), int64(len(resultData)))
	require.NoError(t, err)

	// Проверяем, что изображение добавлено
	var hasImage bool
	var hasContentType bool
	var hasRelationship bool

	for _, file := range zipReader.File {
		switch {
		case strings.HasPrefix(file.Name, "word/media/image_embed"):
			hasImage = true
		case file.Name == "[Content_Types].xml":
			content := readZipFileContent(t, file)
			if strings.Contains(content, `Extension="png"`) && strings.Contains(content, `image/png`) {
				hasContentType = true
			}
		case file.Name == "word/_rels/document.xml.rels":
			content := readZipFileContent(t, file)
			if strings.Contains(content, "rImageId1") && strings.Contains(content, "image") {
				hasRelationship = true
			}
		}
	}

	assert.True(t, hasImage, "Image file should be present")
	assert.True(t, hasContentType, "Content type for image should be registered")
	assert.True(t, hasRelationship, "Relationship for image should be created")
}

func TestTemplate_Integration_CustomFunctions(t *testing.T) {
	tmpl := New("test")

	// Добавляем пользовательскую функцию
	tmpl.Funcs(template.FuncMap{
		"upper": strings.ToUpper,
		"multiply": func(a, b int) int {
			return a * b
		},
	})

	// Создаем DOCX с пользовательскими функциями
	customDocx := createCustomFunctionDocx()
	tmpl.ParseDocxFileData(customDocx)
	require.NoError(t, tmpl.err)

	data := TestData{
		Title: "custom functions",
	}

	var result bytes.Buffer
	err := tmpl.Execute(&result, data)
	require.NoError(t, err)

	docXML := extractDocumentXML(t, result.Bytes())
	assert.Contains(t, docXML, "CUSTOM FUNCTIONS") // upper функция
	assert.Contains(t, docXML, "20")               // multiply 4 * 5
}

func TestTemplate_Integration_EmptyTemplate(t *testing.T) {
	// Тест с пустым шаблоном
	emptyDocx := createEmptyDocx()

	tmpl := New("empty")
	tmpl.ParseDocxFileData(emptyDocx)
	require.NoError(t, tmpl.err)

	var result bytes.Buffer
	err := tmpl.Execute(&result, nil)
	require.NoError(t, err)

	// Проверяем, что результат валидный ZIP
	resultData := result.Bytes()
	_, err = zip.NewReader(bytes.NewReader(resultData), int64(len(resultData)))
	assert.NoError(t, err)
}

func TestTemplate_Integration_ErrorHandling(t *testing.T) {
	// Тест с невалидными данными
	invalidDocx := []byte("not a zip file")

	tmpl := New("invalid")
	tmpl.ParseDocxFileData(invalidDocx)
	assert.Error(t, tmpl.err)

	// Попытка выполнить шаблон с ошибкой
	var result bytes.Buffer
	err := tmpl.Execute(&result, nil)
	assert.Error(t, err)
}

// Вспомогательные функции для тестов

func extractDocumentXML(t *testing.T, docxData []byte) string {
	zipReader, err := zip.NewReader(bytes.NewReader(docxData), int64(len(docxData)))
	require.NoError(t, err)

	for _, file := range zipReader.File {
		if file.Name == "word/document.xml" {
			return readZipFileContent(t, file)
		}
	}
	t.Fatal("document.xml not found")
	return ""
}

func readZipFileContent(t *testing.T, file *zip.File) string {
	rc, err := file.Open()
	require.NoError(t, err)
	defer rc.Close()

	content, err := io.ReadAll(rc)
	require.NoError(t, err)
	return string(content)
}

func createCustomFunctionDocx() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	// Content_Types.xml
	contentTypes := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
    <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
    <Default Extension="xml" ContentType="application/xml"/>
    <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`
	writeZipFile(zw, "[Content_Types].xml", contentTypes)

	// _rels/.rels
	rels := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
    <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`
	writeZipFile(zw, "_rels/.rels", rels)

	// word/document.xml с пользовательскими функциями
	document := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
    <w:body>
        <w:p>
            <w:r><w:t>{{upper .Title}}</w:t></w:r>
        </w:p>
        <w:p>
            <w:r><w:t>Result: {{multiply 4 5}}</w:t></w:r>
        </w:p>
    </w:body>
</w:document>`
	writeZipFile(zw, "word/document.xml", document)

	// word/_rels/document.xml.rels
	docRels := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
</Relationships>`
	writeZipFile(zw, "word/_rels/document.xml.rels", docRels)

	zw.Close()
	return buf.Bytes()
}

func createEmptyDocx() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	// Content_Types.xml
	contentTypes := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
    <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
    <Default Extension="xml" ContentType="application/xml"/>
    <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`
	writeZipFile(zw, "[Content_Types].xml", contentTypes)

	// _rels/.rels
	rels := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
    <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`
	writeZipFile(zw, "_rels/.rels", rels)

	// word/document.xml пустой
	document := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
    <w:body>
    </w:body>
</w:document>`
	writeZipFile(zw, "word/document.xml", document)

	// word/_rels/document.xml.rels
	docRels := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
</Relationships>`
	writeZipFile(zw, "word/_rels/document.xml.rels", docRels)

	zw.Close()
	return buf.Bytes()
}

// Benchmark для полного цикла
func BenchmarkTemplate_FullCycle(b *testing.B) {
	docxData := createMinimalDocx()
	data := TestData{
		Title:  "Benchmark Title",
		Footer: "Benchmark Footer",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmpl := New("benchmark")
		tmpl.ParseDocxFileData(docxData)

		var result bytes.Buffer
		err := tmpl.Execute(&result, data)
		if err != nil {
			b.Fatal(err)
		}
	}
}
