package plugin

import (
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/jackskj/protoc-gen-map/templates"
	"log"
	"os"
	"text/template"
)

var _ generator.Plugin = (*SqlPlugin)(nil)

type (
	alias string // go package alias
)

type SqlPlugin struct {
	*generator.Generator

	SqlTemplates    string
	SqlTemplateName string
	Pkg             map[alias]bool

	matchingTemplates map[string]bool
	currentPackage    string
	method            *Method            // currently generated method
	server            *Server            // currently generated service server
	genTemplate       *template.Template // templates generator
}

func New() (*SqlPlugin, error) {
	return &SqlPlugin{
		SqlTemplateName:   "sqlTemplate",
		Pkg:               make(map[alias]bool),
		matchingTemplates: make(map[string]bool), // [templateName]isMatched
	}, nil
}

func (p *SqlPlugin) Init(g *generator.Generator) {
	p.Generator = g
	p.GetSQLTemplates()
	p.TestSQLTemplates(g)
	p.genTemplate = templates.GeneratorTemplates()
}

func (p *SqlPlugin) Name() string {
	return "map"
}

func (g *SqlPlugin) Error(errorMsg string) {
	log.Print("protoc-gen-map error: ", errorMsg)
	os.Exit(1)
}
