package docx

import (
	"fmt"
	"strings"
)

// ContentTypesManager управляет Content_Types.xml в docx файле
type ContentTypesManager struct {
	original []byte
	buffer   strings.Builder
}

// NewContentTypesManager создает новый менеджер для Content_Types.xml
func NewContentTypesManager(original []byte) *ContentTypesManager {
	return &ContentTypesManager{
		original: original,
	}
}

// Reset сбрасывает буфер и готовит для новых изменений
func (ctm *ContentTypesManager) Reset() {
	ctm.buffer.Reset()
	_, _ = ctm.buffer.WriteString(strings.Replace(string(ctm.original), "</Types>", "", -1))
}

// AddImageType добавляет MIME-тип для изображения, если его еще нет
func (ctm *ContentTypesManager) AddImageType(extension string) error {
	var contentType string
	switch extension {
	case "png":
		contentType = "image/png"
	case "jpg", "jpeg":
		contentType = "image/jpeg"
	case "gif":
		contentType = "image/gif"
	case "webp":
		contentType = "image/webp"
	case "bmp":
		contentType = "image/bmp"
	default:
		contentType = "image/png"
	}

	// Проверяем, есть ли уже такой тип в Content_Types
	defaultEntry := fmt.Sprintf(`<Default Extension="%s" ContentType="%s"/>`, extension, contentType)
	if !strings.Contains(ctm.buffer.String(), defaultEntry) {
		ctm.buffer.WriteString(defaultEntry)
	}

	return nil
}

// Finalize завершает формирование XML и возвращает результат
func (ctm *ContentTypesManager) Finalize() string {
	ctm.buffer.WriteString("</Types>")
	return ctm.buffer.String()
}

// String возвращает текущее содержимое буфера
func (ctm *ContentTypesManager) String() string {
	return ctm.buffer.String()
}
