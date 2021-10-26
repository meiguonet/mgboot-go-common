package reflectx

import (
	"reflect"
	"strings"
	"time"
	"unsafe"
)

func Indirect(arg0 interface{}) interface{} {
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

func Ptr(arg0 interface{}) unsafe.Pointer {
	return unsafe.Pointer(reflect.ValueOf(arg0).Pointer())
}

func SetStructFieldValue(structPtr unsafe.Pointer, field reflect.StructField, value interface{}) bool {
	var success bool

	defer func() {
		if err := recover(); err != nil {
			success = false
		}
	}()

	fieldPtr := structFieldPtr(structPtr, field)
	value = Indirect(value)

	switch v := value.(type) {
	case string:
		if field.Type.Kind() == reflect.String {
			*((*string)(unsafe.Pointer(fieldPtr))) = v
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.String {
			*((**string)(unsafe.Pointer(fieldPtr))) = &v
		}

		success = true
	case bool:
		if field.Type.Kind() == reflect.Bool {
			*((*bool)(unsafe.Pointer(fieldPtr))) = v
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Bool {
			*((**bool)(unsafe.Pointer(fieldPtr))) = &v
			success = true
		}
	case uint8:
		if field.Type.Kind() == reflect.Uint8 {
			*((*uint8)(unsafe.Pointer(fieldPtr))) = v
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Uint8 {
			*((**uint8)(unsafe.Pointer(fieldPtr))) = &v
			success = true
		} else if field.Type.Kind() == reflect.Uint16 {
			*((*uint16)(unsafe.Pointer(fieldPtr))) = uint16(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Uint16 {
			n1 := uint16(v)
			*((**uint16)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		} else if field.Type.Kind() == reflect.Uint32 {
			*((*uint32)(unsafe.Pointer(fieldPtr))) = uint32(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Uint32 {
			n1 := uint32(v)
			*((**uint32)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		} else if field.Type.Kind() == reflect.Uint64 {
			*((*uint64)(unsafe.Pointer(fieldPtr))) = uint64(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Uint64 {
			n1 := uint64(v)
			*((**uint64)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		} else if field.Type.Kind() == reflect.Uint {
			*((*uint)(unsafe.Pointer(fieldPtr))) = uint(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Uint {
			n1 := uint(v)
			*((**uint)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		} else if field.Type.Kind() == reflect.Int8 {
			*((*int8)(unsafe.Pointer(fieldPtr))) = int8(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int8 {
			n1 := int8(v)
			*((**int8)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		} else if field.Type.Kind() == reflect.Int16 {
			*((*int16)(unsafe.Pointer(fieldPtr))) = int16(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int16 {
			n1 := int16(v)
			*((**int16)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		} else if field.Type.Kind() == reflect.Int32 {
			*((*int32)(unsafe.Pointer(fieldPtr))) = int32(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int32 {
			n1 := int32(v)
			*((**int32)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		} else if field.Type.Kind() == reflect.Int64 {
			*((*int64)(unsafe.Pointer(fieldPtr))) = int64(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int64 {
			n1 := int64(v)
			*((**int64)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		} else if field.Type.Kind() == reflect.Int {
			*((*int)(unsafe.Pointer(fieldPtr))) = int(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int {
			n1 := int(v)
			*((**int)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		}
	case uint16:
		if field.Type.Kind() == reflect.Uint16 {
			*((*uint16)(unsafe.Pointer(fieldPtr))) = v
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Uint16 {
			*((**uint16)(unsafe.Pointer(fieldPtr))) = &v
			success = true
		} else if field.Type.Kind() == reflect.Uint32 {
			*((*uint32)(unsafe.Pointer(fieldPtr))) = uint32(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Uint32 {
			n1 := uint32(v)
			*((**uint32)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		} else if field.Type.Kind() == reflect.Uint64 {
			*((*uint64)(unsafe.Pointer(fieldPtr))) = uint64(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Uint64 {
			n1 := uint64(v)
			*((**uint64)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		} else if field.Type.Kind() == reflect.Uint {
			*((*uint)(unsafe.Pointer(fieldPtr))) = uint(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Uint {
			n1 := uint(v)
			*((**uint)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		} else if field.Type.Kind() == reflect.Int16 {
			*((*int16)(unsafe.Pointer(fieldPtr))) = int16(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int16 {
			n1 := int16(v)
			*((**int16)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		} else if field.Type.Kind() == reflect.Int32 {
			*((*int32)(unsafe.Pointer(fieldPtr))) = int32(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int32 {
			n1 := int32(v)
			*((**int32)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		} else if field.Type.Kind() == reflect.Int64 {
			*((*int64)(unsafe.Pointer(fieldPtr))) = int64(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int64 {
			n1 := int64(v)
			*((**int64)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		} else if field.Type.Kind() == reflect.Int {
			*((*int)(unsafe.Pointer(fieldPtr))) = int(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int {
			n1 := int(v)
			*((**int)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		}
	case uint32:
		if field.Type.Kind() == reflect.Uint32 {
			*((*uint32)(unsafe.Pointer(fieldPtr))) = v
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Uint32 {
			*((**uint32)(unsafe.Pointer(fieldPtr))) = &v
			success = true
		} else if field.Type.Kind() == reflect.Uint64 {
			*((*uint64)(unsafe.Pointer(fieldPtr))) = uint64(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Uint64 {
			n1 := uint64(v)
			*((**uint64)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		} else if field.Type.Kind() == reflect.Uint {
			*((*uint)(unsafe.Pointer(fieldPtr))) = uint(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Uint {
			n1 := uint(v)
			*((**uint)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		} else if field.Type.Kind() == reflect.Int32 {
			*((*int32)(unsafe.Pointer(fieldPtr))) = int32(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int32 {
			n1 := int32(v)
			*((**int32)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		} else if field.Type.Kind() == reflect.Int64 {
			*((*int64)(unsafe.Pointer(fieldPtr))) = int64(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int64 {
			n1 := int64(v)
			*((**int64)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		} else if field.Type.Kind() == reflect.Int {
			*((*int)(unsafe.Pointer(fieldPtr))) = int(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int {
			n1 := int(v)
			*((**int)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		}
	case uint:
		if field.Type.Kind() == reflect.Uint {
			*((*uint)(unsafe.Pointer(fieldPtr))) = v
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Uint {
			*((**uint)(unsafe.Pointer(fieldPtr))) = &v
			success = true
		} else if field.Type.Kind() == reflect.Uint64 {
			*((*uint64)(unsafe.Pointer(fieldPtr))) = uint64(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Uint64 {
			n1 := uint64(v)
			*((**uint64)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		} else if field.Type.Kind() == reflect.Int64 {
			*((*int64)(unsafe.Pointer(fieldPtr))) = int64(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int64 {
			n1 := int64(v)
			*((**int64)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		} else if field.Type.Kind() == reflect.Int {
			*((*int)(unsafe.Pointer(fieldPtr))) = int(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int {
			n1 := int(v)
			*((**int)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		}
	case uint64:
		if field.Type.Kind() == reflect.Uint64 {
			*((*uint64)(unsafe.Pointer(fieldPtr))) = v
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Uint64 {
			*((**uint64)(unsafe.Pointer(fieldPtr))) = &v
			success = true
		} else if field.Type.Kind() == reflect.Int64 {
			*((*int64)(unsafe.Pointer(fieldPtr))) = int64(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int64 {
			n1 := int64(v)
			*((**int64)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		}
	case int8:
		if field.Type.Kind() == reflect.Int8 {
			*((*int8)(unsafe.Pointer(fieldPtr))) = v
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int8 {
			*((**int8)(unsafe.Pointer(fieldPtr))) = &v
			success = true
		} else if field.Type.Kind() == reflect.Int16 {
			*((*int16)(unsafe.Pointer(fieldPtr))) = int16(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int16 {
			n1 := int16(v)
			*((**int16)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		} else if field.Type.Kind() == reflect.Int32 {
			*((*int32)(unsafe.Pointer(fieldPtr))) = int32(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int32 {
			n1 := int32(v)
			*((**int32)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		} else if field.Type.Kind() == reflect.Int64 {
			*((*int64)(unsafe.Pointer(fieldPtr))) = int64(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int64 {
			n1 := int64(v)
			*((**int64)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		} else if field.Type.Kind() == reflect.Int {
			*((*int)(unsafe.Pointer(fieldPtr))) = int(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int {
			n1 := int(v)
			*((**int)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		}
	case int16:
		if field.Type.Kind() == reflect.Int16 {
			*((*int16)(unsafe.Pointer(fieldPtr))) = v
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int16 {
			*((**int16)(unsafe.Pointer(fieldPtr))) = &v
			success = true
		} else if field.Type.Kind() == reflect.Int32 {
			*((*int32)(unsafe.Pointer(fieldPtr))) = int32(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int32 {
			n1 := int32(v)
			*((**int32)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		} else if field.Type.Kind() == reflect.Int64 {
			*((*int64)(unsafe.Pointer(fieldPtr))) = int64(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int64 {
			n1 := int64(v)
			*((**int64)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		} else if field.Type.Kind() == reflect.Int {
			*((*int)(unsafe.Pointer(fieldPtr))) = int(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int {
			n1 := int(v)
			*((**int)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		}
	case int32:
		if field.Type.Kind() == reflect.Int32 {
			*((*int32)(unsafe.Pointer(fieldPtr))) = v
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int32 {
			*((**int32)(unsafe.Pointer(fieldPtr))) = &v
			success = true
		} else if field.Type.Kind() == reflect.Int64 {
			*((*int64)(unsafe.Pointer(fieldPtr))) = int64(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int64 {
			n1 := int64(v)
			*((**int64)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		} else if field.Type.Kind() == reflect.Int {
			*((*int)(unsafe.Pointer(fieldPtr))) = int(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int {
			n1 := int(v)
			*((**int)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		}
	case int:
		if field.Type.Kind() == reflect.Int {
			*((*int)(unsafe.Pointer(fieldPtr))) = v
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int {
			*((**int)(unsafe.Pointer(fieldPtr))) = &v
			success = true
		} else if field.Type.Kind() == reflect.Int64 {
			*((*int64)(unsafe.Pointer(fieldPtr))) = int64(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int64 {
			n1 := int64(v)
			*((**int64)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		}
	case int64:
		if field.Type.Kind() == reflect.Int64 {
			*((*int64)(unsafe.Pointer(fieldPtr))) = v
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int64 {
			*((**int64)(unsafe.Pointer(fieldPtr))) = &v
			success = true
		}
	case float32:
		if field.Type.Kind() == reflect.Float32 {
			*((*float32)(unsafe.Pointer(fieldPtr))) = v
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Float32 {
			*((**float32)(unsafe.Pointer(fieldPtr))) = &v
			success = true
		} else if field.Type.Kind() == reflect.Float64 {
			*((*float64)(unsafe.Pointer(fieldPtr))) = float64(v)
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Float64 {
			n1 := float64(v)
			*((**float64)(unsafe.Pointer(fieldPtr))) = &n1
			success = true
		}
	case float64:
		if field.Type.Kind() == reflect.Float64 {
			*((*float64)(unsafe.Pointer(fieldPtr))) = v
			success = true
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Float64 {
			*((**float64)(unsafe.Pointer(fieldPtr))) = &v
			success = true
		}
	case time.Time:
		if isTimeField(field) {
			*((*time.Time)(unsafe.Pointer(fieldPtr))) = v
			success = true
		} else if isTimePtrField(field) {
			*((**time.Time)(unsafe.Pointer(fieldPtr))) = &v
			success = true
		}
	}

	return success
}

func structFieldPtr(structPtr unsafe.Pointer, field reflect.StructField) uintptr {
	return uintptr(structPtr) + field.Offset
}

func isTimeField(field reflect.StructField) bool {
	k := field.Type.Kind()
	tn := field.Type.String()
	return k == reflect.Struct && strings.Contains(tn, "time.Time")
}

func isTimePtrField(field reflect.StructField) bool {
	if field.Type.Kind() != reflect.Ptr {
		return false
	}

	k := field.Type.Elem().Kind()
	tn := field.Type.Elem().String()
	return k == reflect.Struct && strings.Contains(tn, "time.Time")
}
