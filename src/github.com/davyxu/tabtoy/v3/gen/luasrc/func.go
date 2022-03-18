package luasrc

import (
	"github.com/davyxu/tabtoy/v3/gen"
	"github.com/davyxu/tabtoy/v3/model"
	"strings"
	"text/template"
)

var UsefulFunc = template.FuncMap{}

func WrapValue(globals *model.Globals, cell *model.Cell, valueType *model.TypeDefine) string {
	if valueType.IsArray() {

		var sb strings.Builder
		sb.WriteString("[")

		if cell != nil {
			for index, elementValue := range cell.ValueList {
				if index > 0 {
					sb.WriteString(",")
				}
				sb.WriteString(gen.WrapSingleValue(globals, valueType, elementValue))
			}
		}

		sb.WriteString("]")

		return sb.String()

	} else {

		var value string
		if cell != nil {
			value = cell.Value
		}

		return gen.WrapSingleValue(globals, valueType, value)
	}
}

func init() {
	UsefulFunc["WrapTabValue"] = func(globals *model.Globals, dataTable *model.DataTable, allHeaders []*model.TypeDefine, row, col int) (ret string) {

		// 找到完整的表头（按完整表头遍历）
		header := allHeaders[col]

		if header == nil {
			return ""
		}

		// 在单元格找到值
		valueCell := dataTable.GetCell(row, col)

		if valueCell != nil {

			return WrapValue(globals, valueCell, header)
		} else {
			// 这个表中没有这列数据
			return WrapValue(globals, nil, header)
		}
	}

}
