package docx

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"text/template"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/xakepp35/pkg/xerrors"
)

// Template представляет DOCX шаблон с возможностью обработки плейсхолдеров и изображений
type Template struct {
	name string
	err  error

	// Компоненты для обработки
	textTemplate  *template.Template
	xmlProcessor  *XMLProcessor
	imageManager  *ImageManager
	contentTypes  *ContentTypesManager
	relationships *RelationshipsManager

	// Буферы для данных
	fileDataBuffer bytes.Buffer
	docFiles       []*zip.File
}

// New создает новый DOCX шаблон
func New(name string) *Template {
	t := &Template{
		name:         name,
		textTemplate: template.New(name),
		xmlProcessor: NewXMLProcessor(),
		imageManager: NewImageManager(),
	}

	// Регистрируем функции для шаблонов
	t.textTemplate.Funcs(template.FuncMap{
		"addImage": t.addImage,
		"add": func(a, b int) int {
			return a + b
		},
	})

	return t
}

func (t *Template) SetDelimiterPair(open, close rune) {
	t.xmlProcessor.SetDelimiterPair(open, close)
}

// ParseDocxFileData парсит данные DOCX файла
func (t *Template) ParseDocxFileData(data []byte) *Template {
	t.err = nil

	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.err = xerrors.New(err, "parse docx file")
		return t
	}

	t.docFiles = make([]*zip.File, 0, len(reader.File))

	var documentXML []byte
	var relationshipsXML []byte
	var contentTypesXML []byte

	// Читаем все необходимые файлы
	for _, file := range reader.File {
		t.docFiles = append(t.docFiles, file)

		switch file.Name {
		case "word/document.xml":
			documentXML, err = t.readZipFile(file)
			if err != nil {
				t.err = xerrors.New(err, "reading document.xml")
				return t
			}
		case "word/_rels/document.xml.rels":
			relationshipsXML, err = t.readZipFile(file)
			if err != nil {
				t.err = xerrors.New(err, "reading document.xml.rels")
				return t
			}
		case "[Content_Types].xml":
			contentTypesXML, err = t.readZipFile(file)
			if err != nil {
				t.err = xerrors.New(err, "reading Content_Types.xml")
				return t
			}
		}
	}

	// Инициализируем менеджеры
	if len(relationshipsXML) > 0 {
		t.relationships = NewRelationshipsManager(relationshipsXML)
	}
	if len(contentTypesXML) > 0 {
		t.contentTypes = NewContentTypesManager(contentTypesXML)
	}

	// Обрабатываем document.xml
	if len(documentXML) > 0 {
		err = t.processDocumentXML(string(documentXML))
		if err != nil {
			t.err = xerrors.New(err, "processing document.xml")
			return t
		}
	}

	return t
}

// readZipFile читает содержимое файла из ZIP архива
func (t *Template) readZipFile(file *zip.File) ([]byte, error) {
	rc, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rc.Close(); closeErr != nil {
			log.Debug().Err(closeErr).Str("file", file.Name).Msg("Failed to close zip file reader")
		}
	}()

	return io.ReadAll(rc)
}

// processDocumentXML обрабатывает XML документа
func (t *Template) processDocumentXML(xmlContent string) error {
	// Применяем обработку XML в правильном порядке
	processedXML := t.xmlProcessor.FixBrokenTemplateKeys(xmlContent)
	processedXML = t.xmlProcessor.PrepareRangeData(processedXML)
	processedXML = t.xmlProcessor.PrepareAddImageData(processedXML)

	// Парсим как Go template
	_, err := t.textTemplate.Parse(processedXML)
	return err
}

// Execute выполняет шаблон и записывает результат в writer
func (t *Template) Execute(w io.Writer, data interface{}) error {
	if t.err != nil {
		return t.err
	}

	// Подготавливаем компоненты к выполнению
	t.prepareForExecution()

	// Выполняем шаблон
	err := t.textTemplate.Execute(&t.fileDataBuffer, data)
	if err != nil {
		return xerrors.New(err, "executing template")
	}

	// Финализируем XML файлы
	t.finalizeXMLFiles()

	// Создаем результирующий DOCX файл
	return t.writeDocxFile(w)
}

// prepareForExecution подготавливает компоненты к выполнению
func (t *Template) prepareForExecution() {
	t.fileDataBuffer.Reset()
	t.imageManager.Reset()

	if t.relationships != nil {
		t.relationships.Reset()
	}
	if t.contentTypes != nil {
		t.contentTypes.Reset()
	}
}

// finalizeXMLFiles завершает формирование XML файлов
func (t *Template) finalizeXMLFiles() {
	if t.relationships != nil {
		t.relationships.Finalize()
	}
	if t.contentTypes != nil {
		t.contentTypes.Finalize()
	}
}

// writeDocxFile записывает DOCX файл
func (t *Template) writeDocxFile(w io.Writer) error {
	docxWriter := zip.NewWriter(w)
	defer func() {
		if err := docxWriter.Close(); err != nil {
			log.Debug().Err(err).Msg("Failed to close DOCX writer")
		}
	}()

	// Записываем все файлы из оригинального архива
	err := t.writeOriginalFiles(docxWriter)
	if err != nil {
		return xerrors.New(err, "writing original files")
	}

	// Записываем изображения
	err = t.writeImages(docxWriter)
	if err != nil {
		return xerrors.New(err, "writing images")
	}

	return nil
}

// writeOriginalFiles записывает файлы из оригинального архива
func (t *Template) writeOriginalFiles(docxWriter *zip.Writer) error {
	for _, file := range t.docFiles {
		var rc io.ReadCloser
		var err error

		// Определяем, какие файлы заменить нашими обработанными версиями
		switch file.Name {
		case "word/document.xml":
			rc = io.NopCloser(&t.fileDataBuffer)
		case "word/_rels/document.xml.rels":
			if t.relationships != nil {
				rc = io.NopCloser(bytes.NewBufferString(t.relationships.String()))
			} else {
				rc, err = file.Open()
			}
		case "[Content_Types].xml":
			if t.contentTypes != nil {
				rc = io.NopCloser(bytes.NewBufferString(t.contentTypes.String()))
			} else {
				rc, err = file.Open()
			}
		default:
			rc, err = file.Open()
		}

		if err != nil {
			return xerrors.Err(err).Str("filename", file.Name).Msg("opening file")
		}

		err = t.writeZipFile(docxWriter, file, rc)
		if closeErr := rc.Close(); closeErr != nil {
			log.Debug().Err(closeErr).Str("file", file.Name).Msg("Failed to close file reader")
		}

		if err != nil {
			return xerrors.Err(err).Str("filename", file.Name).Msg("writing file")
		}
	}
	return nil
}

// writeZipFile записывает один файл в ZIP архив
func (t *Template) writeZipFile(docxWriter *zip.Writer, file *zip.File, rc io.ReadCloser) error {
	fh := file.FileHeader
	fh.Method = zip.Deflate
	fh.Modified = time.Now() // Используем Modified вместо устаревшего SetModTime
	fh.Flags |= 0x800

	wr, err := docxWriter.CreateHeader(&fh)
	if err != nil {
		return err
	}

	_, err = io.Copy(wr, rc)
	return err
}

// writeImages записывает изображения в архив
func (t *Template) writeImages(docxWriter *zip.Writer) error {
	images := t.imageManager.GetImages()
	for i, imageData := range images {
		imageFormat := DetectImageFormat(imageData.Bytes())
		filename := fmt.Sprintf("word/media/image_embed%d.%s", i+1, imageFormat)

		wr, err := docxWriter.Create(filename)
		if err != nil {
			return xerrors.Err(err).Str("filename", filename).Msg("creating image file")
		}

		rc := io.NopCloser(imageData)
		_, err = io.Copy(wr, rc)
		if closeErr := rc.Close(); closeErr != nil {
			log.Debug().Err(closeErr).Str("filename", filename).Msg("Failed to close image data reader")
		}

		if err != nil {
			return xerrors.Err(err).Str("filename", filename).Msg("writing image file")
		}
	}
	return nil
}

// addImage добавляет изображение в документ (функция для шаблонов)
func (t *Template) addImage(imageData []byte, width, height int) string {
	imageID, filename, err := t.imageManager.AddImage(imageData)
	if err != nil {
		// В реальном приложении нужно лучше обрабатывать ошибки
		return fmt.Sprintf("<!-- Error adding image: %v -->", err)
	}

	// Добавляем relationship для изображения
	if t.relationships != nil {
		t.relationships.AddImageRelationship(imageID, filename)
	}

	// Добавляем content type для изображения
	if t.contentTypes != nil {
		imageFormat := DetectImageFormat(imageData)
		err = t.contentTypes.AddImageType(imageFormat)
		if err != nil {
			return fmt.Sprintf("<!-- Error adding content type: %v -->", err)
		}
	}

	return CreateImageXML(imageID, width, height)
}

// Funcs добавляет пользовательские функции в шаблон
func (t *Template) Funcs(funcMap template.FuncMap) *Template {
	t.textTemplate.Funcs(funcMap)
	return t
}
