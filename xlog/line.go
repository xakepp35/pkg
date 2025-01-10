package xlog

// Line — структура для построения строки лога
type Line struct {
	buf           Buffer // Буфер для текущей строки
	log           *Log   // Родительский логгер
	needSeparator bool   // Нужна ли запятая меджу entries
}

// Bool добавляет булево значение
func (line *Line) Bool(key string, val bool) *Line {
	line.addKey(key)
	line.buf.Bool(val)
	return line
}

// Int добавляет целое число
func (l *Line) Int(key string, val int) *Line {
	l.addKey(key)
	l.buf.Int64(int64(val))
	return l
}

// Float добавляет число с плавающей точкой
func (l *Line) Float(key string, val float64) *Line {
	l.addKey(key)
	l.buf.Float64(val)
	return l
}

// Str добавляет строковое значение
func (l *Line) Str(key string, val string) *Line {
	l.addKey(key)
	l.buf.Byte(LineStringWrapper).String(val).Byte(LineStringWrapper)
	return l
}

// End завершает строку и отправляет её в лог
func (l *Line) End() {
	l.buf.ByteByte(LineObjectCloser, LineNewLine)
	l.log.entries <- l.buf.End()
}

// addKey добавляет ключ
func (l *Line) addKey(key string) {
	if l.needSeparator {
		l.buf.Byte(',')
	}
	l.needSeparator = true
	l.buf.Byte(LineStringWrapper)
	l.buf.String(key)
	l.buf.ByteByte(LineStringWrapper, LinePairSeparator)
}

const (
	LineStringWrapper = '"'
	LineObjectOpener  = '{'
	LineObjectCloser  = '}'
	LineNewLine       = '\n'
	LineItemSeparator = ','
	LinePairSeparator = ':'
)
