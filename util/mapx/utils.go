package mapx

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/meiguonet/mgboot-go-common/enum/RegexConst"
	"github.com/meiguonet/mgboot-go-common/util/castx"
	"github.com/meiguonet/mgboot-go-common/util/reflectx"
	"reflect"
	"regexp"
	"strings"
)


func NewXmlKeyTagMappingOption(mappings [][2]string) *xmlKeyTagOption {
	return &xmlKeyTagOption{mappings: mappings}
}

// @var string[]|string keys
func NewXmlIncludeKeysOption(keys interface{}) *xmlKeysOption {
	includedKeys := make([]string, 0)

	if _keys, ok := keys.([]string); ok && len(_keys) > 0 {
		includedKeys = _keys
	} else if _keys, ok := keys.(string); ok && _keys != "" {
		re := regexp.MustCompile(RegexConst.CommaSep)
		includedKeys = re.Split(_keys, -1)
	}

	return &xmlKeysOption{mode: "include", keys: includedKeys}
}

// @var string[]|string keys
func NewXmlExcludeKeysOption(keys interface{}) *xmlKeysOption {
	excludedKeys := make([]string, 0)

	if _keys, ok := keys.([]string); ok && len(_keys) > 0 {
		excludedKeys = _keys
	} else if _keys, ok := keys.(string); ok && _keys != "" {
		re := regexp.MustCompile(RegexConst.CommaSep)
		excludedKeys = re.Split(_keys, -1)
	}

	return &xmlKeysOption{mode: "exclude", keys: excludedKeys}
}

// @var string[]|string keys
func NewXmlCDataKeysOption(keys interface{}) *xmlCDataOption {
	cdataKeys := make([]string, 0)

	if _keys, ok := keys.([]string); ok && len(_keys) > 0 {
		cdataKeys = _keys
	} else if _keys, ok := keys.(string); ok && _keys != "" {
		re := regexp.MustCompile(RegexConst.CommaSep)
		cdataKeys = re.Split(_keys, -1)
	}

	return &xmlCDataOption{keys: cdataKeys}
}

func NewStructKeyFieldOption(mappings [][2]string) *structKeyFieldOption {
	return &structKeyFieldOption{mappings: mappings}
}

// @var string[]|string keys
func NewSturctIncludeKeysOption(keys interface{}) *structKeysOption {
	includedKeys := make([]string, 0)

	if _keys, ok := keys.([]string); ok && len(_keys) > 0 {
		includedKeys = _keys
	} else if _keys, ok := keys.(string); ok && _keys != "" {
		re := regexp.MustCompile(RegexConst.CommaSep)
		includedKeys = re.Split(_keys, -1)
	}

	return &structKeysOption{mode: "include", keys: includedKeys}
}

// @var string[]|string keys
func NewSturctExcludeKeysOption(keys interface{}) *structKeysOption {
	excludedKeys := make([]string, 0)

	if _keys, ok := keys.([]string); ok && len(_keys) > 0 {
		excludedKeys = _keys
	} else if _keys, ok := keys.(string); ok && _keys != "" {
		re := regexp.MustCompile(RegexConst.CommaSep)
		excludedKeys = re.Split(_keys, -1)
	}

	return &structKeysOption{mode: "exclude", keys: excludedKeys}
}

func FromXml(buf []byte) map[string]string {
	var map1 map[string]string

	if xml.Unmarshal(buf, (*stringMap)(&map1)) != nil {
		return map[string]string{}
	}

	return map1
}

func ToXml(map1 map[string]interface{}, opts ...xmlOption) string {
	if len(map1) < 1 {
		return ""
	}

	keyTagEntries := make([]xmlKeyTagEntry, 0)
	includedKeys := make([]string, 0)
	excludedKeys := make([]string, 0)
	cdataKeys := make([]string, 0)

	for _, opt := range opts {
		if _opt, ok := opt.(*xmlKeyTagOption); ok {
			if _opt != nil && len(keyTagEntries) < 1 {
				keyTagEntries = _opt.getEntries()
			}

			continue
		}

		if _opt, ok := opt.(*xmlKeysOption); ok {
			if _opt != nil && len(includedKeys) < 1 && _opt.mode == "include" && len(_opt.keys) > 0 {
				includedKeys = _opt.keys
			} else if _opt != nil && len(excludedKeys) < 1 && _opt.mode == "exclude" && len(_opt.keys) > 0 {
				excludedKeys = _opt.keys
			}

			continue
		}

		if _opt, ok := opt.(*xmlCDataOption); ok {
			if _opt != nil && len(cdataKeys) < 1 && len(_opt.keys) > 0 {
				cdataKeys = _opt.keys
			}

			continue
		}
	}

	data := make(map[string]string)

	for key, value := range map1 {
		contents := castx.ToString(value)
		isIncluded := true

		if len(includedKeys) > 0 {
			isIncluded = false

			for _, _key := range includedKeys {
				if _key == key {
					isIncluded = true
					break
				}
			}
		}

		if !isIncluded {
			continue
		}

		var isExcluded bool

		for _, _key := range excludedKeys {
			if _key == key {
				isExcluded = true
				break
			}
		}

		if isExcluded {
			continue
		}

		var entry xmlKeyTagEntry

		for _, _entry := range keyTagEntries {
			if _entry.key == key {
				entry = _entry
				break
			}
		}

		if entry.key == "" || entry.tag == "" || entry.key == entry.tag {
			data[key] = contents
		} else {
			for idx, _key := range cdataKeys {
				if _key == entry.key {
					cdataKeys[idx] = entry.tag
					break
				}
			}

			data[entry.tag] = contents
		}
	}

	if len(data) < 1 {
		return ""
	}

	buf, err := xml.Marshal(data)

	if err != nil {
		return ""
	}

	contents := string(buf)

	for _, tag := range cdataKeys {
		re := regexp.MustCompile(fmt.Sprintf("<%s>([^<]*)</%s>", tag, tag))
		matches := re.FindAllStringSubmatch(contents, -1)

		if len(matches) < 1 {
			continue
		}

		for _, m := range matches {
			repl := fmt.Sprintf("<%s><![CDATA[%s]]></%s>", tag, m[1], tag)
			contents = strings.Replace(contents, m[0], repl, 1)
		}
	}

	return contents
}

func ToXmlFromStringMapString(map1 map[string]string, opts ...xmlOption) string {
	if len(map1) < 1 {
		return ""
	}

	keyTagEntries := make([]xmlKeyTagEntry, 0)
	includedKeys := make([]string, 0)
	excludedKeys := make([]string, 0)
	cdataKeys := make([]string, 0)

	for _, opt := range opts {
		if _opt, ok := opt.(*xmlKeyTagOption); ok {
			if _opt != nil && len(keyTagEntries) < 1 {
				keyTagEntries = _opt.getEntries()
			}

			continue
		}

		if _opt, ok := opt.(*xmlKeysOption); ok {
			if _opt != nil && len(includedKeys) < 1 && _opt.mode == "include" && len(_opt.keys) > 0 {
				includedKeys = _opt.keys
			} else if _opt != nil && len(excludedKeys) < 1 && _opt.mode == "exclude" && len(_opt.keys) > 0 {
				excludedKeys = _opt.keys
			}

			continue
		}

		if _opt, ok := opt.(*xmlCDataOption); ok {
			if _opt != nil && len(cdataKeys) < 1 && len(_opt.keys) > 0 {
				cdataKeys = _opt.keys
			}

			continue
		}
	}

	data := make(map[string]string)

	for key, value := range map1 {
		isIncluded := true

		if len(includedKeys) > 0 {
			isIncluded = false

			for _, _key := range includedKeys {
				if _key == key {
					isIncluded = true
					break
				}
			}
		}

		if !isIncluded {
			continue
		}

		var isExcluded bool

		for _, _key := range excludedKeys {
			if _key == key {
				isExcluded = true
				break
			}
		}

		if isExcluded {
			continue
		}

		var entry xmlKeyTagEntry

		for _, _entry := range keyTagEntries {
			if _entry.key == key {
				entry = _entry
				break
			}
		}

		if entry.key == "" || entry.tag == "" || entry.key == entry.tag {
			data[key] = value
		} else {
			for idx, _key := range cdataKeys {
				if _key == entry.key {
					cdataKeys[idx] = entry.tag
					break
				}
			}

			data[entry.tag] = value
		}
	}

	if len(data) < 1 {
		return ""
	}

	buf, err := xml.Marshal(data)

	if err != nil {
		return ""
	}

	contents := string(buf)

	for _, tag := range cdataKeys {
		re := regexp.MustCompile(fmt.Sprintf("<%s>([^<]*)</%s>", tag, tag))
		matches := re.FindAllStringSubmatch(contents, -1)

		if len(matches) < 1 {
			continue
		}

		for _, m := range matches {
			repl := fmt.Sprintf("<%s><![CDATA[%s]]></%s>", tag, m[1], tag)
			contents = strings.Replace(contents, m[0], repl, 1)
		}
	}

	return contents
}

func ToStruct(src map[string]interface{}, dst interface{}, opts ...structOption) {
	if dst == nil || len(src) < 1 {
		return
	}

	rt := reflect.TypeOf(dst)
	var isPtr bool

	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		isPtr = true
	}

	if rt.Kind() != reflect.Struct {
		return
	}

	keyFieldEntries := make([]structKeyFieldEntry, 0)
	includedKeys := make([]string, 0)
	excludedKeys := make([]string, 0)

	for _, opt := range opts {
		if _opt, ok := opt.(*structKeyFieldOption); ok {
			if _opt != nil && len(keyFieldEntries) < 1 {
				keyFieldEntries = _opt.getEntries()
			}

			continue
		}

		if _opt, ok := opt.(*structKeysOption); ok {
			if _opt != nil && len(includedKeys) < 1 && _opt.mode == "include" && len(_opt.keys) > 0 {
				includedKeys = _opt.keys
			} else if _opt != nil && len(excludedKeys) < 1 && _opt.mode == "exclude" && len(_opt.keys) > 0 {
				excludedKeys = _opt.keys
			}
		}
	}

	data := make(map[string]interface{}, len(src))

	for key, value := range src {
		isIncluded := true

		if len(includedKeys) > 0 {
			isIncluded = false

			for _, _key := range includedKeys {
				if _key == key {
					isIncluded = true
					break
				}
			}
		}

		if !isIncluded {
			continue
		}

		var isExcluded bool

		for _, _key := range excludedKeys {
			if _key == key {
				isExcluded = true
				break
			}
		}

		if isExcluded {
			continue
		}

		data[key] = value
	}

	if len(data) < 1 {
		return
	}

	ptr := reflectx.Ptr(dst)

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		mapkey := strings.ToLower(getMapKeyByStructField(field, keyFieldEntries))

		if mapkey == "" {
			continue
		}

		var mapvalue interface{}

		for key, value := range data {
			key = strings.ReplaceAll(key, "-", "")
			key = strings.ReplaceAll(key, "_", "")
			key = strings.TrimSpace(strings.ToLower(key))

			if key == mapkey {
				mapvalue = value
				break
			}
		}

		if mapvalue == nil {
			continue
		}

		if reflectx.SetStructFieldValue(ptr, field, mapvalue) {
			continue
		}

		var rv reflect.Value

		if isPtr {
			rv = reflect.ValueOf(dst).Elem()
		} else {
			rv = reflect.ValueOf(dst)
		}

		unsafeSetStructFieldValue(rv, i, mapvalue)
	}
}

func DeepMerge(src, dst interface{}) (interface{}, error) {
	srcType := reflect.TypeOf(src)
	dstType := reflect.TypeOf(dst)

	if srcType.Kind() != dstType.Kind() {
		return nil, errors.New("type not match")
	}

	switch srcType.Kind() {
	case reflect.Map:
		srcMap := castx.ToStringMap(src)

		for dstKey, dstVal := range castx.ToStringMap(dst) {
			srcVal, ok := srcMap[dstKey]

			if !ok {
				srcMap[dstKey] = dstVal
			} else {
				mergedVal, err := DeepMerge(srcVal, dstVal)

				if err != nil {
					return nil, err
				}

				srcMap[dstKey] = mergedVal
			}
		}

		return srcMap, nil
	case reflect.Slice:
		srcSlice := convertSlice(src)
		dstSlice := convertSlice(dst)
		return append(srcSlice, dstSlice...), nil
	default:
		return src, nil
	}
}

// @var string[]|string keys
func RemoveKeys(map1 map[string]interface{}, keys interface{}) {
	keysSlice1 := make([]string, 0)

	if _keys, ok := keys.([]string); ok && len(_keys) > 0 {
		keysSlice1 = _keys
	} else if _keys, ok := keys.(string); ok && _keys != "" {
		re := regexp.MustCompile(RegexConst.CommaSep)
		keysSlice1 = re.Split(_keys, -1)
	}

	if len(keysSlice1) < 1 {
		return
	}

	keysSlice2 := make([]string, 0)

	for key := range map1 {
		for _, cmpkey := range keysSlice1 {
			if key == cmpkey {
				keysSlice2 = append(keysSlice2, key)
				break
			}
		}
	}

	for _, key := range keysSlice2 {
		delete(map1, key)
	}
}

// @var string[]|string keys
func RemoveKeysFromStringMapString(map1 map[string]string, keys interface{}) {
	keysSlice1 := make([]string, 0)

	if _keys, ok := keys.([]string); ok && len(_keys) > 0 {
		keysSlice1 = _keys
	} else if _keys, ok := keys.(string); ok && _keys != "" {
		re := regexp.MustCompile(RegexConst.CommaSep)
		keysSlice1 = re.Split(_keys, -1)
	}

	if len(keysSlice1) < 1 {
		return
	}

	keysSlice2 := make([]string, 0)

	for key := range map1 {
		for _, cmpkey := range keysSlice1 {
			if key == cmpkey {
				keysSlice2 = append(keysSlice2, key)
				break
			}
		}
	}

	for _, key := range keysSlice2 {
		delete(map1, key)
	}
}

// @var string[]|string keys
func CopyFields(map1 map[string]interface{}, keys interface{}) map[string]interface{} {
	keysToCopy := make([]string, 0)

	if _keys, ok := keys.([]string); ok && len(_keys) > 0 {
		keysToCopy = _keys
	} else if _keys, ok := keys.(string); ok && _keys != "" {
		re := regexp.MustCompile(RegexConst.CommaSep)
		keysToCopy = re.Split(_keys, -1)
	}

	ret := make(map[string]interface{})

	if len(keysToCopy) < 1 {
		return ret
	}

	for key, value := range map1 {
		var matched bool

		for _, _key := range keysToCopy {
			if _key == key {
				matched = true
				break
			}
		}

		if matched {
			ret[key] = value
		}
	}

	return ret
}

// @var string[]|string keys
func CopyFieldsToStringMapString(map1 map[string]interface{}, keys interface{}) map[string]string {
	keysToCopy := make([]string, 0)

	if _keys, ok := keys.([]string); ok && len(_keys) > 0 {
		keysToCopy = _keys
	} else if _keys, ok := keys.(string); ok && _keys != "" {
		re := regexp.MustCompile(RegexConst.CommaSep)
		keysToCopy = re.Split(_keys, -1)
	}

	ret := make(map[string]string)

	if len(keysToCopy) < 1 {
		return ret
	}

	for key, value := range map1 {
		var matched bool

		for _, _key := range keysToCopy {
			if _key == key {
				matched = true
				break
			}
		}

		if matched {
			ret[key] = castx.ToString(value)
		}
	}

	return ret
}

// @var string[]|string keys
func CopyFieldsFromStringMapString(map1 map[string]string, keys interface{}) map[string]string {
	keysToCopy := make([]string, 0)

	if _keys, ok := keys.([]string); ok && len(_keys) > 0 {
		keysToCopy = _keys
	} else if _keys, ok := keys.(string); ok && _keys != "" {
		re := regexp.MustCompile(RegexConst.CommaSep)
		keysToCopy = re.Split(_keys, -1)
	}

	ret := make(map[string]string)

	if len(keysToCopy) < 1 {
		return ret
	}

	for key, value := range map1 {
		var matched bool

		for _, _key := range keysToCopy {
			if _key == key {
				matched = true
				break
			}
		}

		if matched {
			ret[key] = value
		}
	}

	return ret
}

func convertSlice(arg0 interface{}) []interface{} {
	ret := make([]interface{}, 0)

	switch arg0.(type) {
	case []interface{}:
		return arg0.([]interface{})
	case []string:
		for _, v := range arg0.([]string) {
			ret = append(ret, v)
		}

		return ret
	case []int:
		for _, v := range arg0.([]int) {
			ret = append(ret, v)
		}

		return ret
	case []float64:
		for _, v := range arg0.([]float64) {
			ret = append(ret, v)
		}

		return ret
	case []float32:
		for _, v := range arg0.([]float32) {
			ret = append(ret, v)
		}

		return ret
	case []byte:
		return append(ret, arg0)
	}

	return nil
}

func getMapKeyByStructField(field reflect.StructField, entries []structKeyFieldEntry) string {
	for _, entry := range entries {
		if strings.ToLower(entry.fieldName) == strings.ToLower(field.Name) {
			return entry.key
		}
	}

	key := strings.TrimSpace(field.Tag.Get("mapkey"))

	if key != "" {
		return key
	}

	return strings.ToLower(field.Name[:1]) + field.Name[1:]
}

func unsafeSetStructFieldValue(rv reflect.Value, idx int, value interface{}) {
	defer func() {
		recover()
	}()

	rv.Field(idx).Set(reflect.ValueOf(value))
}
