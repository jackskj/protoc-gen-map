package templates

import (
	"text/template"
)

func GeneratorTemplates() *template.Template {
	tmpl := template.New("")
	tmpl.New("server").Parse(server)
	tmpl.New("initFunc").Parse(initFunc)
	tmpl.New("enumValueMaps").Parse(enumValueMaps)
	tmpl.New("unaryResponse").Parse(unaryResponse)
	tmpl.New("streamingResponse").Parse(streamingResponse)
	return tmpl
}
