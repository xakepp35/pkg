package docx

import (
	"bytes"
	"errors"
	"github.com/xakepp35/pkg/xerrors"
	"regexp"
	"strings"
)

// XMLProcessor отвечает за обработку XML содержимого документа
type XMLProcessor struct {
	// Регулярные выражения для работы с XML
	wtRegexp       *regexp.Regexp // для поиска <w:t>...</w:t>
	wrRegexp       *regexp.Regexp // для поиска <w:r>...</w:r>
	rowRegexp      *regexp.Regexp // для поиска строк таблицы
	addImageRegexp *regexp.Regexp // для поиска параграфов с изображениями

	openRune  rune
	closeRune rune
}

// NewXMLProcessor создает новый процессор XML
func NewXMLProcessor() *XMLProcessor {
	return &XMLProcessor{
		wtRegexp:       regexp.MustCompile(`(?s)<w:t[^>]*>(.*?)</w:t>`),
		wrRegexp:       regexp.MustCompile(`(?s)<w:r[\s\S]*?<w:t[\s\S]*?>[\s\S]*?</w:t>[\s\S]*?</w:r>`),
		rowRegexp:      regexp.MustCompile(`(?s)(<w:tr.*?</w:tr>)`),
		addImageRegexp: regexp.MustCompile(`(?s)(<w:p.*?</w:p>)`),
		openRune:       '^',
		closeRune:      '~',
	}
}

func (xp *XMLProcessor) SetDelimiterPair(open, close rune) {
	xp.openRune = open
	xp.closeRune = close
}

func (xp *XMLProcessor) validate(xml string) error {
	open := 0
	for i, c := range xml {
		if c == xp.openRune {
			open++
			if open > 1 {
				return xerrors.Err(nil).Int("place", i).Msg("nested open delimiter")
			}
		}
		if c == xp.closeRune {
			open--
			if open < 0 {
				return xerrors.Err(nil).Int("place", i).Msg("closed delimiter, but not open")
			}
		}
	}

	if open != 0 {
		return errors.New("broken delimiter")
	}

	return nil
}

func (xp *XMLProcessor) FixBrokenTemplateKeys(xml string) string {
	var result bytes.Buffer
	inTemplate := false
	inTag := false

	if !strings.Contains(xml, string(xp.openRune)) && !strings.Contains(xml, string(xp.closeRune)) {
		return xml
	}

	xml = xp.minifyXML(xml)

	for _, c := range xml {
		d := string(c)
		_ = d

		if inTemplate {
			if c == xp.closeRune {
				inTemplate = false
				result.WriteString("}}")
				continue
			}
			if c == '<' {
				inTag = true
				continue
			}
			if c == '>' {
				inTag = false
				continue
			}
			if !inTag {
				result.WriteRune(c)
			}
			continue
		}

		if c == xp.openRune {
			inTemplate = true
			result.WriteString("{{")
			continue
		}

		result.WriteRune(c)
	}

	return result.String()
}

var reBetweenTags = regexp.MustCompile(`>\s+<`)

func (xp *XMLProcessor) minifyXML(xmlContent string) string {
	return reBetweenTags.ReplaceAllString(xmlContent, "><")
}

// isAllowedBetweenRuns проверяет, содержит ли текст только разрешенные элементы между <w:r> тегами
func (xp *XMLProcessor) isAllowedBetweenRuns(text string) bool {
	// Удаляем пробелы и переносы строк
	text = strings.TrimSpace(text)
	if text == "" {
		return true
	}

	// Разрешенные теги между w:r элементами
	allowedTags := []string{
		"<w:proofErr",
		"</w:proofErr>",
		"<w:softHyphen/>",
		"<w:noBreakHyphen/>",
	}

	// Создаем временную строку для обработки
	remaining := text

	for len(remaining) > 0 {
		found := false
		for _, tag := range allowedTags {
			if strings.HasPrefix(remaining, tag) {
				if strings.Contains(tag, "/>") {
					// Самозакрывающийся тег
					remaining = remaining[len(tag):]
				} else if strings.HasPrefix(tag, "</") {
					// Закрывающий тег
					remaining = remaining[len(tag):]
				} else {
					// Открывающий тег - нужно найти его конец
					endPos := strings.Index(remaining, ">")
					if endPos == -1 {
						return false
					}
					remaining = remaining[endPos+1:]
				}
				found = true
				break
			}
		}

		if !found {
			// Если нашли неразрешенный контент, возвращаем false
			return false
		}

		// Убираем пробелы после обработки тега
		remaining = strings.TrimSpace(remaining)
	}

	return true
}

// mergeTemplateRuns объединяет run'ы с плейсхолдером
func (xp *XMLProcessor) mergeTemplateRuns(result *strings.Builder, xml string, lastEnd, bufferStart int, buffer string, matches [][]int, i, j int) {
	result.WriteString(xml[lastEnd:bufferStart])

	// Заменяем <w:t>...</w:t> внутри первого <w:r>
	firstRun := xml[bufferStart:matches[i][1]]
	firstRunNew := xp.wtRegexp.ReplaceAllString(firstRun, "<w:t>"+buffer+"</w:t>")
	result.WriteString(firstRunNew)

	// Остальные <w:r> делаем пустыми
	for k := i + 1; k <= j; k++ {
		emptyRun := xp.wtRegexp.ReplaceAllString(xml[matches[k][0]:matches[k][1]], "<w:t></w:t>")
		result.WriteString(emptyRun)
	}
}

// cleanupEmptyRuns удаляет пустые run'ы и добавляет необходимые
func (xp *XMLProcessor) cleanupEmptyRuns(xml string) string {
	// Удаляем все run'ы, которые содержат только пустые <w:t> (без реального текста)
	emptyRunRegexp := regexp.MustCompile(`(?s)<w:r[^>]*>.*?</w:r>`)
	resultStr := emptyRunRegexp.ReplaceAllStringFunc(xml, func(match string) string {
		// Проверяем, есть ли в этом run'е <w:t> с реальным содержимым
		hasContentT := regexp.MustCompile(`<w:t[^>]*>[^<\s]+.*?</w:t>`)
		if hasContentT.MatchString(match) {
			return match // оставляем как есть - есть реальный текст
		}

		// Проверяем, есть ли хотя бы один <w:t> (пустой или с пробелами)
		hasAnyT := regexp.MustCompile(`<w:t[^>]*>.*?</w:t>`)
		if hasAnyT.MatchString(match) {
			return "" // удаляем - содержит только пустые <w:t>
		}
		return match // оставляем - нет <w:t> вообще
	})

	// Добавляем пустой run в параграфы, которые остались без run'ов
	emptyParaRegexp := regexp.MustCompile(`(?s)(<w:p[^>]*>.*?<w:pPr>.*?</w:pPr>)\s*</w:p>`)
	resultStr = emptyParaRegexp.ReplaceAllStringFunc(resultStr, func(match string) string {
		// Проверяем, есть ли в этом параграфе хотя бы один <w:r>
		if strings.Contains(match, "<w:r") {
			return match // есть run'ы — всё ок
		}
		// Добавляем пустой run перед закрывающим </w:p>
		return strings.Replace(match, "</w:p>", "<w:r><w:t></w:t></w:r></w:p>", 1)
	})

	return resultStr
}

// PrepareRangeData обрабатывает циклы "{{range" в таблицах
func (xp *XMLProcessor) PrepareRangeData(src string) string {
	rows := xp.rowRegexp.FindAllString(src, -1)
	var rebuilt []string

	for _, row := range rows {
		if strings.Contains(row, "{{range") || strings.Contains(row, "{{end}}") {
			re := regexp.MustCompile(`{{\s*(range\s+[^}]+|end)\s*}}`)
			clean := re.FindString(row)
			rebuilt = append(rebuilt, clean)
		} else {
			rebuilt = append(rebuilt, row)
		}
	}

	// полностью заменить все <w:tr> блоки на обновлённые
	return xp.rowRegexp.ReplaceAllStringFunc(src, func(_ string) string {
		if len(rebuilt) == 0 {
			return ""
		}
		s := rebuilt[0]
		rebuilt = rebuilt[1:]
		return s
	})
}

// PrepareAddImageData обрабатывает функции addImage
func (xp *XMLProcessor) PrepareAddImageData(src string) string {
	rows := xp.addImageRegexp.FindAllString(src, -1)
	var rebuilt []string

	for _, row := range rows {
		if strings.Contains(row, "{{addImage") {
			re := regexp.MustCompile(`{{\s*addImage.*?}}`)
			clean := re.FindString(row)
			rebuilt = append(rebuilt, clean)
		} else {
			rebuilt = append(rebuilt, row)
		}
	}

	return xp.addImageRegexp.ReplaceAllStringFunc(src, func(_ string) string {
		if len(rebuilt) == 0 {
			return ""
		}
		s := rebuilt[0]
		rebuilt = rebuilt[1:]
		return s
	})
}
