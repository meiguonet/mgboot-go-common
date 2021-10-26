package castx

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ToStringE(arg0 interface{}) (string, error) {
	if arg0 == nil {
		return "", fmt.Errorf("cannot convert %s to string", "nil")
	}

	srcType := fmt.Sprintf("%t", arg0)
	arg0 = indirect(arg0)

	switch v := arg0.(type) {
	case string:
		return v, nil
	case bool:
		return strconv.FormatBool(v), nil
	case uint:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint64:
		return strconv.FormatUint(v, 10), nil
	case int:
		return strconv.FormatInt(int64(v), 10), nil
	case int8:
		return strconv.FormatInt(int64(v), 10), nil
	case int16:
		return strconv.FormatInt(int64(v), 10), nil
	case int32:
		return strconv.FormatInt(int64(v), 10), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case float32:
		return toFloatString(v), nil
	case float64:
		return toFloatString(v), nil
	}

	return "", fmt.Errorf("cannot convert %s to string", srcType)
}

func ToString(arg0 interface{}) string {
	s1, _ := ToStringE(arg0)
	return s1
}

func ToIntE(arg0 interface{}) (int, error) {
	format1 := "cannot convert %s to int"

	if arg0 == nil {
		return 0, fmt.Errorf(format1, "nil")
	}

	srcType := fmt.Sprintf("%t", arg0)
	arg0 = indirect(arg0)

	switch v := arg0.(type) {
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		return int(v), nil
	case int:
		return v, nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case uint64, int64:
		n1, err := strconv.ParseInt(ToString(v), 10, 32)

		if err != nil {
			return 0, nil
		}

		return int(n1), nil
	case float32, float64:
		n1, err := strconv.Atoi(getIntPart(arg0))

		if err != nil {
			return 0, nil
		}

		return n1, nil
	case string:
		n1, err := strconv.Atoi(v)

		if err != nil {
			return 0, err
		}

		return n1, nil
	}

	return 0, fmt.Errorf(format1, srcType)
}

func ToInt(arg0 interface{}, defaultValue ...int) int {
	var _defaultValue int

	if len(defaultValue) > 0 {
		_defaultValue = defaultValue[0]
	}

	n1, err := ToIntE(arg0)

	if err != nil {
		return _defaultValue
	}

	return n1
}

func ToInt8E(arg0 interface{}) (int8, error) {
	format1 := "cannot convert %s to int8"

	if arg0 == nil {
		return 0, fmt.Errorf(format1, "nil")
	}

	srcType := fmt.Sprintf("%t", arg0)
	arg0 = indirect(arg0)

	switch v := arg0.(type) {
	case uint8:
		return int8(v), nil
	case int8:
		return v, nil
	case uint, uint16, uint32, uint64, int, int16, int32, int64:
		n1, err := strconv.ParseInt(ToString(arg0), 10, 8)

		if err != nil {
			return 0, err
		}

		return int8(n1), nil
	case float32, float64:
		n1, err := strconv.ParseInt(getIntPart(arg0), 10, 8)

		if err != nil {
			return 0, err
		}

		return int8(n1), nil
	case string:
		n1, err := strconv.ParseInt(v, 10, 8)

		if err != nil {
			return 0, err
		}

		return int8(n1), nil
	}

	return 0, fmt.Errorf(format1, srcType)
}

func ToInt8(arg0 interface{}, defaultValue ...int8) int8 {
	var _defaultValue int8

	if len(defaultValue) > 0 {
		_defaultValue = defaultValue[0]
	}

	n1, err := ToInt8E(arg0)

	if err != nil {
		return _defaultValue
	}

	return n1
}

func ToInt16E(arg0 interface{}) (int16, error) {
	format1 := "cannot convert %s to int16"

	if arg0 == nil {
		return 0, fmt.Errorf(format1, "nil")
	}

	srcType := fmt.Sprintf("%t", arg0)
	arg0 = indirect(arg0)

	switch v := arg0.(type) {
	case uint8:
		return int16(v), nil
	case uint16:
		return int16(v), nil
	case int8:
		return int16(v), nil
	case int16:
		return v, nil
	case uint, uint32, uint64, int, int32, int64:
		n1, err := strconv.ParseInt(ToString(arg0), 10, 16)

		if err != nil {
			return 0, err
		}

		return int16(n1), nil
	case float32, float64:
		n1, err := strconv.ParseInt(getIntPart(arg0), 10, 16)

		if err != nil {
			return 0, err
		}

		return int16(n1), nil
	case string:
		n1, err := strconv.ParseInt(v, 10, 16)

		if err != nil {
			return 0, err
		}

		return int16(n1), nil
	}

	return 0, fmt.Errorf(format1, srcType)
}

func ToInt16(arg0 interface{}, defaultValue ...int16) int16 {
	var _defaultValue int16

	if len(defaultValue) > 0 {
		_defaultValue = defaultValue[0]
	}

	n1, err := ToInt16E(arg0)

	if err != nil {
		return _defaultValue
	}

	return n1
}

func ToInt32E(arg0 interface{}) (int32, error) {
	format1 := "cannot convert %s to int32"

	if arg0 == nil {
		return 0, fmt.Errorf(format1, "nil")
	}

	srcType := fmt.Sprintf("%t", arg0)
	arg0 = indirect(arg0)

	switch v := arg0.(type) {
	case uint8:
		return int32(v), nil
	case uint16:
		return int32(v), nil
	case uint32:
		return int32(v), nil
	case int8:
		return int32(v), nil
	case int16:
		return int32(v), nil
	case int32:
		return v, nil
	case uint, uint64, int, int64:
		n1, err := strconv.ParseInt(ToString(arg0), 10, 32)

		if err != nil {
			return 0, err
		}

		return int32(n1), nil
	case float32, float64:
		n1, err := strconv.ParseInt(getIntPart(arg0), 10, 32)

		if err != nil {
			return 0, err
		}

		return int32(n1), nil
	case string:
		n1, err := strconv.ParseInt(v, 10, 32)

		if err != nil {
			return 0, err
		}

		return int32(n1), nil
	}

	return 0, fmt.Errorf(format1, srcType)
}

func ToInt32(arg0 interface{}, defaultValue ...int32) int32 {
	var _defaultValue int32

	if len(defaultValue) > 0 {
		_defaultValue = defaultValue[0]
	}

	n1, err := ToInt32E(arg0)

	if err != nil {
		return _defaultValue
	}

	return n1
}

func ToInt64E(arg0 interface{}) (int64, error) {
	format1 := "cannot convert %s to int64"

	if arg0 == nil {
		return 0, fmt.Errorf(format1, "nil")
	}

	srcType := fmt.Sprintf("%t", arg0)
	arg0 = indirect(arg0)

	switch v := arg0.(type) {
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case float32, float64:
		n1, err := strconv.ParseInt(getIntPart(arg0), 10, 64)

		if err != nil {
			return 0, err
		}

		return n1, nil
	case string:
		n1, err := strconv.ParseInt(v, 10, 64)

		if err != nil {
			return 0, err
		}

		return n1, nil
	}

	return 0, fmt.Errorf(format1, srcType)
}

func ToInt64(arg0 interface{}, defaultValue ...int64) int64 {
	var _defaultValue int64

	if len(defaultValue) > 0 {
		_defaultValue = defaultValue[0]
	}

	n1, err := ToInt64E(arg0)

	if err != nil {
		return _defaultValue
	}

	return n1
}

func ToUintE(arg0 interface{}) (uint, error) {
	format1 := "cannot convert %s to uint"

	if arg0 == nil {
		return 0, fmt.Errorf(format1, "nil")
	}

	srcType := fmt.Sprintf("%t", arg0)
	arg0 = indirect(arg0)

	switch v := arg0.(type) {
	case uint:
		return v, nil
	case uint8:
		return uint(v), nil
	case uint16:
		return uint(v), nil
	case uint32:
		return uint(v), nil
	case uint64, int, int8, int16, int32, int64:
		n1, err := strconv.ParseUint(ToString(v), 10, 32)

		if err != nil {
			return 0, nil
		}

		return uint(n1), nil
	case float32, float64:
		n1, err := strconv.ParseUint(getIntPart(arg0), 10, 32)

		if err != nil {
			return 0, nil
		}

		return uint(n1), nil
	case string:
		n1, err := strconv.ParseUint(v, 10, 32)

		if err != nil {
			return 0, nil
		}

		return uint(n1), nil
	}

	return 0, fmt.Errorf(format1, srcType)
}

func ToUint(arg0 interface{}, defaultValue ...uint) uint {
	var _defaultValue uint

	if len(defaultValue) > 0 {
		_defaultValue = defaultValue[0]
	}

	n1, err := ToUintE(arg0)

	if err != nil {
		return _defaultValue
	}

	return n1
}

func ToUint8E(arg0 interface{}) (uint8, error) {
	format1 := "cannot convert %s to uint8"

	if arg0 == nil {
		return 0, fmt.Errorf(format1, "nil")
	}

	srcType := fmt.Sprintf("%t", arg0)
	arg0 = indirect(arg0)

	switch v := arg0.(type) {
	case uint8:
		return v, nil
	case uint, uint16, uint32, uint64, int, int8, int16, int32, int64:
		n1, err := strconv.ParseUint(ToString(v), 10, 8)

		if err != nil {
			return 0, nil
		}

		return uint8(n1), nil
	case float32, float64:
		n1, err := strconv.ParseUint(getIntPart(arg0), 10, 8)

		if err != nil {
			return 0, nil
		}

		return uint8(n1), nil
	case string:
		n1, err := strconv.ParseUint(v, 10, 8)

		if err != nil {
			return 0, nil
		}

		return uint8(n1), nil
	}

	return 0, fmt.Errorf(format1, srcType)
}

func ToUint8(arg0 interface{}, defaultValue ...uint8) uint8 {
	var _defaultValue uint8

	if len(defaultValue) > 0 {
		_defaultValue = defaultValue[0]
	}

	n1, err := ToUint8E(arg0)

	if err != nil {
		return _defaultValue
	}

	return n1
}

func ToUint16E(arg0 interface{}) (uint16, error) {
	format1 := "cannot convert %s to uint16"

	if arg0 == nil {
		return 0, fmt.Errorf(format1, "nil")
	}

	srcType := fmt.Sprintf("%t", arg0)
	arg0 = indirect(arg0)

	switch v := arg0.(type) {
	case uint8:
		return uint16(v), nil
	case uint16:
		return v, nil
	case uint, uint32, uint64, int, int8, int16, int32, int64:
		n1, err := strconv.ParseUint(ToString(v), 10, 16)

		if err != nil {
			return 0, nil
		}

		return uint16(n1), nil
	case float32, float64:
		n1, err := strconv.ParseUint(getIntPart(arg0), 10, 16)

		if err != nil {
			return 0, nil
		}

		return uint16(n1), nil
	case string:
		n1, err := strconv.ParseUint(v, 10, 16)

		if err != nil {
			return 0, nil
		}

		return uint16(n1), nil
	}

	return 0, fmt.Errorf(format1, srcType)
}

func ToUint16(arg0 interface{}, defaultValue ...uint16) uint16 {
	var _defaultValue uint16

	if len(defaultValue) > 0 {
		_defaultValue = defaultValue[0]
	}

	n1, err := ToUint16E(arg0)

	if err != nil {
		return _defaultValue
	}

	return n1
}

func ToUint32E(arg0 interface{}) (uint32, error) {
	format1 := "cannot convert %s to uint32"

	if arg0 == nil {
		return 0, fmt.Errorf(format1, "nil")
	}

	srcType := fmt.Sprintf("%t", arg0)
	arg0 = indirect(arg0)

	switch v := arg0.(type) {
	case uint8:
		return uint32(v), nil
	case uint16:
		return uint32(v), nil
	case uint32:
		return v, nil
	case uint, uint64, int, int8, int16, int32, int64:
		n1, err := strconv.ParseUint(ToString(v), 10, 32)

		if err != nil {
			return 0, nil
		}

		return uint32(n1), nil
	case float32, float64:
		n1, err := strconv.ParseUint(getIntPart(arg0), 10, 32)

		if err != nil {
			return 0, nil
		}

		return uint32(n1), nil
	case string:
		n1, err := strconv.ParseUint(v, 10, 32)

		if err != nil {
			return 0, nil
		}

		return uint32(n1), nil
	}

	return 0, fmt.Errorf(format1, srcType)
}

func ToUint32(arg0 interface{}, defaultValue ...uint32) uint32 {
	var _defaultValue uint32

	if len(defaultValue) > 0 {
		_defaultValue = defaultValue[0]
	}

	n1, err := ToUint32E(arg0)

	if err != nil {
		return _defaultValue
	}

	return n1
}

func ToUint64E(arg0 interface{}) (uint64, error) {
	format1 := "cannot convert %s to uint64"

	if arg0 == nil {
		return 0, fmt.Errorf(format1, "nil")
	}

	srcType := fmt.Sprintf("%t", arg0)
	arg0 = indirect(arg0)

	switch v := arg0.(type) {
	case uint:
		return uint64(v), nil
	case uint8:
		return uint64(v), nil
	case uint16:
		return uint64(v), nil
	case uint32:
		return uint64(v), nil
	case uint64:
		return v, nil
	case int, int8, int16, int32, int64:
		n1, err := strconv.ParseUint(ToString(v), 10, 64)

		if err != nil {
			return 0, nil
		}

		return n1, nil
	case float32, float64:
		n1, err := strconv.ParseUint(getIntPart(arg0), 10, 64)

		if err != nil {
			return 0, nil
		}

		return n1, nil
	case string:
		n1, err := strconv.ParseUint(v, 10, 64)

		if err != nil {
			return 0, nil
		}

		return n1, nil
	}

	return 0, fmt.Errorf(format1, srcType)
}

func ToUint64(arg0 interface{}, defaultValue ...uint64) uint64 {
	var _defaultValue uint64

	if len(defaultValue) > 0 {
		_defaultValue = defaultValue[0]
	}

	n1, err := ToUint64E(arg0)

	if err != nil {
		return _defaultValue
	}

	return n1
}

func ToFloat32E(arg0 interface{}) (float32, error) {
	format1 := "cannot convert %s to float32"

	if arg0 == nil {
		return 0, fmt.Errorf(format1, "nil")
	}

	srcType := fmt.Sprintf("%t", arg0)
	arg0 = indirect(arg0)

	switch v := arg0.(type) {
	case uint8:
		return float32(v), nil
	case uint16:
		return float32(v), nil
	case uint32:
		return float32(v), nil
	case int8:
		return float32(v), nil
	case int16:
		return float32(v), nil
	case int32:
		return float32(v), nil
	case uint, uint64, int, int64:
		n1, err := strconv.ParseFloat(toFloatString(float64(ToInt64(arg0))), 32)

		if err != nil {
			return 0, nil
		}

		return float32(n1), nil
	case float32:
		return v, nil
	case float64:
		n1, err := strconv.ParseFloat(toFloatString(v), 32)

		if err != nil {
			return 0, nil
		}

		return float32(n1), nil
	case string:
		n1, err := strconv.ParseFloat(v, 32)

		if err != nil {
			return 0, nil
		}

		return float32(n1), nil
	}

	return 0, fmt.Errorf(format1, srcType)
}

func ToFloat32(arg0 interface{}, defaultValue ...float32) float32 {
	var _defaultValue float32

	if len(defaultValue) > 0 {
		_defaultValue = defaultValue[0]
	}

	n1, err := ToFloat32E(arg0)

	if err != nil {
		return _defaultValue
	}

	return n1
}

func ToFloat64E(arg0 interface{}) (float64, error) {
	format1 := "cannot convert %s to float64"

	if arg0 == nil {
		return 0, fmt.Errorf(format1, "nil")
	}

	srcType := fmt.Sprintf("%t", arg0)
	arg0 = indirect(arg0)

	switch v := arg0.(type) {
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case string:
		n1, err := strconv.ParseFloat(v, 64)

		if err != nil {
			return 0, nil
		}

		return n1, nil
	}

	return 0, fmt.Errorf(format1, srcType)
}

func ToFloat64(arg0 interface{}, defaultValue ...float64) float64 {
	var _defaultValue float64

	if len(defaultValue) > 0 {
		_defaultValue = defaultValue[0]
	}

	n1, err := ToFloat64E(arg0)

	if err != nil {
		return _defaultValue
	}

	return n1
}

func ToStringMapE(arg0 interface{}) (map[string]interface{}, error) {
	format1 := "cannot convert %s to map[string]interface{}"

	if arg0 == nil {
		return nil, fmt.Errorf(format1, "nil")
	}

	if map1, ok := arg0.(map[string]interface{}); ok {
		return map1, nil
	}

	srcType := fmt.Sprintf("%t", arg0)
	rt := reflect.TypeOf(arg0)
	rv := reflect.ValueOf(arg0)

	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rv = rv.Elem()
	}

	if rt.Kind() != reflect.Map {
		return nil, fmt.Errorf(format1, srcType)
	}

	map1 := map[string]interface{}{}
	iter := rv.MapRange()

	for iter.Next() {
		key, ok := iter.Key().Interface().(string)

		if !ok || key == "" {
			continue
		}

		map1[key] = iter.Value().Interface()
	}

	return map1, nil
}

func ToStringMap(arg0 interface{}) map[string]interface{} {
	map1, _ := ToStringMapE(arg0)

	if len(map1) < 1 {
		return map[string]interface{}{}
	}

	return map1
}

func ToStringMapStringE(arg0 interface{}) (map[string]string, error) {
	format1 := "cannot convert %s to map[string]string"

	if arg0 == nil {
		return nil, fmt.Errorf(format1, "nil")
	}

	if map1, ok := arg0.(map[string]string); ok {
		return map1, nil
	}

	srcType := fmt.Sprintf("%t", arg0)
	rt := reflect.TypeOf(arg0)
	rv := reflect.ValueOf(arg0)

	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rv = rv.Elem()
	}

	if rt.Kind() != reflect.Map {
		return nil, fmt.Errorf(format1, srcType)
	}

	map1 := map[string]string{}
	iter := rv.MapRange()

	for iter.Next() {
		key, ok := iter.Key().Interface().(string)

		if !ok || key == "" {
			continue
		}

		if value, ok := iter.Value().Interface().(string); ok {
			map1[key] = value
			continue
		}

		value, err := ToStringE(iter.Value().Interface())

		if err != nil {
			return nil, fmt.Errorf("fail to convert %s to map[string]string", srcType)
		}

		map1[key] = value
	}

	return map1, nil
}

func ToStringMapString(arg0 interface{}) map[string]string {
	map1, _ := ToStringMapStringE(arg0)

	if len(map1) < 1 {
		return map[string]string{}
	}

	return map1
}

func ToSliceE(arg0 interface{}) ([]interface{}, error) {
	format1 := "cannot convert %s to []interface{}"

	if arg0 == nil {
		return nil, fmt.Errorf(format1, "nil")
	}

	if list, ok := arg0.([]interface{}); ok {
		return list, nil
	}

	srcType := fmt.Sprintf("%t", arg0)
	rt := reflect.TypeOf(arg0)
	rv := reflect.ValueOf(arg0)

	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rv = rv.Elem()
	}

	if rt.Kind() != reflect.Slice {
		return nil, fmt.Errorf(format1, srcType)
	}

	list := make([]interface{}, 0)

	if rv.Len() < 1 {
		return list, nil
	}

	for i := 0; i < rv.Len(); i++ {
		list = append(list, rv.Index(i).Interface())
	}

	return list, nil
}

func ToSlice(arg0 interface{}) []interface{} {
	list, _ := ToSliceE(arg0)

	if len(list) < 1 {
		return []interface{}{}
	}

	return list
}

func ToStringSliceE(arg0 interface{}) ([]string, error) {
	format1 := "cannot convert %s to []string"

	if arg0 == nil {
		return nil, fmt.Errorf(format1, "nil")
	}

	if list, ok := arg0.([]string); ok {
		return list, nil
	}

	srcType := fmt.Sprintf("%t", arg0)
	rt := reflect.TypeOf(arg0)
	rv := reflect.ValueOf(arg0)

	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rv = rv.Elem()
	}

	if rt.Kind() != reflect.Slice {
		return nil, fmt.Errorf(format1, srcType)
	}

	list := make([]string, 0)

	if rv.Len() < 1 {
		return list, nil
	}

	for i := 0; i < rv.Len(); i++ {
		if value, ok := rv.Index(i).Interface().(string); ok {
			list = append(list, value)
			continue
		}

		value, err := ToStringE(rv.Index(i).Interface())

		if err != nil {
			return nil, fmt.Errorf("fail to convert %s to []string", srcType)
		}

		list = append(list, value)
	}

	return list, nil
}

func ToStringSlice(arg0 interface{}) []string {
	list, _ := ToStringSliceE(arg0)

	if len(list) < 1 {
		return []string{}
	}

	return list
}

func ToIntSliceE(arg0 interface{}) ([]int, error) {
	format1 := "cannot convert %s to []int"

	if arg0 == nil {
		return nil, fmt.Errorf(format1, "nil")
	}

	if list, ok := arg0.([]int); ok {
		return list, nil
	}

	srcType := fmt.Sprintf("%t", arg0)
	rt := reflect.TypeOf(arg0)
	rv := reflect.ValueOf(arg0)

	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rv = rv.Elem()
	}

	if rt.Kind() != reflect.Slice {
		return nil, fmt.Errorf(format1, srcType)
	}

	list := make([]int, 0)

	if rv.Len() < 1 {
		return list, nil
	}

	for i := 0; i < rv.Len(); i++ {
		if value, ok := rv.Index(i).Interface().(int); ok {
			list = append(list, value)
			continue
		}

		value, err := ToIntE(rv.Index(i).Interface())

		if err != nil {
			return nil, fmt.Errorf("fail to convert %s to []int", srcType)
		}

		list = append(list, value)
	}

	return list, nil
}

func ToIntSlice(arg0 interface{}) []int {
	list, _ := ToIntSliceE(arg0)

	if len(list) < 1 {
		return []int{}
	}

	return list
}

func ToMapSliceE(arg0 interface{}) ([]map[string]interface{}, error) {
	format1 := "cannot convert %s to []map[string]interface{}"

	if arg0 == nil {
		return nil, fmt.Errorf(format1, "nil")
	}

	if list, ok := arg0.([]map[string]interface{}); ok {
		return list, nil
	}

	srcType := fmt.Sprintf("%t", arg0)
	rt := reflect.TypeOf(arg0)
	rv := reflect.ValueOf(arg0)

	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rv = rv.Elem()
	}

	if rt.Kind() != reflect.Slice {
		return nil, fmt.Errorf(format1, srcType)
	}

	list := make([]map[string]interface{}, 0)

	if rv.Len() < 1 {
		return list, nil
	}

	for i := 0; i < rv.Len(); i++ {
		if value, ok := rv.Index(i).Interface().(map[string]interface{}); ok {
			list = append(list, value)
			continue
		}

		value, err := ToStringMapE(rv.Index(i).Interface())

		if err != nil {
			return nil, fmt.Errorf("fail to convert %s to []map[string]interface{}", srcType)
		}

		list = append(list, value)
	}

	return list, nil
}

func ToMapSlice(arg0 interface{}) []map[string]interface{} {
	list, _ := ToMapSliceE(arg0)

	if len(list) < 1 {
		return []map[string]interface{}{}
	}

	return list
}

func ToBoolE(arg0 interface{}) (bool, error) {
	format1 := "cannot convert %s to bool"
	
	if arg0 == nil {
		return false, fmt.Errorf(format1, "nil")
	}

	srcType := fmt.Sprintf("%t", arg0)

	switch v := arg0.(type) {
	case bool:
		return v, nil
	case int:
		if v == 1 {
			return true, nil
		}
		
		return false, nil
	case string:
		return strconv.ParseBool(v)
	}

	return false, fmt.Errorf(format1, srcType)
}

func ToBool(arg0 interface{}) bool {
	b1, _ := ToBoolE(arg0)
	return b1
}

func ToDurationE(arg0 interface{}) (time.Duration, error) {
	format1 := "cannot convert %s to time.Duration"

	if arg0 == nil {
		return 0, fmt.Errorf(format1, "nil")
	}

	srcType := fmt.Sprintf("%t", arg0)

	switch v := arg0.(type) {
	case time.Duration:
		return v, nil
	case *time.Duration:
		return ToDurationE(indirect(arg0))
	case int:
		return time.Duration(int64(v)) * time.Second, nil
	case int64:
		return time.Duration(v) * time.Millisecond, nil
	case string:
		return time.ParseDuration(v)
	}

	return 0, fmt.Errorf(format1, srcType)
}

func ToDuration(arg0 interface{}) time.Duration {
	d1, _ := ToDurationE(arg0)
	return d1
}

func ToDataSizeE(arg0 interface{}) (int64, error) {
	format1 := "cannot convert %s to data size"

	if arg0 == nil {
		return 0, fmt.Errorf(format1, "nil")
	}

	srcType := fmt.Sprintf("%t", arg0)
	str, err := ToStringE(arg0)

	if err != nil {
		return 0, fmt.Errorf(format1, srcType)
	}

	if str == "" {
		return 0, nil
	}

	regex1 := regexp.MustCompile("[^0-9]")
	digits := regex1.ReplaceAllString(str, "")

	if digits == "" {
		return 0, nil
	}

	num := ToInt64(digits)
	regex2 := regexp.MustCompile("[^A-Za-z]")
	unit := strings.ToLower(regex2.ReplaceAllString(str, ""))

	switch unit {
	case "":
		return num, nil
	case "k":
		return num * 1024, nil
	case "m":
		return num * 1024 * 1024, nil
	case "g":
		return num * 1024 * 1024 * 1024, nil
	default:
		return 0, nil
	}
}

func ToDataSize(arg0 interface{}) int64 {
	n1, _ := ToDataSizeE(arg0)
	return n1
}

func toFloatString(arg0 interface{}) string {
	parts := strings.Split(fmt.Sprintf("%0.12f", arg0), ".")
	p2 := strings.TrimSuffix(parts[1], "0")

	if p2 == "" {
		return parts[0]
	}

	return parts[0] + "." + p2
}

func getIntPart(arg0 interface{}) string {
	parts := strings.Split(fmt.Sprintf("%0.2f", arg0), ".")
	return parts[0]
}

func indirect(arg0 interface{}) interface{} {
	if arg0 == nil {
		return nil
	}

	if rt := reflect.TypeOf(arg0); rt.Kind() != reflect.Ptr {
		return arg0
	}

	rv := reflect.ValueOf(arg0)

	for rv.Kind() == reflect.Ptr && !rv.IsNil() {
		rv = rv.Elem()
	}

	return rv.Interface()
}
