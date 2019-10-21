package plugin

import (
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/jackskj/protoc-gen-map/templates"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

func (p *SqlPlugin) GetSQLTemplates() error {
	sqlLocations := strings.Split(p.Param["sql"], ",")
	if len(sqlLocations) < 1 {
		p.Error("missing location of SQL templates. \n " +
			"Use --sqlmap_out=sql=my/sql,template/locations to identify sql template directories")
	}

	var sqlFiles []string
	for _, loc := range sqlLocations {
		err := filepath.Walk(loc, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			sqlFiles = append(sqlFiles, path)
			return nil
		})
		if err != nil {
			p.Error(err.Error())
		}
	}

	if len(sqlFiles) < 1 {
		p.Error("No SQL files found or missing location of SQL templates.\n" +
			"Use --sqlmap_out=sql=my/sql,template/locations to identify sql template directories")
	}
	var sql strings.Builder
	for _, sqlFile := range sqlFiles {
		sqlr, _ := ioutil.ReadFile(sqlFile)
		if _, err := template.New("sqlTemplate").Funcs(sprig.TxtFuncMap()).Funcs(templates.Funcs()).Parse(string(sqlr)); err != nil {
			p.Error(fmt.Sprintf("error parsing template file %s: %s", sqlFile, err))
		}
		sql.WriteString(string(sqlr))
	}
	p.SqlTemplates = sql.String()
	return nil
}

func (p *SqlPlugin) TestSQLTemplates(g *generator.Generator) {
	rpcNames := make(map[string]bool)
	for _, fd := range g.AllFiles().GetFile() {
		for _, svc := range fd.GetService() {
			for _, rpc := range svc.GetMethod() {
				rpcNames[rpc.GetName()] = true
			}
		}
	}
	sqlTemplate, err := template.New("sqlTemplate").Funcs(sprig.TxtFuncMap()).Funcs(templates.Funcs()).Parse(p.SqlTemplates)
	if err != nil {
		p.Error(fmt.Sprintf("template parsing error: %s", err))
	}
	tplStr := strings.Trim(sqlTemplate.DefinedTemplates(), "; defined templates are: ")
	tpls := strings.Split(tplStr, ",")
	//removing quotes
	for i, tpl := range tpls {
		trimmedTpl := strings.TrimSpace(tpl)
		tpls[i] = trimmedTpl[1 : len(trimmedTpl)-1]
	}
	var matchingNames []string
	for _, name := range tpls {
		if rpcNames[name] {
			matchingNames = append(matchingNames, name)
		}
	}
	sort.Strings(matchingNames)
	if len(matchingNames) != 0 {
		log.Printf("protoc-gen-map: following rpc-sql template pairs found \n" + strings.Join(matchingNames, ", "))
		for _, name := range matchingNames {
			p.matchingTemplates[name] = false
		}
	} else {
		p.Error(fmt.Sprint("protoc-gen-map: no rpc-sql template pairs found \n",
			"head over to https://github.com/jackskj/protoc-gen-map#sqlproto-definition "+
				"for instructions on rpc-sql definition guide.",
		))
	}
}

// Generate sql template only in the first protofile
func (p *SqlPlugin) PrintSQLTemplates(file *generator.FileDescriptor) error {
	if p.Request.FileToGenerate[0] != file.GetName() {
		return nil
	}
	p.setTemplateImports()
	p.P(fmt.Sprintf("var %[1]s, _ = template.New(\"%[1]s\").Funcs(sprig.TxtFuncMap()).Funcs(mappertmpl.Funcs()).Parse(`", p.SqlTemplateName))
	p.P(p.SqlTemplates)
	p.P("`)")
	return nil
}

func (p *SqlPlugin) GenerateSQLTemplates() error {
	p.P(fmt.Sprintf("var %s = `", p.SqlTemplateName))
	p.P(p.SqlTemplates)

	p.P("`")
	return nil
}
