package mapper_test

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"
	"text/template"

	"github.com/Masterminds/sprig"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jackskj/protoc-gen-map/mapper"
	mappertmpl "github.com/jackskj/protoc-gen-map/templates"
	td "github.com/jackskj/protoc-gen-map/testdata"
)

type ParamTest struct {
	NilVal     *string
	ShortVal   string
	ComplexVal string
	LongVal    string
	IntVal     int
	TimeVal    timestamp.Timestamp
	SliceVal   []string
}

func TestParametarizedQuery(t *testing.T) {
	testingVal := ParamTest{
		ShortVal:   "a",
		ComplexVal: "!#$%&()*+,-./:;<=>?@[\\]^_{|}~",
		LongVal:    strings.Repeat("a", 100),
		IntVal:     6,
		TimeVal:    td.GetSampleTS(),
		SliceVal:   []string{"a", "b", "c", "d", "e"},
	}
	paramBuff := &bytes.Buffer{}

	_ = sqlTemplate.ExecuteTemplate(paramBuff, "ParamTest", testingVal)
	for _, dialect := range []string{"postgres", "mysql"} {
		if preparedSql, args, err := mapper.PrepareQuery(dialect, paramBuff.Bytes()); err == nil {
			testResults["ParamMap_Raw_"+dialect] = paramBuff.String()
			testResults["ParamMap_Args_"+dialect] = fmt.Sprintf("%s", args)
			testResults["ParamMap_Prepared_"+dialect] = preparedSql
		} else {
			log.Panicln("Parameterize Error" + err.Error())
		}
	}
	for _, dialect := range []string{"", "unknown"} {
		if _, _, err := mapper.PrepareQuery(dialect, paramBuff.Bytes()); err == nil {
			log.Panicln("Expected error")
		} else {
			testResults["ParamMap_Error_Raw_"+dialect] = err.Error()
		}
	}
}

var sqlTemplate, _ = template.New("sqlTemplate").Funcs(sprig.TxtFuncMap()).Funcs(mappertmpl.Funcs()).Parse(`
{{ define "ParamTest" }}
select {{ param .NilVal }}

select {{ param "" }}

select {{ param .ShortVal }}

select {{ param .ComplexVal }}

select {{ param .LongVal }}

select {{ param .IntVal }}

select {{ param .TimeVal }}

select {{ param .SliceVal }}

select
{{- range .SliceVal }}
{{- param . -}} , 
{{- end }} 
{{ end }} 
`)
