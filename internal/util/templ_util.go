package util

import (
	"fmt"

	"github.com/a-h/templ"
)

func ConstructTemplUrl(format string, args ...any) templ.SafeURL {
	return templ.URL(fmt.Sprintf(format, args...))
}
