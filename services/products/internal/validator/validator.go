package validator

import (
	"regexp"
	"unicode"
)

var HyphenatedRegex = regexp.MustCompile("^[a-z0-9]+(?:-[a-z0-9]+)*$")

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: map[string]string{}}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func IsAscii(value string) bool {
	for _, c := range value {
		if c > unicode.MaxASCII {
			return false
		}
	}

	return true
}
