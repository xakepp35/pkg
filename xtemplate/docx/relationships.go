package docx

import (
	"fmt"
	"strings"
)

// RelationshipsManager управляет document.xml.rels файлом
type RelationshipsManager struct {
	original []byte
	buffer   strings.Builder
}

// NewRelationshipsManager создает новый менеджер для relationships
func NewRelationshipsManager(original []byte) *RelationshipsManager {
	return &RelationshipsManager{
		original: original,
	}
}

// Reset сбрасывает буфер и готовит для новых изменений
func (rm *RelationshipsManager) Reset() {
	rm.buffer.Reset()
	_, _ = rm.buffer.WriteString(strings.Replace(string(rm.original), "</Relationships>", "", -1))
}

// AddImageRelationship добавляет связь для изображения
func (rm *RelationshipsManager) AddImageRelationship(imageID, filename string) {
	relationshipXML := fmt.Sprintf(
		`<Relationship Id="%s" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="media/%s"/>`,
		imageID, filename,
	)
	rm.buffer.WriteString(relationshipXML)
}

// Finalize завершает формирование XML и возвращает результат
func (rm *RelationshipsManager) Finalize() string {
	rm.buffer.WriteString("</Relationships>")
	return rm.buffer.String()
}

// String возвращает текущее содержимое буфера
func (rm *RelationshipsManager) String() string {
	return rm.buffer.String()
}
