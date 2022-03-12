package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var b strings.Builder
	var last rune
	for _, r := range s {
		if unicode.IsDigit(r) {
			if last == 0 {
				return "", ErrInvalidString
			}

			counter, _ := strconv.Atoi(string(r))
			if counter > 1 {
				b.WriteString(strings.Repeat(string(last), counter-1))
			} else if counter == 0 {
				current := b.String()
				current = current[:len(current)-len(string(last))]
				b.Reset()
				b.WriteString(current)
			}

			last = 0
		} else {
			b.WriteRune(r)
			last = r
		}
	}

	return b.String(), nil
}
