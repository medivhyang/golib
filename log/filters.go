package log

import (
	"strings"
	"time"
)

func MakeLevelFilter(l Level) Filter {
	return func(t Template) bool {
		if !l.Valid() || l > t.Level {
			return false
		}
		return true
	}
}

func MakeModuleFilter(modules ...string) Filter {
	return func(t Template) bool {
		if len(modules) == 0 {
			return true
		}
		if len(modules) < len(t.Modules) {
			return false
		}
		for i := 0; i < len(modules); i++ {
			if modules[i] != t.Modules[i] {
				return false
			}
		}
		return true
	}
}

func MakeTimeFilter(begin, end time.Time) Filter {
	return func(t Template) bool {
		if begin.After(end) {
			return false
		}
		return t.Time.After(begin) && t.Time.Before(end)
	}
}

func MakeMessageFilter(keyword string) Filter {
	return func(t Template) bool {
		return strings.Contains(t.Message, keyword)
	}
}
