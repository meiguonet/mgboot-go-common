package slicex

import (
	"github.com/meiguonet/mgboot-go-common/util/castx"
	"strings"
)

func InStringSlice(needle string, array []string, ignoreCase ...bool) bool {
	if needle == "" || len(array) < 1 {
		return false
	}

	_ignoreCase := false

	if len(ignoreCase) > 0 {
		_ignoreCase = ignoreCase[0]
	}

	var s1 string

	if _ignoreCase {
		s1 = strings.ToLower(needle)
	} else {
		s1 = needle
	}

	for _, s2 := range array {
		if _ignoreCase {
			s2 = strings.ToLower(s2)
		}

		if s1 == s2 {
			return true
		}
	}

	return false
}

func InIntSlice(needle int, array []int) bool {
	if len(array) < 1 {
		return false
	}

	for _, num := range array {
		if needle == num {
			return true
		}
	}

	return false
}

func StringSliceUnique(list []string) []string {
	ret := make([]string, 0)

	if len(list) < 1 {
		return ret
	}

	for _, s1 := range list {
		if InStringSlice(s1, ret) {
			continue
		}

		ret = append(ret, s1)
	}

	return ret
}

func IntSliceUnique(list []int) []int {
	ret := make([]int, 0)

	if len(list) < 1 {
		return ret
	}

	for _, n1 := range list {
		if InIntSlice(n1, ret) {
			continue
		}

		ret = append(ret, n1)
	}

	return ret
}

func StringColumn(
	columnName string,
	list []map[string]interface{},
	matcher func(arg0 string) bool,
	unique ...bool,
) []string {
	if columnName == "" || len(list) < 1 {
		return []string{}
	}

	var _unique bool

	if len(unique) > 0 {
		_unique = unique[0]
	}

	ret := make([]string, 0)

	for _, item := range list {
		s1, err := castx.ToStringE(item[columnName])

		if err != nil || s1 == "" {
			continue
		}

		if matcher != nil && !matcher(s1) {
			continue
		}

		if _unique && InStringSlice(s1, ret) {
			continue
		}

		ret = append(ret, s1)
	}

	return ret
}

func IntColumn(
	columnName string,
	list []map[string]interface{},
	matcher func(arg0 int) bool,
	unique ...bool,
) []int {
	if columnName == "" || len(list) < 1 {
		return []int{}
	}

	var _unique bool

	if len(unique) > 0 {
		_unique = unique[0]
	}

	ret := make([]int, 0)

	for _, item := range list {
		n1, err := castx.ToIntE(item[columnName])

		if err != nil {
			continue
		}

		if matcher != nil && !matcher(n1) {
			continue
		}

		if _unique && InIntSlice(n1, ret) {
			continue
		}

		ret = append(ret, n1)
	}

	return ret
}

func MapSliceForPage(list []map[string]interface{}, page, pageSize int) []map[string]interface{} {
	n1, n2 := ForPage(len(list), page, pageSize)

	if n1 < 0 && n2 < 0 {
		return []map[string]interface{}{}
	}

	return list[n1:n2]
}

func ForPage(cnt, page, pageSize int) (int, int) {
	if cnt < 1 || page < 1 || pageSize < 1 {
		return -1, -1
	}

	pageTotal := cnt / pageSize

	if cnt % pageSize != 0 {
		pageTotal += 1
	}

	if page > pageTotal {
		return -1, -1
	}

	n1 := (page - 1) * pageSize
	n2 := page * pageSize

	if n2 > cnt {
		n2 = cnt
	}

	return n1, n2
}
