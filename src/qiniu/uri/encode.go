package uri

import (
	"syscall"
)

const (
	needEscape = 0xff
	dontEscape = 16
)

const (
	escapeChar = '\''
)

func genEncoding() []byte {
	var encoding [256]byte
	for c := 0; c <= 0xff; c++ {
		encoding[c] = needEscape
	}
	for c := 'a'; c <= 'f'; c++ {
		encoding[c] = byte(c - ('a' - 10))
	}
	for c := 'A'; c <= 'F'; c++ {
		encoding[c] = byte(c - ('A' - 10))
	}
	for c := 'g'; c <= 'z'; c++ {
		encoding[c] = dontEscape
	}
	for c := 'G'; c <= 'Z'; c++ {
		encoding[c] = dontEscape
	}
	for c := '0'; c <= '9'; c++ {
		encoding[c] = byte(c - '0')
	}
	for _, c := range []byte{'-', '_', '.', '~', '*', '(', ')', '$', '&', '+', ',', ':', ';', '=', '@'} {
		encoding[c] = dontEscape
	}
	encoding['/'] = '!'
	return encoding[:]
}

var encoding = genEncoding()

func encode(v string) string {
	n := 0
	hasEscape := false
	for i := 0; i < len(v); i++ {
		c := v[i]
		switch encoding[c] {
		case needEscape:
			n++
		case '!':
			hasEscape = true
		}
	}
	if !hasEscape && n == 0 {
		return v
	}

	t := make([]byte, len(v)+2*n)
	j := 0
	for i := 0; i < len(v); i++ {
		c := v[i]
		switch encoding[c] {
		case needEscape:
			t[j] = escapeChar
			t[j+1] = "0123456789ABCDEF"[c>>4]
			t[j+2] = "0123456789ABCDEF"[c&15]
			j += 3
		case '!':
			t[j] = encoding[c]
			j++
		default:
			t[j] = c
			j++
		}
	}
	return string(t)
}

func decode(s string) (v string, err error) {
	n := 0
	hasEscape := false
	for i := 0; i < len(s); {
		switch s[i] {
		case escapeChar:
			n++
			if i+2 >= len(s) || encoding[s[i+1]] >= 16 || encoding[s[i+2]] >= 16 {
				return "", syscall.EINVAL
			}
			i += 3
		case '!':
			hasEscape = true
			i++
		default:
			i++
		}
	}
	if !hasEscape && n == 0 {
		return s, nil
	}

	t := make([]byte, len(s)-2*n)

	j := 0
	for i := 0; i < len(s); {
		switch s[i] {
		case escapeChar:
			t[j] = (encoding[s[i+1]] << 4) | encoding[s[i+2]]
			i += 3
		case '!':
			t[j] = '/'
			i++
		default:
			t[j] = s[i]
			i++
		}
		j++
	}
	return string(t), nil
}
