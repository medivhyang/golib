package naming

import (
	"strings"
	"unicode"
)

type Case string

const (
	CaseSnake  Case = "snake"
	CaseKebab  Case = "kebab"
	CaseCamel  Case = "camel"
	CasePascal Case = "pascal"
)

func ToCase(c Case, s string, abbrs ...string) string {
	switch c {
	case CaseSnake:
		return strings.Join(parseWords(s), "_")
	case CaseKebab:
		return strings.Join(parseWords(s), "-")
	case CaseCamel:
		words := parseWords(s)
		for i, word := range words {
			if i == 0 {
				words[i] = strings.ToLower(word)
				continue
			}
			matched := false
			for _, abbr := range abbrs {
				if strings.ToLower(abbr) == strings.ToLower(word) {
					words[i] = strings.ToUpper(abbr)
					matched = true
					break
				}
			}
			if matched {
				continue
			}
			words[i] = upperFirstChar(word)
		}
		return strings.Join(words, "")
	case CasePascal:
		words := parseWords(s)
		for i, word := range words {
			matched := false
			for _, abbr := range abbrs {
				if strings.ToLower(abbr) == strings.ToLower(word) {
					words[i] = strings.ToUpper(abbr)
					matched = true
					break
				}
			}
			if matched {
				continue
			}
			words[i] = upperFirstChar(word)
		}
		return strings.Join(words, "")
	}
	return s
}

func parseWords(s string) (words []string) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return nil
	}
	rs := []rune(s)
	word := ""
	for i := 0; i < len(rs); i++ {
		r := rs[i]
		if r == '_' || r == '-' || r == ' ' {
			if word != "" {
				words = append(words, word)
			}
			word = ""
			continue
		}
		if unicode.IsUpper(r) && ((i-1 > 0 && unicode.IsLower(rs[i-1])) || (i+1 < len(rs) && unicode.IsLower(rs[i+1]))) {
			if word != "" {
				words = append(words, word)
			}
			word = string(r)
			continue
		}
		word += string(r)
	}
	if word != "" {
		words = append(words, word)
	}
	for i := range words {
		words[i] = strings.ToLower(words[i])
	}
	return words
}

func upperFirstChar(s string) string {
	runes := []rune(s)
	return string(unicode.ToUpper(runes[0])) + strings.ToLower(string(runes[1:]))
}
