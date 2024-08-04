// Credit: This file is copied from github.com/NathanBaulch/protoc-gen-cobra/naming
// It therefore retains the original copyright and license.
package naming

import (
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
)

type Namer func(string) string

var (
	fixNumbers = func() func(string) string {
		r := regexp.MustCompile(`([a-zA-Z])[-_](\d)`)
		return func(s string) string { return r.ReplaceAllString(s, "$1$2") }
	}()
	Lower      Namer = strings.ToLower
	Upper      Namer = strings.ToUpper
	Pascal     Namer = strcase.ToCamel
	Camel      Namer = strcase.ToLowerCamel
	LowerKebab Namer = func(s string) string { return fixNumbers(strcase.ToKebab(s)) }
	UpperKebab Namer = func(s string) string { return fixNumbers(strcase.ToScreamingKebab(s)) }
	LowerSnake Namer = func(s string) string { return fixNumbers(strcase.ToSnake(s)) }
	UpperSnake Namer = func(s string) string { return fixNumbers(strcase.ToScreamingSnake(s)) }
)

func Composite(s string, namers ...Namer) string {
	for _, namer := range namers {
		s = namer(s)
	}
	return s
}
