package mapx

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/meiguonet/mgboot-go-common/enum/DatetimeFormat"
	"github.com/meiguonet/mgboot-go-common/enum/RegexConst"
	"github.com/meiguonet/mgboot-go-common/util/castx"
	"github.com/meiguonet/mgboot-go-common/util/numberx"
	"github.com/meiguonet/mgboot-go-common/util/reflectx"
	"github.com/meiguonet/mgboot-go-common/util/stringx"
	"math"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unsafe"
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
	if len(map1) < 1 {
		return
	}

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
	if len(map1) < 1 {
		return
	}

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
func CopyFields(srcMap map[string]interface{}, keys interface{}) map[string]interface{} {
	if len(srcMap) < 1 {
		return map[string]interface{}{}
	}

	_keys := make([]string, 0)

	if a1, ok := keys.([]string); ok && len(a1) > 0 {
		_keys = a1
	} else if s1, ok := keys.(string); ok && s1 != "" {
		re := regexp.MustCompile(RegexConst.CommaSep)
		_keys = re.Split(s1, -1)
	}

	if len(_keys) < 1 {
		return map[string]interface{}{}
	}

	dstKeys := map[string]string{}

	for idx, key := range _keys {
		if strings.Contains(key, ":") {
			p1 := stringx.SubstringBefore(key, ":")
			p2 := stringx.SubstringAfter(key, ":")
			dstKeys[strings.ToLower(p1)] = p2
			_keys[idx] = p1
			continue
		}

		if strings.Contains(key, "#") {
			p1 := stringx.SubstringBefore(key, "#")
			p2 := stringx.SubstringAfter(key, "#")
			dstKeys[strings.ToLower(p1)] = p2
			_keys[idx] = p1
			continue
		}
	}

	dstMap := map[string]interface{}{}

	for _, key := range _keys {
		var matchedKey string

		for cmpkey := range srcMap {
			if strings.ToLower(cmpkey) == strings.ToLower(key) {
				matchedKey = cmpkey
				break
			}
		}

		if matchedKey == "" {
			continue
		}

		dstKey := dstKeys[strings.ToLower(matchedKey)]

		if dstKey == "" {
			dstKey = matchedKey
		}

		dstMap[dstKey] = srcMap[matchedKey]
	}

	return dstMap
}

// @var string[]|string keys
func CopyFieldsToStringMapString(srcMap map[string]interface{}, keys interface{}) map[string]string {
	if len(srcMap) < 1 {
		return map[string]string{}
	}

	_keys := make([]string, 0)

	if a1, ok := keys.([]string); ok && len(a1) > 0 {
		_keys = a1
	} else if s1, ok := keys.(string); ok && s1 != "" {
		re := regexp.MustCompile(RegexConst.CommaSep)
		_keys = re.Split(s1, -1)
	}

	if len(_keys) < 1 {
		return map[string]string{}
	}

	dstKeys := map[string]string{}

	for idx, key := range _keys {
		if strings.Contains(key, ":") {
			p1 := stringx.SubstringBefore(key, ":")
			p2 := stringx.SubstringAfter(key, ":")
			dstKeys[strings.ToLower(p1)] = p2
			_keys[idx] = p1
			continue
		}

		if strings.Contains(key, "#") {
			p1 := stringx.SubstringBefore(key, "#")
			p2 := stringx.SubstringAfter(key, "#")
			dstKeys[strings.ToLower(p1)] = p2
			_keys[idx] = p1
			continue
		}
	}

	dstMap := map[string]string{}

	for _, key := range _keys {
		var matchedKey string

		for cmpkey := range srcMap {
			if strings.ToLower(cmpkey) == strings.ToLower(key) {
				matchedKey = cmpkey
				break
			}
		}

		if matchedKey == "" {
			continue
		}

		value, err := castx.ToStringE(srcMap[matchedKey])

		if err != nil {
			continue
		}

		dstKey := dstKeys[strings.ToLower(matchedKey)]

		if dstKey == "" {
			dstKey = matchedKey
		}

		dstMap[dstKey] = value
	}

	return dstMap
}

// @var string[]|string keys
func CopyFieldsFromStringMapString(srcMap map[string]string, keys interface{}) map[string]string {
	if len(srcMap) < 1 {
		return map[string]string{}
	}

	_keys := make([]string, 0)

	if a1, ok := keys.([]string); ok && len(a1) > 0 {
		_keys = a1
	} else if s1, ok := keys.(string); ok && s1 != "" {
		re := regexp.MustCompile(RegexConst.CommaSep)
		_keys = re.Split(s1, -1)
	}

	if len(_keys) < 1 {
		return map[string]string{}
	}

	dstKeys := map[string]string{}

	for idx, key := range _keys {
		if strings.Contains(key, ":") {
			p1 := stringx.SubstringBefore(key, ":")
			p2 := stringx.SubstringAfter(key, ":")
			dstKeys[strings.ToLower(p1)] = p2
			_keys[idx] = p1
			continue
		}

		if strings.Contains(key, "#") {
			p1 := stringx.SubstringBefore(key, "#")
			p2 := stringx.SubstringAfter(key, "#")
			dstKeys[strings.ToLower(p1)] = p2
			_keys[idx] = p1
			continue
		}
	}

	dstMap := map[string]string{}

	for _, key := range _keys {
		var matchedKey string

		for cmpkey := range srcMap {
			if strings.ToLower(cmpkey) == strings.ToLower(key) {
				matchedKey = cmpkey
				break
			}
		}

		if matchedKey == "" {
			continue
		}

		dstKey := dstKeys[strings.ToLower(matchedKey)]

		if dstKey == "" {
			dstKey = matchedKey
		}

		dstMap[dstKey] = srcMap[matchedKey]
	}

	return dstMap
}

// @var string[]|string rules
func FromRequestParam(srcMap map[string]interface{}, rules ...interface{}) map[string]interface{} {
	if len(srcMap) < 1 {
		return map[string]interface{}{}
	}

	var _rules []string

	if len(rules) > 0 {
		if a1, ok := rules[0].([]string); ok && len(a1) > 0 {
			_rules = a1
		} else if s1, ok := rules[0].(string); ok && s1 != "" {
			re := regexp.MustCompile(RegexConst.CommaSep)
			_rules = re.Split(s1, -1)
		}
	}

	if len(_rules) < 1 {
		return srcMap
	}

	dstKeys := make([]string, 0)

	for idx, rule := range _rules {
		if strings.Contains(rule, "#") {
			dstKeys = append(dstKeys, stringx.SubstringBefore(rule, "#"))
			_rules[idx] = stringx.SubstringAfter(rule, "#")
		} else {
			dstKeys = append(dstKeys, "")
		}
	}

	dstMap := map[string]interface{}{}
	re1 := regexp.MustCompile(`:[^:]+$`)
	re2 := regexp.MustCompile(`:[0-9]+$`)

	for idx, rule := range _rules {
		name := rule
		typ := 1
		mode := 2
		dv := ""

		if strings.HasPrefix(name, "i:") {
			name = stringx.SubstringAfter(name, ":")
			typ = 2

			if re1.MatchString(name) {
				dv = stringx.SubstringAfterLast(name, ":")
				name = re1.ReplaceAllString(name, "")
			}
		} else if strings.HasPrefix(name, "d:") {
			name = stringx.SubstringAfter(name, ":")
			typ = 3

			if re1.MatchString(name) {
				dv = stringx.SubstringAfterLast(name, ":")
				name = re1.ReplaceAllString(name, "")
			}
		} else if strings.HasPrefix(name, "s:") {
			name = stringx.SubstringAfter(name, ":")

			if re2.MatchString(name) {
				s1 := stringx.SubstringAfterLast(name, ":")
				mode = castx.ToInt(s1, 2)
				name = re2.ReplaceAllString(name, "")
			}
		} else if re2.MatchString(name) {
			s1 := stringx.SubstringAfterLast(name, ":")
			mode = castx.ToInt(s1, 2)
			name = re2.ReplaceAllString(name, "")
		}

		if strings.Contains(name, ":") {
			name = stringx.SubstringBefore(name, ":")
		}

		if name == "" {
			continue
		}

		var dstKey string

		if dstKeys[idx] != "" {
			dstKey = dstKeys[idx]
		} else {
			dstKey = name
		}

		switch typ {
		case 1:
			value := castx.ToString(srcMap[name])

			switch mode {
			case 1, 2:
				value = stringx.StripTags(value)
			}

			dstMap[dstKey] = value
		case 2:
			var value int

			if n1, err := castx.ToIntE(dv); err == nil {
				value = castx.ToInt(srcMap[name], n1)
			} else {
				value = castx.ToInt(srcMap[name])
			}

			dstMap[dstKey] = value
		case 3:
			var value float64

			if n1, err := castx.ToFloat64E(dv); err == nil {
				value = castx.ToFloat64(srcMap[name], n1)
			} else {
				value = castx.ToFloat64(srcMap[name])
			}

			dstMap[dstKey] = numberx.ToDecimalString(value)
		}
	}

	return dstMap
}

func BindToDto(srcMap map[string]interface{}, dto interface{}) error {
	if len(srcMap) < 1 || dto == nil {
		return errors.New("in mapx.BindToDto function, argument #0 must be a not empty map[string]interface{}")
	}

	ex1 := errors.New("in mapx.BindToDto function, argument #1 must be a struct pointer")
	rt := reflect.TypeOf(dto)

	if rt.Kind() != reflect.Ptr {
		return ex1
	}

	rt = rt.Elem()

	if rt.Kind() != reflect.Struct || rt.NumField() < 1 {
		return ex1
	}

	defer func() {
		recover()
	}()

	structPtr := getStructPtr(dto)
	re1 := regexp.MustCompile(`MapKey:([^\x20\t]+)`)
	re2 := regexp.MustCompile(`SecurityMode:([^\x20\t]+)`)
	re3 := regexp.MustCompile(`DefaultValue:([^\x20\t]+)`)

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		tag := field.Tag.Get("DtoBind")
		var key string
		matches := re1.FindStringSubmatch(tag)

		if len(matches) > 1 {
			key = matches[1]
		} else {
			key = field.Name
		}

		mode := 2
		matches = re2.FindStringSubmatch(tag)

		if len(matches) > 1 {
			if n1, err := castx.ToIntE(matches[1]); err == nil {
				mode = n1
			}
		}

		dv := ""

		matches = re3.FindStringSubmatch(tag)

		if len(matches) > 1 {
			dv = matches[1]
		}

		switch field.Type.Kind() {
		case reflect.String:
			bindDtoStringField(structPtr, field, srcMap[key], mode, dv)
			continue
		case reflect.Bool:
			var defaultValue bool

			if b1, err := castx.ToBoolE(dv); err == nil {
				defaultValue = b1
			}

			bindDtoBoolField(structPtr, field, srcMap[key], defaultValue)
			continue
		case reflect.Int:
			defaultValue := math.MinInt32

			if n1, err := castx.ToIntE(dv); err == nil {
				defaultValue = n1
			}

			bindDtoIntField(structPtr, field, srcMap[key], defaultValue)
			continue
		case reflect.Int64:
			defaultValue := int64(math.MinInt64)

			if n1, err := castx.ToInt64E(dv); err == nil {
				defaultValue = n1
			}

			bindDtoInt64Field(structPtr, field, srcMap[key], defaultValue)
			continue
		case reflect.Float32:
			defaultValue := float32(math.SmallestNonzeroFloat32)

			if n1, err := castx.ToFloat32E(dv); err == nil {
				defaultValue = n1
			}

			bindDtoFloat32Field(structPtr, field, srcMap[key], defaultValue)
			continue
		case reflect.Float64:
			defaultValue := math.SmallestNonzeroFloat64

			if n1, err := castx.ToFloat64E(dv); err == nil {
				defaultValue = n1
			}

			bindDtoFloat64Field(structPtr, field, srcMap[key], defaultValue)
			continue
		}

		if strings.Contains(field.Type.String(), "*time.Time") {
			bindDtoTimePtrField(structPtr, field, srcMap[key], tag)
			continue
		}

		if strings.Contains(field.Type.String(), "time.Time") {
			bindDtoTimeField(structPtr, field, srcMap[key], tag)
			continue
		}

		return errors.New("in mapx.BindToDto function, unsupported dto field type: " + field.Type.String())
	}

	return nil
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

func bindDtoStringField(
	structPtr unsafe.Pointer,
	field reflect.StructField,
	value interface{},
	securityMode int,
	defaultValue string,
) {
	var realValue string

	if s1, err := castx.ToStringE(value); err == nil && s1 != "" {
		realValue = s1
	} else {
		realValue = defaultValue
	}

	if securityMode == 1 || securityMode == 2 {
		realValue = stringx.StripTags(realValue)
	}

	fieldPtr := getStructFieldPtr(structPtr, field)
	*((*string)(unsafe.Pointer(fieldPtr))) = realValue
}

func bindDtoBoolField(
	structPtr unsafe.Pointer,
	field reflect.StructField,
	value interface{},
	defaultValue bool,
) {
	var realValue bool

	if b1, err := castx.ToBoolE(value); err == nil {
		realValue = b1
	} else {
		realValue = defaultValue
	}

	fieldPtr := getStructFieldPtr(structPtr, field)
	*((*bool)(unsafe.Pointer(fieldPtr))) = realValue
}

func bindDtoIntField(
	structPtr unsafe.Pointer,
	field reflect.StructField,
	value interface{},
	defaultValue int,
) {
	var realValue int

	if n1, err := castx.ToIntE(value); err == nil {
		realValue = n1
	} else {
		realValue = defaultValue
	}

	fieldPtr := getStructFieldPtr(structPtr, field)
	*((*int)(unsafe.Pointer(fieldPtr))) = realValue
}

func bindDtoInt64Field(
	structPtr unsafe.Pointer,
	field reflect.StructField,
	value interface{},
	defaultValue int64,
) {
	var realValue int64

	if n1, err := castx.ToInt64E(value); err == nil {
		realValue = n1
	} else {
		realValue = defaultValue
	}

	fieldPtr := getStructFieldPtr(structPtr, field)
	*((*int64)(unsafe.Pointer(fieldPtr))) = realValue
}

func bindDtoFloat32Field(
	structPtr unsafe.Pointer,
	field reflect.StructField,
	value interface{},
	defaultValue float32,
) {
	var realValue float32

	if n1, err := castx.ToFloat32E(value); err == nil {
		realValue = n1
	} else {
		realValue = defaultValue
	}

	fieldPtr := getStructFieldPtr(structPtr, field)
	*((*float32)(unsafe.Pointer(fieldPtr))) = realValue
}

func bindDtoFloat64Field(
	structPtr unsafe.Pointer,
	field reflect.StructField,
	value interface{},
	defaultValue float64,
) {
	var realValue float64

	if n1, err := castx.ToFloat64E(value); err == nil {
		realValue = n1
	} else {
		realValue = defaultValue
	}

	fieldPtr := getStructFieldPtr(structPtr, field)
	*((*float64)(unsafe.Pointer(fieldPtr))) = realValue
}

func bindDtoTimePtrField(
	structPtr unsafe.Pointer,
	field reflect.StructField,
	value interface{},
	tag string,
) {
	fieldPtr := getStructFieldPtr(structPtr, field)

	if d1, ok := value.(time.Time); ok && !d1.IsZero() {
		*((**time.Time)(unsafe.Pointer(fieldPtr))) = &d1
		return
	}

	if d1, ok := value.(*time.Time); ok && d1 != nil && !d1.IsZero() {
		*((**time.Time)(unsafe.Pointer(fieldPtr))) = d1
		return
	}

	s1 := castx.ToString(value)

	if s1 == "" {
		return
	}

	var d1 time.Time
	var err error

	if strings.Contains(tag, "DateOnly") {
		d1, err = time.Parse(DatetimeFormat.DateOnly, s1)
	} else if strings.Contains(tag, "TimeOnly") {
		d1, err = time.Parse(DatetimeFormat.TimeOnly, s1)
	} else {
		d1, err = time.Parse(DatetimeFormat.Full, s1)
	}

	if err != nil {
		return
	}

	*((**time.Time)(unsafe.Pointer(fieldPtr))) = &d1
}

func bindDtoTimeField(
	structPtr unsafe.Pointer,
	field reflect.StructField,
	value interface{},
	tag string,
) {
	fieldPtr := getStructFieldPtr(structPtr, field)

	if d1, ok := value.(time.Time); ok && !d1.IsZero() {
		*((*time.Time)(unsafe.Pointer(fieldPtr))) = d1
		return
	}

	if d1, ok := value.(*time.Time); ok && d1 != nil && !d1.IsZero() {
		*((*time.Time)(unsafe.Pointer(fieldPtr))) = indirect(d1).(time.Time)
		return
	}

	s1 := castx.ToString(value)

	if s1 == "" {
		return
	}

	var d1 time.Time
	var err error

	if strings.Contains(tag, "DateOnly") {
		d1, err = time.Parse(DatetimeFormat.DateOnly, s1)
	} else if strings.Contains(tag, "TimeOnly") {
		d1, err = time.Parse(DatetimeFormat.TimeOnly, s1)
	} else {
		d1, err = time.Parse(DatetimeFormat.Full, s1)
	}

	if err != nil {
		return
	}

	*((*time.Time)(unsafe.Pointer(fieldPtr))) = d1
}

func getStructPtr(arg0 interface{}) unsafe.Pointer {
	return unsafe.Pointer(reflect.ValueOf(arg0).Pointer())
}

func getStructFieldPtr(structPtr unsafe.Pointer, field reflect.StructField) uintptr {
	return uintptr(structPtr) + field.Offset
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
