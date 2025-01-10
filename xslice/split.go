package xslice

func SplitBytesLines(s []byte) [][]byte {
	return SplitBytes(s, '\n')
}

// Optimized split for a single byte separator (strict O(N), no extra iterations)
func SplitBytes(s []byte, sep byte) [][]byte {
	n := 1
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			n++
		}
	}
	a := make([][]byte, n)
	start := 0
	count := 0
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			a[count] = s[start:i]
			start = i + 1
			count++
		}
	}
	// edge case
	a[n-1] = s[start:]
	return a
}

func LastIndexByteString(s string, c byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == c {
			return i
		}
	}
	return -1
}
