package docx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXMLProcessor_NewXMLProcessor(t *testing.T) {
	xp := NewXMLProcessor()
	assert.NotNil(t, xp)
	assert.NotNil(t, xp.wtRegexp)
	assert.NotNil(t, xp.wrRegexp)
	assert.NotNil(t, xp.rowRegexp)
	assert.NotNil(t, xp.addImageRegexp)
}

func TestXMLProcessor_FixBrokenTemplateKeys_NoTemplate(t *testing.T) {
	xp := NewXMLProcessor()

	// XML без шаблонов
	xml := `<w:p><w:r><w:t>Simple text</w:t></w:r></w:p>`
	result := xp.FixBrokenTemplateKeys(xml)
	assert.Equal(t, xml, result)
}

func TestXMLProcessor_FixBrokenTemplateKeys_CompleteTemplate(t *testing.T) {
	xp := NewXMLProcessor()

	// XML с полным шаблоном в одном run
	xml := `<w:p><w:r><w:t>^.Title~</w:t></w:r></w:p>`
	result := xp.FixBrokenTemplateKeys(xml)
	expected := `<w:p><w:r><w:t>{{.Title}}</w:t></w:r></w:p>`
	assert.Equal(t, expected, result)
}

func TestXMLProcessor_FixBrokenTemplateKeys_BrokenAcrossRuns(t *testing.T) {
	xp := NewXMLProcessor()

	// XML с разбитым шаблоном - функция объединяет содержимое
	xml := `<w:p><w:r><w:t>^.Ti</w:t></w:r><w:r><w:t>tle~</w:t></w:r></w:p>`
	result := xp.FixBrokenTemplateKeys(xml)

	// Должен объединить содержимое шаблона в один run
	expected := `<w:p><w:r><w:t>{{.Title}}</w:t></w:r></w:p>`
	assert.Equal(t, expected, result)
}

func TestXMLProcessor_FixBrokenTemplateKeys_MultipleRuns(t *testing.T) {
	xp := NewXMLProcessor()

	// Шаблон разбит на 3 run'а - функция объединяет содержимое
	xml := `<w:p><w:r><w:t>^.</w:t></w:r><w:r><w:t>Tit</w:t></w:r><w:r><w:t>le~</w:t></w:r></w:p>`
	result := xp.FixBrokenTemplateKeys(xml)

	// Должен объединить содержимое шаблона в один run
	expected := `<w:p><w:r><w:t>{{.Title}}</w:t></w:r></w:p>`
	assert.Equal(t, expected, result)
}

func TestXMLProcessor_FixBrokenTemplateKeys_EmptyRunsCleanup(t *testing.T) {
	processor := NewXMLProcessor()

	input := `<w:p>
		<w:r><w:t>Normal text</w:t></w:r>
		<w:r><w:rPr><w:lang w:val="en-US"/></w:rPr><w:t></w:t></w:r>
		<w:r><w:t>  </w:t></w:r>
	</w:p>`

	result := processor.FixBrokenTemplateKeys(input)

	// Новая логика не изменяет XML без шаблонов
	expected := `<w:p>
		<w:r><w:t>Normal text</w:t></w:r>
		<w:r><w:rPr><w:lang w:val="en-US"/></w:rPr><w:t></w:t></w:r>
		<w:r><w:t>  </w:t></w:r>
	</w:p>`

	assert.Equal(t, expected, result)
}

func TestXMLProcessor_FixBrokenTemplateKeys_EmptyParagraphFix(t *testing.T) {
	xp := NewXMLProcessor()

	// Параграф с пустым run'ом
	xml := `<w:p w14:paraId="123"><w:pPr><w:rPr><w:lang w:val="en-US"/></w:rPr></w:pPr><w:r><w:rPr><w:lang w:val="en-US"/></w:rPr><w:t></w:t></w:r></w:p>`

	result := xp.FixBrokenTemplateKeys(xml)

	// Функция не должна изменять XML без шаблонов
	assert.Equal(t, xml, result)
}

func TestXMLProcessor_PrepareRangeData(t *testing.T) {
	xp := NewXMLProcessor()

	xml := `<w:tbl>
		<w:tr><w:tc><w:p><w:r><w:t>Header</w:t></w:r></w:p></w:tc></w:tr>
		<w:tr><w:tc><w:p><w:r><w:t>{{range .Items}}</w:t></w:r></w:p></w:tc></w:tr>
		<w:tr><w:tc><w:p><w:r><w:t>{{.Name}}</w:t></w:r></w:p></w:tc></w:tr>
		<w:tr><w:tc><w:p><w:r><w:t>{{end}}</w:t></w:r></w:p></w:tc></w:tr>
	</w:tbl>`

	result := xp.PrepareRangeData(xml)

	// Должен заменить строки с range/end на чистые команды
	assert.Contains(t, result, "{{range .Items}}")
	assert.Contains(t, result, "{{end}}")
	assert.NotContains(t, result, "<w:tc><w:p><w:r><w:t>{{range .Items}}</w:t></w:r></w:p></w:tc>")
}

func TestXMLProcessor_PrepareAddImageData(t *testing.T) {
	xp := NewXMLProcessor()

	xml := `<w:p><w:r><w:t>{{addImage .ImageData 200 200}}</w:t></w:r></w:p>
		<w:p><w:r><w:t>Regular text</w:t></w:r></w:p>`

	result := xp.PrepareAddImageData(xml)

	// Должен заменить параграф с addImage на чистую команду
	assert.Contains(t, result, "{{addImage .ImageData 200 200}}")
	assert.Contains(t, result, "Regular text")
}

func TestXMLProcessor_FixBrokenTemplateKeys_EdgeCases(t *testing.T) {
	xp := NewXMLProcessor()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No broken templates",
			input:    `<w:p><w:r><w:t>Normal text</w:t></w:r></w:p>`,
			expected: `<w:p><w:r><w:t>Normal text</w:t></w:r></w:p>`,
		},
		{
			name:     "Template_with_spaces",
			input:    `<w:p><w:r><w:t>^ .Ti</w:t></w:r><w:r><w:t>tle ~</w:t></w:r></w:p>`,
			expected: `<w:p><w:r><w:t>{{ .Title }}</w:t></w:r></w:p>`,
		},
		{
			name:     "Multiple templates",
			input:    `<w:p><w:r><w:t>^.Na</w:t></w:r><w:r><w:t>me~ and ^.Ti</w:t></w:r><w:r><w:t>tle~</w:t></w:r></w:p>`,
			expected: `<w:p><w:r><w:t>{{.Name}} and {{.Title}}</w:t></w:r></w:p>`,
		},
		{
			name:     "Template with XML tags inside",
			input:    `<w:p><w:r><w:t>^<w:r>test</w:r>~</w:t></w:r></w:p>`,
			expected: `<w:p><w:r><w:t>{{test}}</w:t></w:r></w:p>`,
		},
		{
			name:     "Template with nested tags",
			input:    `<w:p><w:r><w:t>^<w:r><w:t>test</w:t></w:r>~</w:t></w:r></w:p>`,
			expected: `<w:p><w:r><w:t>{{test}}</w:t></w:r></w:p>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := xp.FixBrokenTemplateKeys(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestXMLProcessor_isAllowedBetweenRuns(t *testing.T) {
	processor := NewXMLProcessor()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: true,
		},
		{
			name:     "Only whitespace",
			input:    "  \n\t  ",
			expected: true,
		},
		{
			name:     "ProofErr tags",
			input:    `<w:proofErr w:type="spellStart"/><w:proofErr w:type="spellEnd"/>`,
			expected: true,
		},
		{
			name:     "Closing proofErr tag",
			input:    `</w:proofErr>`,
			expected: true,
		},
		{
			name:     "Complex proofErr sequence",
			input:    `<w:proofErr w:type="spellStart"/><w:proofErr w:type="spellEnd"/><w:proofErr w:type="gramEnd"/>`,
			expected: true,
		},
		{
			name:     "Non-allowed content",
			input:    `<w:r><w:t>text</w:t></w:r>`,
			expected: false,
		},
		{
			name:     "Mixed allowed and non-allowed",
			input:    `<w:proofErr w:type="spellStart"/>some text`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.isAllowedBetweenRuns(tt.input)
			if result != tt.expected {
				t.Errorf("isAllowedBetweenRuns(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestXMLProcessor_FixBrokenTemplateKeys_Debug(t *testing.T) {
	processor := NewXMLProcessor()

	// Тестируем простой случай
	input := `<w:p><w:r><w:t>^.Title~</w:t></w:r></w:p>`

	result := processor.FixBrokenTemplateKeys(input)
	t.Logf("Simple case:")
	t.Logf("Input:  %s", input)
	t.Logf("Result: %s", result)

	// Тестируем разорванный случай
	input2 := `<w:p><w:r><w:t>^.Ti</w:t></w:r><w:r><w:t>tle~</w:t></w:r></w:p>`

	result2 := processor.FixBrokenTemplateKeys(input2)
	t.Logf("Broken case:")
	t.Logf("Input:  %s", input2)
	t.Logf("Result: %s", result2)
}

func TestXMLProcessor_FixBrokenTemplateKeys_SimplePrefix(t *testing.T) {
	processor := NewXMLProcessor()

	// Простой случай без переносов строк
	input := `<w:p><w:r><w:t>Title: ^.Ti</w:t></w:r><w:r><w:t>tle~</w:t></w:r></w:p>`

	result := processor.FixBrokenTemplateKeys(input)
	t.Logf("Simple prefix case:")
	t.Logf("Input:  %s", input)
	t.Logf("Result: %s", result)

	// Проверяем что содержимое шаблона объединилось
	expected := `<w:p><w:r><w:t>Title: {{.Title}}</w:t></w:r></w:p>`
	assert.Equal(t, expected, result)
}

func TestXMLProcessor_FixBrokenTemplateKeys_CustomDelimiters(t *testing.T) {
	processor := NewXMLProcessor()

	// Устанавливаем кастомные разделители
	processor.SetDelimiterPair('{', '}')

	// Тестируем с кастомными разделителями
	input := `<w:p><w:r><w:t>{.Title}</w:t></w:r></w:p>`
	result := processor.FixBrokenTemplateKeys(input)

	// Должен преобразовать в стандартные разделители
	expected := `<w:p><w:r><w:t>{{.Title}}</w:t></w:r></w:p>`
	assert.Equal(t, expected, result)
}

func TestXMLProcessor_FixBrokenTemplateKeys_ComplexNestedTags(t *testing.T) {
	processor := NewXMLProcessor()

	// Тестируем сложный случай с вложенными тегами внутри шаблона
	input := `<w:p><w:r><w:t>^<w:r><w:t>Name</w:t></w:r>: <w:r><w:t>Value</w:t></w:r>~</w:t></w:r></w:p>`
	result := processor.FixBrokenTemplateKeys(input)

	// Должен извлечь только текст, игнорируя XML теги
	expected := `<w:p><w:r><w:t>{{Name: Value}}</w:t></w:r></w:p>`
	assert.Equal(t, expected, result)
}

func TestXMLProcessor_FixBrokenTemplateKeys_MultipleTemplatesInOneRun(t *testing.T) {
	processor := NewXMLProcessor()

	// Тестируем несколько шаблонов в одном run
	input := `<w:p><w:r><w:t>^Name~: ^Value~</w:t></w:r></w:p>`
	result := processor.FixBrokenTemplateKeys(input)

	// Должен обработать оба шаблона
	expected := `<w:p><w:r><w:t>{{Name}}: {{Value}}</w:t></w:r></w:p>`
	assert.Equal(t, expected, result)
}

func TestXMLProcessor_FixBrokenTemplateKeys_InsideXMLTags(t *testing.T) {
	processor := NewXMLProcessor()

	// Тестируем случай, когда разделители находятся внутри XML тегов
	input := `<w:p><w:r><w:t>^</w:t></w:r><w:r><w:t>Name</w:t></w:r><w:r><w:t>~</w:t></w:r></w:p>`
	result := processor.FixBrokenTemplateKeys(input)

	// Должен объединить содержимое шаблона в один run
	expected := `<w:p><w:r><w:t>{{Name}}</w:t></w:r></w:p>`
	assert.Equal(t, expected, result)
}

func TestXMLProcessor_FixBrokenTemplateKeys_OnlyClosingDelimiter(t *testing.T) {
	processor := NewXMLProcessor()

	// Тестируем случай только с закрывающим разделителем
	input := `<w:p><w:r><w:t>Name~</w:t></w:r></w:p>`
	result := processor.FixBrokenTemplateKeys(input)

	// Должен оставить как есть, так как нет открывающего разделителя
	expected := `<w:p><w:r><w:t>Name~</w:t></w:r></w:p>`
	assert.Equal(t, expected, result)
}

func TestXMLProcessor_FixBrokenTemplateKeys_SimpleTest(t *testing.T) {
	processor := NewXMLProcessor()

	// Простой тест для понимания логики
	input := `<w:p><w:r><w:t>^test~</w:t></w:r></w:p>`
	result := processor.FixBrokenTemplateKeys(input)
	t.Logf("Simple test:")
	t.Logf("Input:  %s", input)
	t.Logf("Result: %s", result)

	// Тест с разбитым шаблоном
	input2 := `<w:p><w:r><w:t>^te</w:t></w:r><w:r><w:t>st~</w:t></w:r></w:p>`
	result2 := processor.FixBrokenTemplateKeys(input2)
	t.Logf("Broken test:")
	t.Logf("Input:  %s", input2)
	t.Logf("Result: %s", result2)
}
