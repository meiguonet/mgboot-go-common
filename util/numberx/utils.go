package numberx

import (
	"fmt"
	"math/big"
	"meiguonet/mgboot-go-common/util/castx"
	"meiguonet/mgboot-go-common/util/stringx"
	"reflect"
	"strings"
)

func IsNumber(arg0 interface{}) bool {
	kinds := []reflect.Kind{
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128,
	}

	kind := reflect.TypeOf(arg0).Kind()

	for _, v := range kinds {
		if kind == v {
			return true
		}
	}

	return false
}

func ToDecimalString(num interface{}, args ...interface{}) string {
	fractionDigitsNum := 2
	var stripTailsZero interface{}

	for _, arg := range args {
		if n1, ok := arg.(int); ok {
			if n1 >= 0 {
				fractionDigitsNum = n1
			}

			continue
		}

		if b1, ok := arg.(bool); ok {
			if _, isBool := stripTailsZero.(bool); !isBool {
				stripTailsZero = b1
			}
		}
	}

	if fractionDigitsNum > 12 {
		fractionDigitsNum = 12
	}

	if _, ok := stripTailsZero.(bool); !ok {
		stripTailsZero = true
	}

	str := fmt.Sprintf("%0.12f", castx.ToFloat64(num))
	parts := strings.Split(str, ".")
	p1 := parts[0]

	if fractionDigitsNum < 1 {
		return p1
	}

	p2 := parts[1][:fractionDigitsNum]

	if yes, ok := stripTailsZero.(bool); ok && yes {
		p2 = strings.TrimSuffix(p2, "0")
	}

	if p2 == "" {
		return p1
	}

	return p1 + "." + p2
}

func ToDecimalStringWithThousandSep(num interface{}, args ...interface{}) string {
	s1 := ToDecimalString(num, args...)

	if !strings.Contains(s1, ".") {
		return thousandSep(s1)
	}

	parts := strings.Split(s1, ".")
	return thousandSep(parts[0]) + "." + parts[1]
}

func ToFriendlyString(num interface{}, args ...interface{}) string {
	n1, err := castx.ToFloat64E(num)

	if err != nil {
		return ""
	}

	fractionDigitsNum := 2
	units := []string{"K", "W"}

	for _, arg := range args {
		if n1, ok := arg.(int); ok {
			if n1 > 0 {
				fractionDigitsNum = n1
			}

			continue
		}

		if a1, ok := arg.([]string); ok && len(a1) == 2 {
			units = a1
		}
	}

	n2 := big.NewFloat(n1)

	if n2.Cmp(big.NewFloat(1000.0)) == -1 {
		return stringx.SubstringBefore(ToDecimalString(n1), ".")
	}

	if n2.Cmp(big.NewFloat(10000.0)) == -1 {
		n1, _ = n2.Quo(n2, big.NewFloat(1000.0)).Float64()
		return ToDecimalString(n1, fractionDigitsNum) + units[0]
	}

	n1, _ = n2.Quo(n2, big.NewFloat(10000.0)).Float64()
	return ToDecimalString(n1, fractionDigitsNum) + units[1]
}

func ToFriendlyDistanceString(num interface{}, units ...[]string) string {
	n1, err := castx.ToFloat64E(num)

	if err != nil {
		return ""
	}

	_units := []string{"m", "km"}

	if len(units) > 0 && len(units[0]) == 2 {
		_units = units[0]
	}

	n2 := big.NewFloat(n1)

	if n2.Cmp(big.NewFloat(1000.0)) == -1 {
		return stringx.SubstringBefore(ToDecimalString(n1), ".") + _units[0]
	}

	n1, _ = n2.Quo(n2, big.NewFloat(1000.0)).Float64()
	return ToDecimalString(n1, 1) + _units[1]
}

func Ojld(m, n int) int {
	if m == 0 || n == 0 {
		return 1
	}

	if m == n {
		return m
	}

	for {
		if n == 0 {
			break
		}

		m, n = n, m % n
	}

	return m
}

func thousandSep(str string) string {
	if str == "" || len(str) <= 3 {
		return str
	}

	n1 := len(str) % 3
	n2 := len(str) / 3

	if n1 == 0 {
		n1 = 3
		n2 -= 1
	}

	sb := strings.Builder{}
	sb.WriteString(str[0:n1])
	str = str[n1:]

	for i := 1; i <= n2; i++ {
		sb.WriteString(",")

		if i == n2 {
			sb.WriteString(str[(i-1)*3:])
		} else {
			sb.WriteString(str[(i-1)*3:i*3])
		}
	}

	return sb.String()
}
