package jsonx

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strings"
)

var re1 = regexp.MustCompile(`"(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2})[^"]+"`)

type ToJsonOption struct {
	HandleTimeField   bool
	StripZeroTimePart bool
}

func MapFrom(arg0 interface{}) map[string]interface{} {
	var buf []byte

	if _buf, ok := arg0.([]byte); ok && len(_buf) > 0 {
		buf = _buf
	} else if s1, ok := arg0.(string); ok && s1 != "" {
		buf = []byte(s1)
	}

	var map1 map[string]interface{}

	if len(buf) < 1 || json.Unmarshal(buf, &map1) != nil {
		return map[string]interface{}{}
	}

	return map1
}

func ArrayFrom(arg0 interface{}) []interface{} {
	var buf []byte

	if _buf, ok := arg0.([]byte); ok && len(_buf) > 0 {
		buf = _buf
	} else if s1, ok := arg0.(string); ok && s1 != "" {
		buf = []byte(s1)
	}

	var list []interface{}

	if len(buf) < 1 || json.Unmarshal(buf, &list) != nil {
		return []interface{}{}
	}

	return list
}

func ToJson(arg0 interface{}, opts ...ToJsonOption) string {
	buf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)

	if encoder.Encode(arg0) != nil {
		return "{}"
	}

	contents := buf.String()
	var opt ToJsonOption

	if len(opts) > 0 {
		opt = opts[0]
	}

	if !opt.HandleTimeField {
		return contents
	}

	matches := re1.FindAllStringSubmatch(contents, -1)

	if len(matches) < 1 {
		return contents
	}

	for _, groups := range matches {
		if len(groups) < 3 {
			continue
		}

		if groups[2] == "00:00:00" && opt.StripZeroTimePart {
			contents = strings.Replace(contents, groups[0], groups[1], 1)
			continue
		}

		contents = strings.Replace(contents, groups[0], groups[1] + " " + groups[2], 1)
	}

	return contents
}
