package gengo

// 报错行号+3
const goCodeTemplate = `// Auto generated by github.com/davyxu/cellmesh/protogen
// DO NOT EDIT!

package {{.PackageName}}

import (	
	"github.com/davyxu/cellnet"	
	"github.com/davyxu/cellnet/codec"{{range ProtoImportList $}}
	_ "github.com/davyxu/cellnet/codec/{{.}}"{{end}}
	"reflect"
)

// Make compiler import happy
var(
	_ cellnet.Event
	_ codec.CodecRecycler
	_ reflect.Type
)

{{range ServiceGroup $}}
// {{$svcName := .Key}}{{$svcName}}
var ( {{range .Group}}
	Handle_{{ExportSymbolName $svcName}}_{{.Name}} = func(ev cellnet.Event){ panic("'{{.Name}}' not handled") } {{end}}
	Handle_{{ExportSymbolName $svcName}}_Default func(ev cellnet.Event)
)
{{end}}

func GetMessageHandler(svcName string) cellnet.EventCallback {

	switch svcName { {{range ServiceGroup $}}
	case "{{$svcName := .Key}}{{$svcName}}":
		return func(ev cellnet.Event) {
			switch ev.Message().(type) { {{range .Group}}
			case *{{.Name}}:
				Handle_{{ExportSymbolName $svcName}}_{{.Name}}(ev) {{end}}
			default:
				if Handle_{{ExportSymbolName $svcName}}_Default != nil {
					Handle_{{ExportSymbolName $svcName}}_Default(ev)
				}
			}
		} {{end}}
	} 

	return nil
}


func init() {
	{{range .Structs}} {{ if IsMessage . }}
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("{{StructCodec .}}"),	
		Type:  reflect.TypeOf((*{{.Name}})(nil)).Elem(),
		ID:    {{StructMsgID .}},
	}) {{end}} {{end}}
}

`
