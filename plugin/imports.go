package plugin

import (
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

var imports = map[alias]string{
	"bytes":      "bytes",
	"context":    "context",
	"codes":      "google.golang.org/grpc/codes",
	"status":     "google.golang.org/grpc/status",
	"mapper":     "github.com/jackskj/protoc-gen-map/mapper",
	"mappertmpl": "github.com/jackskj/protoc-gen-map/templates",
	"sprig":      "github.com/Masterminds/sprig",
	"log":        "log",
	"sync":       "sync",
	"sql":        "database/sql",
	"template":   "text/template",
}

func (p *SqlPlugin) GenerateImports(file *generator.FileDescriptor) {
	p.P()
	p.P("\t //protoc-gen-map packages")
	for alias, _ := range p.Pkg {
		p.PrintImport(
			generator.GoPackageName(alias),
			generator.GoImportPath(imports[alias]),
		)
	}
	p.Pkg = make(map[alias]bool)
}

func (p *SqlPlugin) setStreamingImports() {
	p.Pkg["bytes"] = true
	p.Pkg["codes"] = true
	p.Pkg["status"] = true
	p.Pkg["log"] = true
	p.Pkg["context"] = true
}

func (p *SqlPlugin) setUnaryImports() {
	p.Pkg["bytes"] = true
	p.Pkg["codes"] = true
	p.Pkg["status"] = true
	p.Pkg["log"] = true
	p.Pkg["context"] = true
}
func (p *SqlPlugin) setTemplateImports() {
	p.Pkg["mappertmpl"] = true
	p.Pkg["template"] = true
	p.Pkg["sprig"] = true
}

func (p *SqlPlugin) setEnumImports() {
	p.Pkg["mapper"] = true
}
