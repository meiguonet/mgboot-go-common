package AppConf

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"math"
	"meiguonet/mgboot-go-common/util/castx"
	"meiguonet/mgboot-go-common/util/stringx"
	"os"
	"strings"
	"time"
)

var _env = "dev"
var _dataDir string
var _data map[string]interface{}

func SetEnv(env string) {
	_env = env
}

func GetEnv() string {
	return _env
}

func SetDataDir(dir string) {
	if dir == "" {
		return
	}

	if stat, err := os.Stat(dir); err != nil || !stat.IsDir() {
		return
	}

	_dataDir = dir
}

func GetDataDir() string {
	return _dataDir
}

func InitFromMap(arg0 interface{}) {
	if arg0 == nil {
		return
	}

	map1 := castx.ToStringMap(arg0)
	
	if len(map1) > 0 {
		_data = map1
	}
}

func InitFromJson(arg0 interface{}) {
	var buf []byte

	if _buf, ok := arg0.([]byte); ok && len(_buf) > 0 {
		buf = _buf
	} else if reader, ok := arg0.(io.Reader); ok {
		buf, _ = ioutil.ReadAll(reader)
	} else if s1, ok := arg0.(string); ok && s1 != "" {
		buf = []byte(s1)
	}

	var data map[string]interface{}

	if len(buf) < 1 || json.Unmarshal(buf, &data) != nil {
		return
	}
	
	_data = data
}

func InitFromYaml(arg0 interface{}) {
	var buf []byte

	if _buf, ok := arg0.([]byte); ok && len(_buf) > 0 {
		buf = _buf
	} else if reader, ok := arg0.(io.Reader); ok {
		buf, _ = ioutil.ReadAll(reader)
	} else if s1, ok := arg0.(string); ok && s1 != "" {
		if stat, err := os.Stat(s1); err == nil {
			if !stat.IsDir() {
				buf, _ = ioutil.ReadFile(s1)
			}
		} else {
			buf = []byte(s1)
		}
	}

	var data map[string]interface{}

	if len(buf) < 1 || yaml.Unmarshal(buf, &data) != nil {
		return
	}
	
	_data = data
}

func GetMap(path string) map[string]interface{} {
	if len(_data) < 1 {
		return map[string]interface{}{}
	}

	if !strings.Contains(path, ".") {
		return castx.ToStringMap(getValueInternal(path))
	}

	lastKey := stringx.SubstringAfterLast(path, ".")
	keys := strings.Split(strings.TrimSuffix(path, "." + lastKey), ".")
	var map1 map[string]interface{}

	for idx, key := range keys {
		if idx == 0 {
			map1 = castx.ToStringMap(getValueInternal(key))
			continue
		}

		if len(map1) < 1 {
			break
		}

		map1 = castx.ToStringMap(getValueInternal(key, map1))
	}

	if len(map1) < 1 {
		return map[string]interface{}{}
	}

	return castx.ToStringMap(getValueInternal(lastKey, map1))
}

func GetStringMap(path string) map[string]string {
	return castx.ToStringMapString(GetMap(path))
}

func GetSlice(path string) []interface{} {
	if !strings.Contains(path, ".") {
		if len(_data) < 1 {
			return []interface{}{}
		}
		
		return castx.ToSlice(getValueInternal(path))
	}

	map1 := GetMap(stringx.SubstringBeforeLast(path, "."))
	key := stringx.SubstringAfterLast(path, ".")
	return castx.ToSlice(getValueInternal(key, map1))
}

func GetStringSlice(path string) []string {
	if !strings.Contains(path, ".") {
		if len(_data) < 1 {
			return []string{}
		}
		
		return castx.ToStringSlice(getValueInternal(path))
	}

	map1 := GetMap(stringx.SubstringBeforeLast(path, "."))
	key := stringx.SubstringAfterLast(path, ".")
	return castx.ToStringSlice(getValueInternal(key, map1))
}

func GetIntSlice(path string) []int {
	if !strings.Contains(path, ".") {
		if len(_data) < 1 {
			return []int{}
		}

		return castx.ToIntSlice(getValueInternal(path))
	}

	map1 := GetMap(stringx.SubstringBeforeLast(path, "."))
	key := stringx.SubstringAfterLast(path, ".")
	return castx.ToIntSlice(getValueInternal(key, map1))
}

func GetMapSlice(path string) []map[string]interface{} {
	if !strings.Contains(path, ".") {
		if len(_data) < 1 {
			return []map[string]interface{}{}
		}

		return castx.ToMapSlice(getValueInternal(path))
	}

	map1 := GetMap(stringx.SubstringBeforeLast(path, "."))
	key := stringx.SubstringAfterLast(path, ".")
	return castx.ToMapSlice(getValueInternal(key, map1))
}

func GetString(path string, defaultValue ...string) string {
	_defaultValue := ""

	if len(defaultValue) > 0 && defaultValue[0] != "" {
		_defaultValue = defaultValue[0]
	}

	if !strings.Contains(path, ".") {
		if len(_data) < 1 {
			return _defaultValue
		}

		if s1, err := castx.ToStringE(getValueInternal(path)); err == nil && s1 != "" {
			return s1
		}

		return _defaultValue
	}

	map1 := GetMap(stringx.SubstringBeforeLast(path, "."))
	key := stringx.SubstringAfterLast(path, ".")

	if s1, err := castx.ToStringE(getValueInternal(key, map1)); err == nil && s1 != "" {
		return s1
	}

	return _defaultValue
}

func GetInt(path string, defaultValue ...int) int {
	_defaultValue := math.MinInt32

	if len(defaultValue) > 0 {
		_defaultValue = defaultValue[0]
	}

	if !strings.Contains(path, ".") {
		if len(_data) < 1 {
			return _defaultValue
		}

		if n1, err := castx.ToIntE(getValueInternal(path)); err == nil {
			return n1
		}

		return _defaultValue
	}

	map1 := GetMap(stringx.SubstringBeforeLast(path, "."))
	key := stringx.SubstringAfterLast(path, ".")

	if n1, err := castx.ToIntE(getValueInternal(key, map1)); err == nil {
		return n1
	}

	return _defaultValue
}

func GetInt64(path string, defaultValue ...int64) int64 {
	_defaultValue := int64(math.MinInt64)

	if len(defaultValue) > 0 {
		_defaultValue = defaultValue[0]
	}

	if !strings.Contains(path, ".") {
		if len(_data) < 1 {
			return _defaultValue
		}

		if n1, err := castx.ToInt64E(getValueInternal(path)); err == nil {
			return n1
		}

		return _defaultValue
	}

	map1 := GetMap(stringx.SubstringBeforeLast(path, "."))
	key := stringx.SubstringAfterLast(path, ".")

	if n1, err := castx.ToInt64E(getValueInternal(key, map1)); err == nil {
		return n1
	}

	return _defaultValue
}

func GetFloat(path string, defaultValue ...float64) float64 {
	_defaultValue := math.SmallestNonzeroFloat64

	if len(defaultValue) > 0 {
		_defaultValue = defaultValue[0]
	}

	if !strings.Contains(path, ".") {
		if len(_data) < 1 {
			return _defaultValue
		}

		if n1, err := castx.ToFloat64E(getValueInternal(path)); err == nil {
			return n1
		}

		return _defaultValue
	}

	map1 := GetMap(stringx.SubstringBeforeLast(path, "."))
	key := stringx.SubstringAfterLast(path, ".")

	if n1, err := castx.ToFloat64E(getValueInternal(key, map1)); err == nil {
		return n1
	}

	return _defaultValue
}

func GetBoolean(path string, defaultValue ...bool) bool {
	_defaultValue := false

	if len(defaultValue) > 0 {
		_defaultValue = defaultValue[0]
	}

	if !strings.Contains(path, ".") {
		if len(_data) < 1 {
			return _defaultValue
		}

		if b1, err := castx.ToBoolE(getValueInternal(path)); err == nil {
			return b1
		}

		return _defaultValue
	}

	map1 := GetMap(stringx.SubstringBeforeLast(path, "."))
	key := stringx.SubstringAfterLast(path, ".")

	if b1, err := castx.ToBoolE(getValueInternal(key, map1)); err == nil {
		return b1
	}

	return _defaultValue
}

func GetDataSize(path string) int64 {
	return castx.ToDataSize(GetString(path))
}

func GetDuration(path string) time.Duration {
	return castx.ToDuration(GetString(path))
}

func getValueInternal(key string, source ...map[string]interface{}) interface{} {
	var data map[string]interface{}

	if len(source) > 0 {
		data = source[0]
	} else {
		data = _data
	}

	if len(data) < 1 {
		return nil
	}

	key = strings.ReplaceAll(key, "-", "")
	key = strings.ReplaceAll(key, "_", "")
	key = strings.ToLower(key)

	for compkey, value := range data {
		compkey = strings.ReplaceAll(compkey, "-", "")
		compkey = strings.ReplaceAll(compkey, "_", "")
		compkey = strings.ToLower(compkey)

		if compkey == key {
			return value
		}
	}

	return nil
}
