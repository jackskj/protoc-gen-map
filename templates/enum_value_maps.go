package templates

var enumValueMaps = `
var EnumValueMaps = map[string]map[string]int32{
	{{- range $enumName, $valueMap := .EnumValueMaps }}
	"{{ $enumName }}": map[string]int32{
		{{- range $value, $valueInt := $valueMap }}
		"{{ $value }}": {{ $valueInt }},
		{{- end }}
	},
	{{- end }}
}
`
