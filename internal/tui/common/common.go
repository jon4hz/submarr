package common

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const Ellipsis = "…"

func Title(s string) string {
	return cases.Title(language.English).String(s)
}
