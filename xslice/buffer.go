package xslice

const (
	BufferInitialSize = 64
)

type Buffer struct {
	buf []byte
}

func NewBuffer() *Buffer {
	return NewBufferSized(BufferInitialSize)
}

func NewBufferSized(size int) *Buffer {
	return &Buffer{
		buf: make([]byte, 0, size),
	}
}

func (b *Buffer) End() []byte {
	return b.buf
}

// Byte добавляет один байт
func (b *Buffer) Byte(val byte) *Buffer {
	b.grow(1)
	b.buf[len(b.buf)-1] = val
	return b
}

// ByteByte добавляет два байта подряд
func (b *Buffer) ByteByte(val1, val2 byte) *Buffer {
	b.grow(2)
	b.buf[len(b.buf)-2] = val1
	b.buf[len(b.buf)-1] = val2
	return b
}

// Bytes добавляет срез байтов
func (b *Buffer) Bytes(val []byte) *Buffer {
	n := len(val)
	b.grow(n)
	copy(b.buf[len(b.buf)-n:], val)
	return b
}

// String добавляет строку
func (b *Buffer) String(val string) *Buffer {
	n := len(val)
	b.grow(n)
	copy(b.buf[len(b.buf)-n:], val)
	return b
}

const (
	bufferRune1 = 1 << 7
	bufferRune2 = 1 << 11
	bufferRune3 = 1 << 16
)

// Rune добавляет руну (UNICODE символ)
func (b *Buffer) Rune(val rune) *Buffer {
	switch {
	case val < bufferRune1:
		b.grow(1)
		b.buf[len(b.buf)-1] = byte(val)
	case val < bufferRune2:
		b.grow(2)
		b.buf[len(b.buf)-2] = 0xC0 | byte(val>>6)
		b.buf[len(b.buf)-1] = 0x80 | byte(val&0x3F)
	case val < bufferRune3:
		b.grow(3)
		b.buf[len(b.buf)-3] = 0xE0 | byte(val>>12)
		b.buf[len(b.buf)-2] = 0x80 | byte((val>>6)&0x3F)
		b.buf[len(b.buf)-1] = 0x80 | byte(val&0x3F)
	default:
		b.grow(4)
		b.buf[len(b.buf)-4] = 0xF0 | byte(val>>18)
		b.buf[len(b.buf)-3] = 0x80 | byte((val>>12)&0x3F)
		b.buf[len(b.buf)-2] = 0x80 | byte((val>>6)&0x3F)
		b.buf[len(b.buf)-1] = 0x80 | byte(val&0x3F)
	}
	return b
}

// Int64 добавляет целое число
func (b *Buffer) Int64(val int64) *Buffer {
	var tmp [20]byte // Максимум 20 символов для int64
	sign := val < 0
	if sign {
		val = -val
	}
	i := len(tmp)
	for {
		i--
		tmp[i] = '0' + byte(val%10)
		val /= 10
		if val == 0 {
			break
		}
	}
	if sign {
		i--
		tmp[i] = '-'
	}
	b.Bytes(tmp[i:])
	return b
}

// Float64 добавляет число с плавающей точкой (режим 'g', максимальная точность)
func (b *Buffer) Float64(val float64) *Buffer {
	if val < 0 {
		b.Byte('-')
		val = -val
	}

	// Обработка специальных значений
	if val != val { // NaN
		return b.String("NaN")

	}
	if val > 1e308 { // +Inf
		return b.String("+Inf")
	}
	if val < -1e308 { // -Inf
		return b.String("-Inf")
	}
	// Буфер для записи числа
	var tmp [32]byte
	pos := len(tmp)
	// Разбиение на целую и дробную части
	intPart := int64(val)
	fracPart := val - float64(intPart)
	// Обработка целой части
	for intPart != 0 {
		pos--
		tmp[pos] = '0' + byte(intPart%10)
		intPart /= 10
	}
	if pos == len(tmp) {
		pos--
		tmp[pos] = '0'
	}
	// Обработка дробной части
	if fracPart != 0 {
		b.Bytes(tmp[pos:])
		b.Byte('.')
		for i := 0; i < 16; i++ {
			fracPart *= 10
			digit := int64(fracPart)
			b.Byte('0' + byte(digit))
			fracPart -= float64(digit)
			if fracPart == 0 {
				break
			}
		}
		return b
	}
	// Записываем целую часть
	return b.Bytes(tmp[pos:])
}

// Bool добавляет значение true или false
func (b *Buffer) Bool(val bool) *Buffer {
	if val {
		b.String("true")
		return b
	}
	b.String("false")
	return b
}

// grow увеличивает буфер на n байт
func (b *Buffer) grow(n int) {
	l := len(b.buf)
	newL := n + l
	// Попытка увеличить буфер через реслайсинг
	if newL <= cap(b.buf) {
		b.buf = b.buf[:newL]
		return
	}
	newCap := newL
	if newCap < cap(b.buf)*2 {
		newCap = cap(b.buf) * 2
	}
	newBuf := make([]byte, l, newCap)
	copy(newBuf, b.buf)
	b.buf = newBuf[:newL]
}
