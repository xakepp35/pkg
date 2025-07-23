package docx

import (
	"bytes"
	_ "embed"
	"os"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"
)

//go:embed example.docx
var exampleGoodFile []byte

//go:embed example.png
var exampleImage []byte

func TestTemplate(t *testing.T) {
	docTemplate := New("docx")

	docTemplate.Funcs(template.FuncMap{
		"add": func(a, b int64) int64 {
			return a + b
		},
	})

	docTemplate.ParseDocxFileData(exampleGoodFile)

	buf := new(bytes.Buffer)

	for i := 0; i < 100; i++ {
		buf.Reset()
		err := docTemplate.Execute(buf, DocumentData{
			Title:     "title",
			Footer:    "footer",
			ImageData: exampleImage,

			Items: []Item{
				{
					Id:     1,
					Name:   "a",
					Email:  "111@mail.ru",
					Number: 111,
				},
				{
					Id:     2,
					Name:   "b",
					Email:  "222@mail.ru",
					Number: 222,
				},
				{
					Id:     3,
					Name:   "c",
					Email:  "333@mail.ru",
					Number: 333,
				},
				{
					Id:     4,
					Name:   "d",
					Email:  "444@mail.ru",
					Number: 444,
				},
			},
		})
		require.NoError(t, err)
	}

	err := os.WriteFile("example.result.docx", buf.Bytes(), 0644)
	require.NoError(t, err)
}

type DocumentData struct {
	Title     string
	Footer    string
	Items     []Item
	ImageData []byte
}

type Item struct {
	Id     int64
	Name   string
	Email  string
	Number float64
}
