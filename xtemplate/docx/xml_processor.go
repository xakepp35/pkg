package docx

import (
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
func (xp *XMLProcessor) FixBrokenTemplateKeys(xml string) string {
	matches := xp.wrRegexp.FindAllStringIndex(xml, -1)
	if len(matches) == 0 {
		return xml
	}

	var result strings.Builder
	lastEnd := 0
	n := len(matches)
	i := 0

	for i < n {
		start, end := matches[i][0], matches[i][1]
		run := xml[start:end]

		// Извлекаем текст из <w:t>
		tMatch := xp.wtRegexp.FindStringSubmatch(run)
		if tMatch == nil || len(tMatch) < 2 {
			result.WriteString(xml[lastEnd:end])
			lastEnd = end
			i++
			continue
		}

		text := tMatch[1]
		openIdx := strings.Index(text, "{{")
		if openIdx == -1 {
			// Не начало плейсхолдера — просто копируем
			result.WriteString(xml[lastEnd:end])
			lastEnd = end
			i++
			continue
		}

		// Кандидат на склейку
		buffer := text
		bufferStart := start
		bufferEnd := end
		j := i + 1

		for j < n {
			// Проверяем, что между end предыдущего и start следующего только whitespace
			inter := xml[bufferEnd:matches[j][0]]
			if strings.TrimSpace(inter) != "" {
				break
			}

			nextRun := xml[matches[j][0]:matches[j][1]]
			nextTMatch := xp.wtRegexp.FindStringSubmatch(nextRun)
			if nextTMatch == nil || len(nextTMatch) < 2 {
				break
			}

			nextText := nextTMatch[1]
			buffer += nextText
			bufferEnd = matches[j][1]

			if strings.Contains(buffer, "}}") {
				// Нашли конец плейсхолдера, склеиваем
				xp.mergeTemplateRuns(&result, xml, lastEnd, bufferStart, buffer, matches, i, j)
				lastEnd = bufferEnd
				i = j + 1
				break
			}
			j++
		}

		if lastEnd < bufferStart {
			// Не удалось склеить — выводим как есть
			result.WriteString(xml[lastEnd:end])
			lastEnd = end
			i++
		}
	}

	if lastEnd < len(xml) {
		result.WriteString(xml[lastEnd:])
	}

	resultStr := result.String()
	return xp.cleanupEmptyRuns(resultStr)
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
