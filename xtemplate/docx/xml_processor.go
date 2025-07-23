package docx

import (
	"regexp"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/xml"
)

// XMLProcessor отвечает за обработку XML содержимого документа
type XMLProcessor struct {
	// Регулярные выражения для работы с XML
	wtRegexp       *regexp.Regexp // для поиска <w:t>...</w:t>
	wrRegexp       *regexp.Regexp // для поиска <w:r>...</w:r>
	rowRegexp      *regexp.Regexp // для поиска строк таблицы
	addImageRegexp *regexp.Regexp // для поиска параграфов с изображениями
}

// NewXMLProcessor создает новый процессор XML
func NewXMLProcessor() *XMLProcessor {
	return &XMLProcessor{
		wtRegexp:       regexp.MustCompile(`(?s)<w:t[^>]*>(.*?)</w:t>`),
		wrRegexp:       regexp.MustCompile(`(?s)<w:r[\s\S]*?<w:t[\s\S]*?>[\s\S]*?</w:t>[\s\S]*?</w:r>`),
		rowRegexp:      regexp.MustCompile(`(?s)(<w:tr.*?</w:tr>)`),
		addImageRegexp: regexp.MustCompile(`(?s)(<w:p.*?</w:p>)`),
	}
}

// FixBrokenTemplateKeys объединяет разбитые плейсхолдеры типа {{.Title}}
// Использует regexp для поиска и замены разорванных шаблонов
func (xp *XMLProcessor) FixBrokenTemplateKeys(xml string) string {
	// Сначала убираем все форматирование XML для упрощения обработки
	result := xp.minifyXML(xml)

	// Повторяем до тех пор, пока есть изменения
	changed := true
	for changed {
		oldResult := result

		// Ищем все возможные части шаблонов
		// Паттерн для поиска { + что-то + { + что-то + } + что-то + }
		templatePattern := regexp.MustCompile(`\{[^{}]*\{[^{}]*\}[^{}]*\}`)

		result = templatePattern.ReplaceAllStringFunc(result, func(match string) string {
			// Удаляем все XML теги из match
			cleanContent := regexp.MustCompile(`<[^>]*>`).ReplaceAllString(match, "")
			cleanContent = strings.TrimSpace(cleanContent)

			// Проверяем что это валидный шаблон
			if len(cleanContent) >= 4 && cleanContent[:2] == "{{" && cleanContent[len(cleanContent)-2:] == "}}" {
				// Извлекаем содержимое между {{ и }}
				templateContent := strings.TrimSpace(cleanContent[2 : len(cleanContent)-2])
				if templateContent != "" {
					// Находим первый w:t тег в оригинальном match
					firstTagPattern := regexp.MustCompile(`<w:t>[^<]*</w:t>`)
					firstTag := firstTagPattern.FindString(match)

					if firstTag != "" {
						// Заменяем содержимое первого тега на восстановленный шаблон
						newFirstTag := regexp.MustCompile(`(<w:t>)[^<]*(</w:t>)`).ReplaceAllString(firstTag, "${1}{{"+templateContent+"}}${2}")

						// Возвращаем только первый тег с восстановленным шаблоном
						return newFirstTag
					}

					return "{{" + templateContent + "}}"
				}
			}

			return match
		})

		// Убираем дублированные теги w:t
		result = regexp.MustCompile(`<w:t><w:t>`).ReplaceAllString(result, `<w:t>`)
		result = regexp.MustCompile(`</w:t></w:t>`).ReplaceAllString(result, `</w:t>`)

		changed = (oldResult != result)
	}

	return result
}

// minifyXML удаляет все пробелы, табы и переносы строк между XML тегами
func (xp *XMLProcessor) minifyXML(xmlContent string) string {
	m := minify.New()
	m.AddFunc("text/xml", xml.Minify)

	result, err := m.String("text/xml", xmlContent)
	if err != nil {
		// Если минификация не удалась, возвращаем оригинал
		return xmlContent
	}

	return result
}

// isLetter проверяет, является ли символ буквой
func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
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
