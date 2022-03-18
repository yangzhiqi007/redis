package jsondata

import (
	"github.com/davyxu/protoplus/codegen"
	"github.com/davyxu/tabtoy/v3/model"
)

func Generate(globals *model.Globals) (data []byte, err error) {

	err = codegen.NewCodeGen("jsontext").
		RegisterTemplateFunc(codegen.UsefulFunc).
		RegisterTemplateFunc(UsefulFunc).
		ParseTemplate(templateText, globals).
		WriteBytes(&data).Error()

	return
}
