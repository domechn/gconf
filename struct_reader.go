//Author : dmc
//
//Date: 2018/9/4 上午10:27
//
//Description:
package gconf

import (
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"strings"
)

const (
	str          = "string"
	strs         = "[]string"
	it           = "int"
	it64         = "int64"
	flt          = "float"
	bol          = "bool"
	defaultValue = "default"
	jsonc        = "json"
)

type read2struct struct {
	cf Configer
}

// Read2Struct 将配置信息读取到结构体(暂时不支持组合)
func Read2Struct(cf Configer, out interface{}) error {
	r2s := &read2struct{cf: cf}
	rv := reflect.ValueOf(out)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("invaild received param , need ptr")
	}

	err := r2s.value(rv)

	return err
}

func (r2s *read2struct) value(rv reflect.Value) error {
	te := rv.Type().Elem()
	for i := 0; i < te.NumField(); i++ {
		var err error
		field := te.Field(i)
		name := field.Name
		js := field.Tag.Get(jsonc)
		dv := field.Tag.Get(defaultValue)
		rf := rv.Elem().Field(i)
		if js != "" {
			name = js
		}
		err = r2s.assign(name, dv, rf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r2s *read2struct) assign(key, dv string, rv reflect.Value) error {
	switch rv.Type().String() {
	case str:
		return r2s.String(key, dv, rv)
	case strs:
		return r2s.Strings(key, dv, rv)
	case it:
		return r2s.Int(key, dv, rv)
	case it64:
		return r2s.Int(key, dv, rv)
	case flt:
		return r2s.Float(key, dv, rv)
	case bol:
		return r2s.Bool(key, dv, rv)
	default:
		return r2s.Interface(key, dv, rv)
	}

}

func (r2s *read2struct) Int(key, dv string, rv reflect.Value) error {
	var dint int
	var err error
	if dv == "" {
		dint = 0
	} else {
		dint, err = strconv.Atoi(dv)
		if err != nil {
			return err
		}

	}
	v := r2s.cf.DefaultInt(key, dint)
	rv.SetInt(int64(v))
	return nil
}


func (r2s *read2struct) String(key, dv string, rv reflect.Value) error {
	v := r2s.cf.DefaultString(key, dv)
	rv.SetString(v)
	return nil
}

func (r2s *read2struct) Strings(key, dv string, rv reflect.Value) error {
	dvs := strings.Split(dv, ";")
	v := r2s.cf.DefaultStrings(key, dvs)
	rv.Set(reflect.ValueOf(v))
	return nil
}
func (r2s *read2struct) Float(key, dv string, rv reflect.Value) error {
	var df float64
	var err error
	if dv == "" {
		df = float64(0)
	} else {
		df, err = strconv.ParseFloat(dv, 0)
		if err != nil {
			return err
		}

	}
	v := r2s.cf.DefaultFloat(key, df)
	rv.SetFloat(v)
	return nil
}
func (r2s *read2struct) Bool(key, dv string, rv reflect.Value) error {
	b, _ := parseBool(dv)
	v := r2s.cf.DefaultBool(key, b)
	rv.SetBool(v)
	return nil
}

func (r2s *read2struct) Interface(key, dv string, rv reflect.Value) error {
	it, err := r2s.cf.Interface(key)
	if err != nil {
		return nil
	}
	rv.Set(reflect.ValueOf(it))
	return nil
}
