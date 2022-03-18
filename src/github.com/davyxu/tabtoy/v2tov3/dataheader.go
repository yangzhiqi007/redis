package v2tov3

import (
	"github.com/davyxu/golexer"
	"github.com/davyxu/tabtoy/v2tov3/model"
	"github.com/davyxu/tabtoy/v3/helper"
	v3model "github.com/davyxu/tabtoy/v3/model"
	"github.com/tealeg/xlsx"
	"strings"
)

func importDataHeader(globals *model.Globals, sourceSheet, targetSheet *xlsx.Sheet, tableName string) (headerList []model.ObjectFieldType) {

	var headerRow *xlsx.Row

	// 遍历所有行
	for col := 0; ; col++ {

		var oft model.ObjectFieldType
		oft.ObjectType = tableName
		oft.Kind = v3model.TypeUsage_HeaderStruct

		oft.FieldName = helper.GetSheetValueString(sourceSheet, 0, col)

		// 空列，终止
		if oft.FieldName == "" {
			break
		}

		// 列头中带有#的，特别是最后一行
		if strings.HasPrefix(oft.FieldName, "#") {
			continue
		}

		if headerRow == nil {
			headerRow = targetSheet.AddRow()
		}

		oft.FieldType = helper.GetSheetValueString(sourceSheet, 1, col)

		// 元信息
		meta := helper.GetSheetValueString(sourceSheet, 2, col)

		oft.Meta = golexer.NewKVPair()
		if err := oft.Meta.Parse(meta); err != nil {
			continue
		}

		if strings.HasPrefix(oft.FieldType, "[]") {
			oft.FieldType = oft.FieldType[2:]
			oft.ArraySplitter = oft.Meta.GetString("ListSpliter")

			if oft.ArraySplitter == "" {
				log.Warnln("array list no ListSpliter:", oft.FieldName, oft.ObjectType)
			}
		}

		oft.Name = helper.GetSheetValueString(sourceSheet, 3, col)

		if oft.Name == "" {
			log.Warnf("v2的字段注释为空, %s | %s", oft.FieldName, tableName)
			oft.Name = oft.FieldName
		}

		var disabledForV3 string

		// 添加V3表头
		if globals.TypeIsNoneKind(oft.FieldType) {
			disabledForV3 = "#"
		}

		// 结构体等类型，标记为nong，输出为#
		if !model.IsNativeType(oft.FieldType) {

			targetOft := globals.ObjectTypeByName(oft.FieldType)
			// 类型已经被前置定义，且不是枚举（那就是结构体）时，标记为空，后面不会被使用
			if targetOft != nil && targetOft.Kind != v3model.TypeUsage_Enum {
				oft.Kind = v3model.TypeUsage_None
			}

		}

		// 新表的表头加列
		headerRow.AddCell().SetValue(disabledForV3 + oft.Name)

		// 拆分字段填充的数组
		if !globals.SourceTypeExists(oft.ObjectType, oft.FieldName) {

			globals.AddSourceType(oft)

		}

		headerList = append(headerList, oft)
	}

	return
}
