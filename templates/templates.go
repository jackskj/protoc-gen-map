package templates

import (
	"log"
	"text/template"
)

func GeneratorTemplates() *template.Template {
	tmpl := template.New("")
	tmpl.New("server").Parse(server)
	tmpl.New("initFunc").Parse(initFunc)
	tmpl.New("enumValueMaps").Parse(enumValueMaps)
	_, err := tmpl.New("unaryResponse").Parse(unaryResponse)
	if err != nil {
		log.Println(err)

	}
	tmpl.New("streamingResponse").Parse(streamingResponse)
	return tmpl
}
