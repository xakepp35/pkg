package docx

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContentTypesManager_NewContentTypesManager(t *testing.T) {
	original := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
    <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
    <Default Extension="xml" ContentType="application/xml"/>
</Types>`)

	ctm := NewContentTypesManager(original)
	assert.NotNil(t, ctm)
	assert.Equal(t, original, ctm.original)
}

func TestContentTypesManager_Reset(t *testing.T) {
	original := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
    <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
</Types>`)

	ctm := NewContentTypesManager(original)
	ctm.Reset()

	// Проверяем, что в буфере есть содержимое без закрывающего тега
	content := ctm.String()
	assert.Contains(t, content, `<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">`)
	assert.Contains(t, content, `<Default Extension="rels"`)
	assert.NotContains(t, content, "</Types>")
}

func TestContentTypesManager_AddImageType(t *testing.T) {
	original := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
    <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
</Types>`)

	ctm := NewContentTypesManager(original)
	ctm.Reset()

	tests := []struct {
		name         string
		extension    string
		expectedType string
	}{
		{"PNG image", "png", "image/png"},
		{"JPEG image", "jpg", "image/jpeg"},
		{"JPEG image alt", "jpeg", "image/jpeg"},
		{"GIF image", "gif", "image/gif"},
		{"WebP image", "webp", "image/webp"},
		{"BMP image", "bmp", "image/bmp"},
		{"Unknown format", "unknown", "image/png"}, // fallback to PNG
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ctm.AddImageType(tt.extension)
			require.NoError(t, err)

			content := ctm.String()
			expectedEntry := `<Default Extension="` + tt.extension + `" ContentType="` + tt.expectedType + `"/>`
			assert.Contains(t, content, expectedEntry)
		})
	}
}

func TestContentTypesManager_AddImageType_NoDuplicates(t *testing.T) {
	original := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
</Types>`)

	ctm := NewContentTypesManager(original)
	ctm.Reset()

	// Добавляем PNG дважды
	err := ctm.AddImageType("png")
	require.NoError(t, err)

	err = ctm.AddImageType("png")
	require.NoError(t, err)

	// Проверяем, что PNG entry есть только один раз
	content := ctm.String()
	pngEntry := `<Default Extension="png" ContentType="image/png"/>`
	count := strings.Count(content, pngEntry)
	assert.Equal(t, 1, count, "PNG entry should appear only once")
}

func TestContentTypesManager_Finalize(t *testing.T) {
	original := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
</Types>`)

	ctm := NewContentTypesManager(original)
	ctm.Reset()

	err := ctm.AddImageType("png")
	require.NoError(t, err)

	result := ctm.Finalize()

	// Проверяем структуру
	assert.Contains(t, result, `<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">`)
	assert.Contains(t, result, `<Default Extension="png" ContentType="image/png"/>`)
	assert.Contains(t, result, "</Types>")

	// Проверяем, что закрывающий тег в конце
	assert.True(t, strings.HasSuffix(result, "</Types>"))
}
