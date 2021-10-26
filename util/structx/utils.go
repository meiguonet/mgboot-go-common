package structx

import (
	"reflect"
	"strings"
)

type toMapOption interface {
	getToMapOptionName() string
}

type fieldNamesOption struct {
	mode string
	fieldNames []string
}

func (opt *fieldNamesOption) getToMapOptionName() string {
	return "FieldNamesOption"
}

type ignoreZeroValueOption struct {
}

func (opt *ignoreZeroValueOption) getToMapOptionName() string {
	return "IgnoreZeroValueOption"
}

type ignoreNilValueOption struct {
}

func (opt *ignoreNilValueOption) getToMapOptionName() string {
	return "IgnoreNilValueOption"
}

func ToMapIncludeFields(fields []string) *fieldNamesOption {
	return &fieldNamesOption{mode: "include", fieldNames: fields}
}

func ToMapExcludeFields(fields []string) *fieldNamesOption {
	return &fieldNamesOption{mode: "exclude", fieldNames: fields}
}

func ToMapIgnoreZero() *ignoreZeroValueOption {
	return &ignoreZeroValueOption{}
}

func ToMapIgnoreNil() *ignoreNilValueOption {
	return &ignoreNilValueOption{}
}

func ToMap(arg0 interface{}, opts ...toMapOption) map[string]interface{} {
	if arg0 == nil {
		return map[string]interface{}{}
	}

	rt := reflect.TypeOf(arg0)
	var isPtr bool

	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		isPtr = true
	}

	if rt.Kind() != reflect.Struct {
		return map[string]interface{}{}
	}

	rv := reflect.ValueOf(arg0)

	if isPtr {
		rv = rv.Elem()
	}

	var opt1 *fieldNamesOption
	var ignoreZero bool
	var ignoreNil bool

	for _, opt := range opts {
		if _opt, ok := opt.(*fieldNamesOption); ok {
			if opt1 == nil {
				opt1 = _opt
			}

			continue
		}

		if _, ok := opt.(*ignoreZeroValueOption); ok {
			ignoreZero = true
			continue
		}

		if _, ok := opt.(*ignoreNilValueOption); ok {
			ignoreNil = true
		}
	}

	map1 := make(map[string]interface{})

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fieldValue := rv.Field(i)

		if ignoreNil && field.Type.Kind() == reflect.Ptr && fieldValue.IsZero() {
			continue
		}

		if ignoreZero && fieldValue.IsZero() {
			continue
		}

		if !inIncludeFields(field.Name, opt1) {
			continue
		}

		if inExcludeFields(field.Name, opt1) {
			continue
		}

		map1[getMapKey(field)] = fieldValue.Interface()
	}

	return map1
}

func inIncludeFields(fieldName string, opt *fieldNamesOption) bool {
	if opt == nil || opt.mode != "include" || len(opt.fieldNames) < 1 {
		return true
	}

	fieldName = strings.ToLower(fieldName)

	for _, fname := range opt.fieldNames {
		if fieldName == strings.ToLower(fname) {
			return true
		}
	}

	return false
}

func inExcludeFields(fieldName string, opt *fieldNamesOption) bool {
	if opt == nil || opt.mode != "exclude" || len(opt.fieldNames) < 1 {
		return false
	}

	fieldName = strings.ToLower(fieldName)

	for _, fname := range opt.fieldNames {
		if fieldName == strings.ToLower(fname) {
			return true
		}
	}

	return false
}

func getMapKey(field reflect.StructField) string {
	key := strings.TrimSpace(field.Tag.Get("mapkey"))

	if key != "" {
		return key
	}

	return strings.ToLower(field.Name[:1]) + field.Name[1:]
}
