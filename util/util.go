package util

import (
	"bytes"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func Normalize(s string) []byte {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	s, _, _ = transform.String(t, s)
	s = strings.ToLower(s)
	return []byte(s)
}

func LowerCaseAlpha(s string) []byte {
	result := bytes.NewBuffer(nil)
	for _, x := range Normalize(s) {
		if x >= 'a' && x <= 'z' {
			result.WriteByte(x)
		}
	}
	return result.Bytes()
}

func LowerCaseAlphaOrDot(s string) []byte {
	result := bytes.NewBuffer(nil)
	for _, x := range Normalize(s) {
		if x >= 'a' && x <= 'z' || x == '.' {
			result.WriteByte(x)
		}
	}
	return result.Bytes()
}
