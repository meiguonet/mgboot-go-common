package CommonBo

import (
	"math"
	"meiguonet/mgboot-go-common/enum/RegexConst"
	"meiguonet/mgboot-go-common/util/castx"
	"meiguonet/mgboot-go-common/util/slicex"
	"meiguonet/mgboot-go-common/util/stringx"
)

type pager struct {
	recordTotal int
	currentPage int
	pageSize    int
	pageStep    int
}

func NewPager(cnt, page, pageSize int, step ...int) *pager {
	pageStep := 5

	if len(step) > 0 && step[0] > 0 {
		pageStep = step[0]
	}

	return &pager{
		recordTotal: cnt,
		currentPage: page,
		pageSize:    pageSize,
		pageStep:    pageStep,
	}
}

func (bo *pager) ToMap(includeFields ...interface{}) map[string]interface{} {
	_includeFields := make([]string, 0)

	if len(includeFields) > 0 {
		if a1, ok := includeFields[0].([]string); ok {
			_includeFields = a1
		} else if s1, ok := includeFields[0].(string); ok && s1 != "" {
			_includeFields = stringx.SplitWithRegexp(s1, RegexConst.CommaSep)
		}
	}

	if len(_includeFields) < 1 {
		return map[string]interface{}{}
	}

	pageTotal := 0

	if bo.recordTotal > 0 {
		n1 := math.Ceil(castx.ToFloat64(bo.recordTotal) / castx.ToFloat64(bo.pageSize))
		pageTotal = castx.ToInt(n1)
	}

	pageList := make([]int, 0)

	if pageTotal > 0 {
		n1 := math.Ceil(castx.ToFloat64(bo.currentPage) / castx.ToFloat64(bo.pageStep))
		i := castx.ToInt(n1)
		j := 0

		for k := bo.pageStep * (i - 1) + 1; k <= pageTotal; k++ {
			if j > bo.pageStep {
				break
			}

			pageList = append(pageList, k)
			j++
		}
	}

	map1 := map[string]interface{}{
		"recordTotal": bo.recordTotal,
		"pageTotal":   pageTotal,
		"currentPage": bo.currentPage,
		"pageSize":    bo.pageSize,
		"pageStep":    bo.pageStep,
		"pageList":    pageList,
	}

	for key := range map1 {
		if !slicex.InStringSlice(key, _includeFields) {
			delete(map1, key)
		}
	}

	return map1
}

func (bo *pager) ToCommonMap() map[string]interface{} {
	return bo.ToMap("recordTotal, pageTotal, currentPage, pageSize")
}
