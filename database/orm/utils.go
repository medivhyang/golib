package orm

import (
	"reflect"
	"strings"
	"unicode"
)

func containStrings(source []string, target ...string) bool {
	for _, t := range target {
		flag := false
		for _, s := range source {
			if t == s {
				flag = true
				break
			}
		}
		if !flag {
			return false
		}
	}
	return true
}

func repeatString(s string, n int) []string {
	r := make([]string, 0, n)
	for i := 0; i < n; i++ {
		r = append(r, s)
	}
	return r
}

const (
	caseSnake  = "snake"
	caseKebab  = "kebab"
	caseCamel  = "camel"
	casePascal = "pascal"
)

func toCase(typo string, s string, abbrs ...string) string {
	switch typo {
	case caseSnake:
		return strings.Join(parseWords(s), "_")
	case caseKebab:
		return strings.Join(parseWords(s), "-")
	case caseCamel:
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
	case casePascal:
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

func mapStructToMap(src interface{}) (map[string]interface{}, error) {
	value := unrefValue(reflect.ValueOf(src))
	if value.Kind() != reflect.Struct {
		return nil, ErrRequireStructType
	}
	result := map[string]interface{}{}
	for i := 0; i < value.NumField(); i++ {
		result[value.Type().Field(i).Name] = value.Field(i).Interface()
	}
	return result, nil
}

func unrefValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

func unrefValueAndInit(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	return v
}
