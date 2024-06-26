package jxorm

import (
	"bytes"
	"strings"
	"unicode"
)

type Strings struct {
	source string
}

func (s Strings) Source() string {
	return s.source
}

// ToCamel converts the input text into camel case
func (s Strings) ToCamel() string {
	list := s.splitBy(func(r rune) bool {
		return r == '_'
	}, true)
	var target []string
	for _, item := range list {
		target = append(target, From(item).Title())
	}
	return strings.Join(target, "")
}

// From converts the input text to String and returns it
func From(data string) Strings {
	return Strings{source: data}
}

// it will not ignore spaces
func (s Strings) splitBy(fn func(r rune) bool, remove bool) []string {
	if s.IsEmptyOrSpace() {
		return nil
	}
	var list []string
	buffer := new(bytes.Buffer)
	for _, r := range s.source {
		if fn(r) {
			if buffer.Len() != 0 {
				list = append(list, buffer.String())
				buffer.Reset()
			}
			if !remove {
				buffer.WriteRune(r)
			}
			continue
		}
		buffer.WriteRune(r)
	}
	if buffer.Len() != 0 {
		list = append(list, buffer.String())
	}
	return list
}

// Untitle return the original string if rune is not letter at index 0
func (s Strings) Untitle() string {
	if s.IsEmptyOrSpace() {
		return s.source
	}
	r := rune(s.source[0])
	if !unicode.IsUpper(r) && !unicode.IsLower(r) {
		return s.source
	}
	return string(unicode.ToLower(r)) + s.source[1:]
}

// Title calls the strings.Title
func (s Strings) Title() string {
	if s.IsEmptyOrSpace() {
		return s.source
	}
	return strings.Title(s.source)
}

// IsEmptyOrSpace returns true if the length of the string value is 0 after call strings.TrimSpace, or else returns false
func (s Strings) IsEmptyOrSpace() bool {
	if len(s.source) == 0 {
		return true
	}
	if strings.TrimSpace(s.source) == "" {
		return true
	}
	return false
}

// Lower calls the strings.ToLower
func (s Strings) Lower() string {
	return strings.ToLower(s.source)
}

// Upper calls the strings.ToUpper
func (s Strings) Upper() string {
	return strings.ToUpper(s.source)
}

// ToSnake converts the input text into snake case
func (s Strings) ToSnake() string {
	list := s.splitBy(unicode.IsUpper, false)
	var target []string
	for _, item := range list {
		target = append(target, From(item).Lower())
	}
	return strings.Join(target, "_")
}
