package pagerx

import (
	"github.com/meiguonet/mgboot-go-common/util/castx"
)

func IsNoPage(arg0 interface{}) bool {
	if b1, ok := arg0.(bool); ok {
		return b1
	}

	if n1, ok := arg0.(int); ok {
		return n1 == 1
	}

	if map1, ok := arg0.(map[string]interface{}); ok {
		if _, ok := map1["disablePagination"]; ok {
			return map1["disablePagination"] == 1
		} else if _, ok := map1["nopage"]; ok {
			return map1["nopage"] == 1
		}
	}

	return false
}

func GetPageAndPageSize(args ...interface{}) (int, int) {
	switch len(args) {
	case 1:
		if map1, ok := args[0].(map[string]interface{}); ok {
			page := castx.ToInt(map1["page"])

			if page < 1 {
				page = 1
			}

			pageSize := castx.ToInt(map1["pageSize"])

			if pageSize < 1 {
				pageSize = 20
			}

			return page, pageSize
		}

		if pageSize, ok := args[0].(int); ok && pageSize > 0 {
			return 1, pageSize
		}

		return 1, 20
	case 2:
		page, ok1 := args[0].(int)
		pageSize, ok2 := args[1].(int)

		if !ok1 || !ok2 {
			return 1, 20
		}

		if page < 1 {
			page = 1
		}

		if pageSize < 1 {
			pageSize = 20
		}

		return page, pageSize
	default:
		return 1, 20
	}
}

func ForPage(cnt, page, pageSize int) (int, int) {
	if cnt < 1 || page < 1 || pageSize < 1 {
		return -1, -1
	}

	pageTotal := cnt / pageSize

	if cnt % pageSize != 0 {
		pageTotal++
	}

	if page > pageTotal {
		return -1, -1
	}

	if page == pageTotal {
		return (page - 1) * pageSize, -1
	}

	n1 := (page - 1) * pageSize
	n2 := page * pageSize

	if n2 >= cnt {
		n2 = -1
	}

	return n1, n2
}
