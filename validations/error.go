package validations

import (
	"github.com/astaxie/beego/validation"
	"strings"
)

/*
	add dash id at all errors.Name without it ex : Nom => facturation-Nom
	with param id = facturation and sym = -
*/
func ErrorAddId(valid *validation.Validation, id, sym string) {
	for k, v := range valid.Errors {
		if !strings.Contains(valid.Errors[k].Field, sym) {
			valid.Errors[k].Field = id + sym + v.Field
			valid.Errors[k].Key = id + sym + v.Key
		}
	}
}

