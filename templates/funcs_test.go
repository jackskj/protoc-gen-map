package templates_test

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/golang/protobuf/ptypes"
	"github.com/jackskj/protoc-gen-map/templates"
	"testing"
	"text/template"
	"time"
)

var sampleTime = time.Date(2014, 11, 14, 9, 15, 52, 0, time.FixedZone("UTC-8", -8*60*60))

func TestSquoteall(t *testing.T) {
	tmpl := "{{ .Ids | squoteall | join \" , \" }}"
	if err := runtv(tmpl, "'1' , '2' , '3' , '4'", map[string]interface{}{"Ids": []int{1, 2, 3, 4}}); err != nil {
		t.Error(err)
	}
	if err := runtv(tmpl, "'1'", map[string]interface{}{"Ids": []int{1}}); err != nil {
		t.Error(err)
	}
	if err := runtv(tmpl, "", map[string]interface{}{"Ids": []int{}}); err != nil {
		t.Error(err)
	}

}
func TestQuoteall(t *testing.T) {
	tmpl := "{{ .Ids | quoteall | join \" , \" }}"
	if err := runtv(tmpl, "\"1\" , \"2\" , \"3\" , \"4\"", map[string]interface{}{"Ids": []int{1, 2, 3, 4}}); err != nil {
		t.Error(err)
	}
	if err := runtv(tmpl, "\"1\"", map[string]interface{}{"Ids": []int{1}}); err != nil {
		t.Error(err)
	}
	if err := runtv(tmpl, "", map[string]interface{}{"Ids": []int{}}); err != nil {
		t.Error(err)
	}
}
func TestTimestamp(t *testing.T) {
	protoTS, _ := ptypes.TimestampProto(sampleTime)
	tmpl := "{{ .Timestamp | timestamp }}"
	if err := runtv(tmpl, "2014-11-14 17:15:52", map[string]interface{}{"Timestamp": protoTS}); err != nil {
		t.Error(err)
	}
}
func TestDate(t *testing.T) {
	protoTS, _ := ptypes.TimestampProto(sampleTime)
	tmpl := "{{ .Timestamp | date }}"
	if err := runtv(tmpl, "2014-11-14", map[string]interface{}{"Timestamp": protoTS}); err != nil {
		t.Error(err)
	}
}
func TestTime(t *testing.T) {
	protoTS, _ := ptypes.TimestampProto(sampleTime)
	tmpl := "{{ .Timestamp | time }}"
	if err := runtv(tmpl, "17:15:52", map[string]interface{}{"Timestamp": protoTS}); err != nil {
		t.Error(err)
	}
}

//inspired by sprig testing
func runtv(tpl, expect string, vars interface{}) error {
	fmap := templates.Funcs()
	t := template.Must(template.New("test").Funcs(sprig.TxtFuncMap()).Funcs(fmap).Parse(tpl))
	var b bytes.Buffer
	err := t.Execute(&b, vars)
	if err != nil {
		return err
	}
	if expect != b.String() {
		return fmt.Errorf("Expected '%s', got '%s'", expect, b.String())
	}
	return nil
}
