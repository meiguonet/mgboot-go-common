package validatex

import (
	"github.com/meiguonet/mgboot-go-common/util/castx"
	"github.com/meiguonet/mgboot-go-common/util/stringx"
	"github.com/meiguonet/mgboot-go-common/util/structx"
	"reflect"
	"strings"
)

func Validate(validator *validator, data map[string]interface{}, rules []string, autoPanic ...bool) map[string]string {
	if validator == nil {
		validator = NewValidator()
	}

	errors := map[string]string{}

	if len(data) < 1 || len(rules) < 1 {
		return errors
	}

	handleData(data)

	for _, rule := range rules {
		var checkOnNotEmpty bool
		fn := "Required"
		var checkValue, msg, fieldName string

		if strings.Contains(rule, "@CheckOnNotEmpty") || strings.Contains(rule, "@WithNotEmpty") {
			checkOnNotEmpty = true
			rule = strings.ReplaceAll(rule, "@CheckOnNotEmpty", "")
			rule = strings.ReplaceAll(rule, "@WithNotEmpty", "")
		}

		if strings.Contains(rule, "@msg:") {
			msg = stringx.SubstringAfterLast(rule, "@")
			msg = strings.TrimPrefix(msg, "msg:")
			msg = strings.TrimSpace(msg)
			rule = stringx.SubstringBeforeLast(rule, "@")
		}

		if strings.Contains(rule, "@") {
			fieldName = stringx.SubstringBefore(rule, "@")
			fieldName = strings.TrimSpace(fieldName)
			fn = stringx.SubstringAfter(rule, "@")
			fn = strings.TrimSpace(fn)

			if strings.Contains(fn, ":") {
				checkValue = stringx.SubstringAfter(fn, ":")
				checkValue = strings.TrimSpace(checkValue)
				fn = stringx.SubstringBefore(fn, ":")
				fn = strings.TrimSpace(fn)
			}
		} else {
			fieldName = strings.TrimSpace(rule)
		}

		if msg == "" {
			switch fn {
			case "Mobile":
				msg = "不是有效的手机号码"
				break
			case "Email":
				msg = "不是有效的邮箱地址"
				break
			case "PasswordTooSimple":
				msg = "密码过于简单"
				break
			case "Idcard":
				msg = "不是有效的身份证号码"
				break
			default:
				msg = "必须填写"
				break
			}
		}

		fieldValue := getMapFieldValue(data, fieldName)

		if checkOnNotEmpty && fieldValue == "" {
			continue
		}

		if fieldValue == "" {
			errors[fieldName] = msg
			continue
		}

		if fn == "Required" {
			continue
		}

		if fn == "EqualsWith" {
			if fieldValue != getMapFieldValue(data, checkValue) {
				errors[fieldName] = msg
			}

			continue
		}

		funcName := "Is" + fn
		var passed bool

		if checkValue == "" {
			passed = validateByFn(validator, funcName, fieldValue)
		} else {
			passed = validateByFn(validator, funcName, fieldValue, checkValue)
		}

		if len(validator.checkers) > 0 {
			for _, checker := range validator.checkers {
				if checker.GetRuleName() != fn {
					continue
				}

				if checkValue == "" {
					passed = checker.Check(fieldValue)
				} else {
					passed = checker.Check(fieldValue, checkValue)
				}

				break
			}
		}

		if passed {
			continue
		}

		errors[fieldName] = msg
	}

	if len(autoPanic) > 0 && autoPanic[0] && len(errors) > 0 {
		panic(NewValidateExceptionWithErrors(errors))
	}

	return errors
}

func ValidateByStruct(validator *validator, arg0 interface{}, rules []string, autoPanic ...bool) map[string]string {
	return Validate(validator, structx.ToMap(arg0), rules, autoPanic...)
}

func FailfastValidate(validator *validator, data map[string]interface{}, rules []string, autoPanic ...bool) string {
	if validator == nil {
		validator = NewValidator()
	}

	var errorTips string

	if len(data) < 1 || len(rules) < 1 {
		return errorTips
	}

	handleData(data)

	for _, rule := range rules {
		var checkOnNotEmpty bool
		fn := "Required"
		var checkValue, msg, fieldName string

		if strings.Contains(rule, "@CheckOnNotEmpty") || strings.Contains(rule, "@WithNotEmpty") {
			checkOnNotEmpty = true
			rule = strings.ReplaceAll(rule, "@CheckOnNotEmpty", "")
			rule = strings.ReplaceAll(rule, "@WithNotEmpty", "")
		}

		if strings.Contains(rule, "@msg:") {
			msg = stringx.SubstringAfterLast(rule, "@")
			msg = strings.TrimPrefix(msg, "msg:")
			msg = strings.TrimSpace(msg)
			rule = stringx.SubstringBeforeLast(rule, "@")
		}

		if strings.Contains(rule, "@") {
			fieldName = stringx.SubstringBefore(rule, "@")
			fieldName = strings.TrimSpace(fieldName)
			fn = stringx.SubstringAfter(rule, "@")
			fn = strings.TrimSpace(fn)

			if strings.Contains(fn, ":") {
				checkValue = stringx.SubstringAfter(fn, ":")
				checkValue = strings.TrimSpace(checkValue)
				fn = stringx.SubstringBefore(fn, ":")
				fn = strings.TrimSpace(fn)
			}
		} else {
			fieldName = strings.TrimSpace(rule)
		}

		if msg == "" {
			switch fn {
			case "Mobile":
				msg = "不是有效的手机号码"
				break
			case "Email":
				msg = "不是有效的邮箱地址"
				break
			case "PasswordTooSimple":
				msg = "密码过于简单"
				break
			case "Idcard":
				msg = "不是有效的身份证号码"
				break
			default:
				msg = "必须填写"
				break
			}
		}

		fieldValue := getMapFieldValue(data, fieldName)

		if checkOnNotEmpty && fieldValue == "" {
			continue
		}

		if fieldValue == "" {
			errorTips = msg
			break
		}

		if fn == "Required" {
			continue
		}

		if fn == "EqualsWith" {
			if fieldValue != getMapFieldValue(data, checkValue) {
				errorTips = msg
				break
			}

			continue
		}

		funcName := "Is" + fn
		var passed bool

		if checkValue == "" {
			passed = validateByFn(validator, funcName, fieldValue)
		} else {
			passed = validateByFn(validator, funcName, fieldValue, checkValue)
		}

		if len(validator.checkers) > 0 {
			for _, checker := range validator.checkers {
				if checker.GetRuleName() != fn {
					continue
				}

				if checkValue == "" {
					passed = checker.Check(fieldValue)
				} else {
					passed = checker.Check(fieldValue, checkValue)
				}

				break
			}
		}

		if !passed {
			errorTips = msg
			break
		}
	}

	if len(autoPanic) > 0 && autoPanic[0] && errorTips != "" {
		panic(NewValidateExceptionWithErrorTips(errorTips))
	}

	return errorTips
}

func FailfastValidateByStruct(validator *validator, arg0 interface{}, rules []string, autoPanic ...bool) string {
	return FailfastValidate(validator, structx.ToMap(arg0), rules, autoPanic...)
}

func handleData(data map[string]interface{}) {
	for key, value := range data {
		newKey := strings.ToLower(key)

		if key == newKey {
			continue
		}

		data[newKey] = value
		delete(data, key)
	}
}

func getMapFieldValue(data map[string]interface{}, key string) string {
	if _, ok := data[key]; ok {
		return castx.ToString(data[key])
	}

	key = strings.ToLower(key)

	if _, ok := data[key]; ok {
		return castx.ToString(data[key])
	}

	for mapKey, value := range data {
		if key == strings.ReplaceAll(mapKey, "_", "") {
			return castx.ToString(value)
		}
	}

	return ""
}

func validateByFn(validator *validator, funcName string, args ...string) bool {
	rv := reflect.ValueOf(validator)
	method := rv.MethodByName(funcName)

	if !method.IsValid() {
		return true
	}

	callArgs := make([]reflect.Value, 0)

	for _, arg := range args {
		callArgs = append(callArgs, reflect.ValueOf(arg))
	}

	results := method.Call(callArgs)

	if len(results) < 1 {
		return false
	}

	flag, ok := results[0].Interface().(bool)

	if !ok {
		return false
	}

	return flag
}
