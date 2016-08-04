package validations

import (
	"github.com/astaxie/beego/validation"
)

type ValidFormer interface {
	Valid(*validation.Validation) (bool, error)
}
