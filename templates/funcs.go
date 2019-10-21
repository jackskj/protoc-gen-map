package templates

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"reflect"
	"text/template"

	protots "github.com/golang/protobuf/ptypes/timestamp"
)

type quoteType int

const (
	Single quoteType = iota
	Double
)

// List of protoc-gen-map helper functions
// Timestamp is required for proto timestamp fields
var functionMap = map[string]interface{}{
	//proto timestamp
	"timestamp": timestamp,
	"date":      date,
	"time":      time,

	//quoting helpers
	"quoteall":  quoteall,
	"squoteall": squoteall,
}

// Adds single quotes to all elements in array
func squoteall(str interface{}) []string {
	return quotedList(str, Single)
}

// Adds double quotes to all elements in array
func quoteall(str interface{}) []string {
	return quotedList(str, Double)
}

// Formats  from proto rimestamp to the standard SQL timestamp format
func timestamp(ts interface{}) interface{} {
	switch ts.(type) {
	case *protots.Timestamp:
		if sqltime, err := ptypes.Timestamp(ts.(*protots.Timestamp)); err == nil {
			return sqltime.Format("2006-01-02 15:04:05")
		} else {
			return ts
		}
	default:
		return ts
	}
}

// Formats  rom proto rimestamp to the standard SQL date format
func date(ts interface{}) interface{} {
	switch ts.(type) {
	case *protots.Timestamp:
		if sqltime, err := ptypes.Timestamp(ts.(*protots.Timestamp)); err == nil {
			return sqltime.Format("2006-01-02")
		} else {
			return ts
		}
	default:
		return ts
	}
}

// Formats from proto rimestamp to the standard SQL time format
func time(ts interface{}) interface{} {
	switch ts.(type) {
	case *protots.Timestamp:
		if sqltime, err := ptypes.Timestamp(ts.(*protots.Timestamp)); err == nil {
			return sqltime.Format("15:04:05")
		} else {
			return ts
		}
	default:
		return ts
	}
}

func Funcs() template.FuncMap {
	return template.FuncMap(functionMap)
}

//inspired by sprig testing
func quotedList(list interface{}, qtype quoteType) []string {
	tp := reflect.TypeOf(list).Kind()
	switch tp {
	case reflect.Slice, reflect.Array:
		l2 := reflect.ValueOf(list)

		l := l2.Len()
		nl := make([]string, l)
		for i := 0; i < l; i++ {
			if qtype == Single {
				nl[i] = fmt.Sprintf("'%v'", l2.Index(i))
			} else {
				nl[i] = fmt.Sprintf("%q", fmt.Sprintf("%v", l2.Index(i)))
			}
		}
		return nl
	default:
		nl := make([]string, 1)
		if qtype == Single {
			nl[0] = fmt.Sprintf("'%v'", list)
		} else {
			nl[0] = fmt.Sprintf("%q", fmt.Sprintf("%v", list))
		}
		return nl
	}
}
