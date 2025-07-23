# DOCX Template Package

Пакет для работы с DOCX шаблонами, обработки плейсхолдеров и работы с изображениями в документах Word.

## Возможности

- ✅ **Обработка разбитых плейсхолдеров** - автоматическое склеивание плейсхолдеров типа `{{.Title}}`, разбитых между `<w:r>` элементами
- ✅ **Добавление изображений** - вставка изображений с автоматической регистрацией MIME-типов и relationships
- ✅ **Поддержка циклов** - обработка конструкций `{{range}}` и `{{end}}` в таблицах
- ✅ **Пользовательские функции** - возможность добавления собственных функций шаблонов
- ✅ **Валидная XML структура** - сохранение корректной структуры WordprocessingML

## Архитектура

Пакет разделен на логические модули:

### Основные компоненты

- **`Template`** - основной класс для работы с DOCX шаблонами
- **`XMLProcessor`** - обработка XML содержимого, склейка плейсхолдеров
- **`ImageManager`** - управление изображениями в документе
- **`ContentTypesManager`** - управление MIME-типами в `[Content_Types].xml`
- **`RelationshipsManager`** - управление связями в `document.xml.rels`

### Структура файлов

```
docx/
├── template.go              # Основной класс Template
├── xml_processor.go         # Обработка XML и плейсхолдеров
├── image_utils.go          # Работа с изображениями
├── content_types.go        # Управление Content_Types.xml
├── relationships.go        # Управление relationships
├── template_test.go        # Интеграционные тесты
├── xml_processor_test.go   # Тесты XMLProcessor
├── image_utils_test.go     # Тесты ImageManager
├── content_types_test.go   # Тесты ContentTypesManager
└── integration_test.go     # Полные интеграционные тесты
```

## Использование

### Базовое использование

```go
import "internal/utils/docx"

// Читаем DOCX файл
docxData, err := os.ReadFile("template.docx")
if err != nil {
    log.Fatal(err)
}

// Создаем шаблон
tmpl := docx.New("my-template")

// Добавляем пользовательские функции (опционально)
tmpl.Funcs(template.FuncMap{
    "upper": strings.ToUpper,
    "add": func(a, b int) int { return a + b },
})

// Парсим DOCX
tmpl.ParseDocxFileData(docxData)

// Подготавливаем данные
data := struct {
    Title  string
    Items  []Item
    Image  []byte
}{
    Title: "Мой документ",
    Items: []Item{{Name: "Пункт 1"}, {Name: "Пункт 2"}},
    Image: imageData,
}

// Выполняем шаблон
var result bytes.Buffer
err = tmpl.Execute(&result, data)
if err != nil {
    log.Fatal(err)
}

// Сохраняем результат
os.WriteFile("result.docx", result.Bytes(), 0644)
```

### Работа с изображениями

В DOCX шаблоне используйте функцию `addImage`:

```
{{addImage .ImageData 200 150}}
```

Где:
- `.ImageData` - поле с данными изображения (`[]byte`)
- `200` - ширина в пикселях
- `150` - высота в пикселях

### Работа с циклами

Для создания таблиц с динамическими данными:

```xml
<w:tr>
    {{range .Items}}
</w:tr>
<w:tr>
    <w:tc><w:p><w:r><w:t>{{.Name}}</w:t></w:r></w:p></w:tc>
    <w:tc><w:p><w:r><w:t>{{.Price}}</w:t></w:r></w:p></w:tc>
</w:tr>
<w:tr>
    {{end}}
</w:tr>
```

## Обработка разбитых плейсхолдеров

Основная проблема, которую решает пакет - Word часто разбивает плейсхолдеры между XML элементами:

**Проблема:**
```xml
<w:r><w:t>{{.Ti</w:t></w:r>
<w:r><w:t>tle}}</w:t></w:r>
```

**Решение:**
```xml
<w:r><w:t>{{.Title}}</w:t></w:r>
```

XMLProcessor автоматически:
1. Находит разбитые плейсхолдеры
2. Склеивает их в первом `<w:r>` элементе
3. Удаляет пустые `<w:r>` элементы
4. Добавляет необходимые пустые элементы для валидности XML

## Тестирование

Пакет покрыт comprehensive тестами:

```bash
# Запуск всех тестов
go test -v

# Запуск конкретной группы тестов
go test -v -run TestXMLProcessor
go test -v -run TestImageManager
go test -v -run TestIntegration

# Запуск бенчмарков
go test -bench=.
```

### Типы тестов

- **Unit тесты** - тестирование отдельных компонентов
- **Integration тесты** - полный цикл работы с DOCX
- **Edge case тесты** - граничные случаи и ошибки
- **Benchmarks** - производительность

## Известные ограничения

1. **Data sharing в ImageManager** - текущая реализация не создает копию данных изображений
2. **Сложные вложенные плейсхолдеры** - поддерживаются только базовые конструкции Go templates
3. **Большие файлы** - весь DOCX загружается в память

## Лицензия

Внутренний пакет для проекта invoicer-be. 