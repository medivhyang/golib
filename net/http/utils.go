package http

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"unicode"
)

var errorPrefix = "http: "

func errorf(format string, args ...interface{}) error {
	return fmt.Errorf(errorPrefix+format, args...)
}

var (
	debug              = false
	debugLogger Logger = log.New(os.Stdout, "", log.LstdFlags)
	debugPrefix        = "http: "
)

type Logger interface {
	Printf(format string, args ...interface{})
}

func EnableDebug(b bool) {
	debug = b
}

func SetLogger(l Logger) {
	debugLogger = l
}

func debugf(format string, args ...interface{}) {
	if debug && debugLogger != nil {
		debugLogger.Printf(debugPrefix+format, args...)
	}
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

func parseStructTagFirstWord(i interface{}, k string) map[string]string {
	rv := reflect.ValueOf(i)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		panic(ErrRequireStructType)
	}
	m := map[string]string{}
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if v, ok := field.Tag.Lookup(k); ok {
			m[field.Name] = strings.Split(v, " ")[0]
		}
	}
	return m
}
