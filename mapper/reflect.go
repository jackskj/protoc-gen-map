package mapper

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"reflect"
	"strconv"
	"time"
)

// This set of functions casts a single value reviewed from the DB
// onto a single proto field type and send the proto field.
// Inspired by the fmt package
// Neither proto filed nor value are known types, hence I use type assertions

// TODO, return errors for faulty conversions, for example, casting string to int
// TODO, convert "null" responses (mssql) to nil responses

type Value struct {
	rvalue   reflect.Value
	intVal   int64
	uintVal  uint64
	floatVal float64
	strVal   string
	boolVal  bool
	timeVal  reflect.Value
	err      error
}

func setProto(field reflect.Value, value interface{}) error {
	var v *Value
	switch field.Interface().(type) {
	case int, int8, int16, int32, int64:
		v = NewValue(value, "int")
		field.SetInt(v.intVal)
	case uint, uint8, uint16, uint32, uint64:
		v = NewValue(value, "uint")
		field.SetUint(v.uintVal)
	case float32, float64:
		v = NewValue(value, "float")
		field.SetFloat(v.floatVal)
	case string:
		v = NewValue(value, "string")
		field.SetString(v.strVal)
	case bool:
		v = NewValue(value, "bool")
		field.SetBool(v.boolVal)
	case *timestamp.Timestamp:
		v = NewValue(value, "timestamp")
		field.Set(v.timeVal)
	}
	// TODO, return errors for more faulty conversions, for example, invalid string to int
	// As on now, errors are reurn for invalid datetime
	if v.err != nil {
		return v.err
	} else {
		return nil
	}
}

func NewValue(ivalue interface{}, respType string) *Value {
	v := Value{rvalue: reflect.ValueOf(ivalue)}
	switch respType {
	case "int":
		v.castInt(ivalue)
	case "uint":
		v.castUint(ivalue)
	case "float":
		v.castFLoat(ivalue)
	case "string":
		v.castString(ivalue)
	case "bool":
		v.castBool(ivalue)
	case "timestamp":
		v.castTimestamp(ivalue)
	}
	return &v
}

func (v *Value) castInt(ivalue interface{}) {
	switch ivalue.(type) {
	case int, int8, int16, int32, int64:
		v.intVal = v.rvalue.Int()
	case uint, uint8, uint16, uint32, uint64:
		v.intVal = int64(v.rvalue.Uint())
	case float32, float64:
		v.intVal = int64(v.rvalue.Float())
	case string:
		if s, err := strconv.Atoi(v.rvalue.String()); err == nil {
			v.intVal = int64(s)
		} else {
			v.intVal = int64(0)
		}
	case bool:
		if v.rvalue.Bool() {
			v.intVal = 1
		} else {
			v.intVal = 0
		}
	case time.Time:
		v.intVal = ivalue.(time.Time).Unix()
	}
}

func (v *Value) castUint(ivalue interface{}) {
	switch ivalue.(type) {
	case int, int8, int16, int32, int64:
		v.uintVal = uint64(v.rvalue.Int())
	case uint, uint8, uint16, uint32, uint64:
		v.uintVal = v.rvalue.Uint()
	case float32, float64:
		v.uintVal = uint64(v.rvalue.Float())
	case string:
		if s, err := strconv.Atoi(v.rvalue.String()); err == nil {
			v.uintVal = uint64(s)
		} else {
			v.uintVal = uint64(0)
		}
	case bool:
		if v.rvalue.Bool() {
			v.uintVal = uint64(1)
		} else {
			v.uintVal = uint64(0)
		}
	case time.Time:
		v.uintVal = uint64(ivalue.(time.Time).Unix())
	}
}

func (v *Value) castFLoat(ivalue interface{}) {
	switch ivalue.(type) {
	case int, int8, int16, int32, int64:
		v.floatVal = float64(v.rvalue.Int())
	case uint, uint8, uint16, uint32, uint64:
		v.floatVal = float64(v.rvalue.Uint())
	case float32, float64:
		v.floatVal = v.rvalue.Float()
	case string:
		if s, err := strconv.Atoi(v.rvalue.String()); err == nil {
			v.floatVal = float64(s)
		} else {
			v.floatVal = float64(0)
		}
	case bool:
		if v.rvalue.Bool() {
			v.floatVal = float64(1)
		} else {
			v.floatVal = float64(0)
		}
	case time.Time:
		v.floatVal = float64(ivalue.(time.Time).Unix())
	}
}

func (v *Value) castString(ivalue interface{}) {
	switch ivalue.(type) {
	case int, int8, int16, int32, int64:
		v.strVal = strconv.FormatInt(v.rvalue.Int(), 10)
	case uint, uint8, uint16, uint32, uint64:
		v.strVal = strconv.FormatUint(v.rvalue.Uint(), 10)
	case float32, float64:
		v.strVal = strconv.FormatFloat(v.rvalue.Float(), 'E', -1, 64)
	case string:
		v.strVal = v.rvalue.String()
	case bool:
		if v.rvalue.Bool() {
			v.strVal = "true"
		} else {
			v.strVal = "false"
		}
	case time.Time:
		v.strVal = ivalue.(time.Time).String()
	}
}

func (v *Value) castBool(ivalue interface{}) {
	switch ivalue.(type) {
	case int, int8, int16, int32, int64:
		if v.rvalue.Int() == 0 {
			v.boolVal = false
		} else {
			v.boolVal = true
		}
	case uint, uint8, uint16, uint32, uint64:
		if v.rvalue.Uint() == uint64(0) {
			v.boolVal = false
		} else {
			v.boolVal = true
		}
	case float32, float64:
		if v.rvalue.Float() == float64(0) {
			v.boolVal = false
		} else {
			v.boolVal = true
		}
	case string:
		if v.rvalue.String() == "" {
			v.boolVal = false
		} else {
			v.boolVal = true
		}
	case bool:
		v.boolVal = v.rvalue.Bool()
	case time.Time:
		v.boolVal = true
	}
}

func (v *Value) castTimestamp(ivalue interface{}) {
	switch ivalue.(type) {
	case time.Time:
		if sqlTime, err := ptypes.TimestampProto(ivalue.(time.Time)); err == nil {
			v.timeVal = reflect.ValueOf(sqlTime)
		} else {
			v.timeVal = reflect.ValueOf(ptypes.TimestampNow())
			v.err = err
		}
	default:
		v.timeVal = reflect.ValueOf(ptypes.TimestampNow())
		v.err = errors.New(fmt.Sprintf("cannot convert 	%s of type %s to time.Time", ivalue, reflect.TypeOf(ivalue)))
	}
}
