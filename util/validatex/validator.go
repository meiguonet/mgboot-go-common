package validatex

import (
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	dateFormatFull     = "2006-01-02 15:04:05"
	dateFormatDateOnly = "2006-01-02"
)

var daysOfMonth = [13]int{0, 31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

var regexpDatetime = struct {
	re1 *regexp.Regexp
	re2 *regexp.Regexp
	re3 *regexp.Regexp
}{
	re1: regexp.MustCompile(`^\d{4}-\d{1,2}-\d{1,2}$`),
	re2: regexp.MustCompile(`^\d{4}/\d{1,2}/\d{1,2}$`),
	re3: regexp.MustCompile(`^\d{2}:\d{2}:\d{2}$`),
}

var regexpInt = regexp.MustCompile("^-?[0-9]+$")
var regexpCommaSep = regexp.MustCompile(`[\x20\t]*,[\x20\t]*`)
var regexpNumbers = regexp.MustCompile("^[0-9]+$")
var regexpAlphas = regexp.MustCompile("^[A-Za-z]+$")
var regexpAlnum = regexp.MustCompile("^[A-Za-z0-9]+$")

var regexpMobileNumber = struct {
	re1 *regexp.Regexp
}{re1: regexp.MustCompile("^[1-9][0-9]+$")}

var regexpEmail = struct {
	re1 *regexp.Regexp
	re2 *regexp.Regexp
}{
	re1: regexp.MustCompile(`^[A-Za-z0-9]]+`),
	re2: regexp.MustCompile(`\.[A-Za-z]{2,}$`),
}

var idcard = struct {
	re1 *regexp.Regexp
	a1  []int
	a2  []string
}{
	re1: regexp.MustCompile(`^[1-9]\d{5}(18|19|20)\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$`),
	a1:  []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2},
	a2:  []string{"1", "0", "x", "9", "8", "7", "6", "5", "4", "3", "2"},
}

type validator struct {
	checkers []RuleChecker
}

func NewValidator() *validator {
	return &validator{}
}

func (v *validator) WithCheckers(checkers []RuleChecker) {
	if len(checkers) < 1 {
		return
	}

	for _, checker := range checkers {
		v.WithChecker(checker)
	}
}

func (v *validator) WithChecker(checker RuleChecker) {
	if checker.GetRuleName() == "" {
		return
	}

	if v.checkers == nil {
		v.checkers = make([]RuleChecker, 0)
	}

	var found bool

	for _, ck := range v.checkers {
		if ck.GetRuleName() == checker.GetRuleName() {
			found = true
			break
		}
	}

	if !found {
		v.checkers = append(v.checkers, checker)
	}
}

func (v *validator) IsDate(arg0 string) bool {
	if arg0 == "" {
		return false
	}
	
	if !regexpDatetime.re1.MatchString(arg0) && !regexpDatetime.re2.MatchString(arg0) {
		return false
	}

	nums := strings.Split(strings.ReplaceAll(arg0, "/", "-"), "-")
	year := -1
	month := -1
	day := -1

	for _, num := range nums {
		n1, err := strconv.Atoi(num)

		if err != nil || n1 < 0 {
			continue
		}

		if year < 1 {
			year = n1
			continue
		}

		if month < 1 {
			month = n1
			continue
		}

		if day < 1 {
			day = n1
		}
	}

	if year < 1000 || month < 1 || month > 12 || day < 1 || day > 31 {
		return false
	}

	if day > daysOfMonth[month] {
		return false
	}

	isLeapYear := (year % 4 == 0 && year % 100 != 0) || year % 400 == 0

	if month == 2 && !isLeapYear && day > 28 {
		return false
	}

	return true
}

func (v *validator) IsDateTime(arg0 string) bool {
	if arg0 == "" {
		return false
	}

	parts := strings.Split(arg0,  " ")

	if len(parts) != 2 {
		return false
	}

	if !v.IsDate(parts[0]) {
		return false
	}

	if !regexpDatetime.re3.MatchString(parts[1]) {
		return false
	}

	nums := strings.Split(parts[1], ":")
	hours := -1
	minutes := -1
	seconds := -1

	for _, num := range nums {
		n1, err := strconv.Atoi(num)

		if err != nil || n1 < 0 {
			continue
		}

		if hours < 1 {
			hours = n1
			continue
		}

		if minutes < 1 {
			minutes = n1
			continue
		}

		if seconds < 1 {
			seconds = n1
		}
	}

	if hours < 0 || hours > 23 || minutes < 0 || minutes > 59 || seconds < 0 || seconds > 59 {
		return false
	}

	return true
}

func (v *validator) IsFutureDate(arg0 string) bool {
	if !v.IsDate(arg0) {
		return false
	}

	parts := strings.Split(strings.ReplaceAll(arg0, "/", "-"), "-")

	if len(parts) < 3 {
		return false
	}

	year, err := strconv.Atoi(parts[0])

	if err != nil {
		return false
	}

	month, err := strconv.Atoi(parts[0])

	if err != nil {
		return false
	}

	date, err := strconv.Atoi(parts[0])

	if err != nil {
		return false
	}

	now := time.Now()
	t1 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	t2, _ := time.ParseInLocation(
		dateFormatDateOnly, fmt.Sprintf("%d-%02d-%02d", year, month, date), now.Location(),
	)

	return t1.After(t2)
}

func (v *validator) IsPastDate(arg0 string) bool {
	parts := strings.Split(strings.ReplaceAll(arg0, "/", "-"), "-")

	if len(parts) < 3 {
		return false
	}

	year, err := strconv.Atoi(parts[0])

	if err != nil {
		return false
	}

	month, err := strconv.Atoi(parts[1])

	if err != nil {
		return false
	}

	date, err := strconv.Atoi(parts[2])

	if err != nil {
		return false
	}

	now := time.Now()
	t1 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	t2, _ := time.ParseInLocation(
		dateFormatDateOnly, fmt.Sprintf("%d-%02d-%02d", year, month, date), now.Location(),
	)

	return t2.After(t1)
}

func (v *validator) IsInt(arg0 string) bool {
	if arg0 == "" {
		return false
	}

	return regexpInt.MatchString(arg0)
}

func (v *validator) IsIntEq(arg0, arg1 string, not ...bool) bool {
	if !v.IsInt(arg0) || !v.IsInt(arg1) {
		return false
	}

	if len(not) > 0 && not[0] {
		return v.toInt(arg0) != v.toInt(arg1)
	}

	return v.toInt(arg0) == v.toInt(arg1)
}

func (v *validator) IsIntNe(arg0, arg1 string) bool {
	return v.IsIntEq(arg0, arg1, true)
}

func (v *validator) IsIntGt(arg0, arg1 string) bool {
	if !v.IsInt(arg0) || !v.IsInt(arg1) {
		return false
	}

	return v.toInt(arg0) > v.toInt(arg1)
}

func (v *validator) IsIntGe(arg0, arg1 string) bool {
	if !v.IsInt(arg0) || !v.IsInt(arg1) {
		return false
	}

	return v.toInt(arg0) >= v.toInt(arg1)
}

func (v *validator) IsIntLt(arg0, arg1 string) bool {
	if !v.IsInt(arg0) || !v.IsInt(arg1) {
		return false
	}

	return v.toInt(arg0) < v.toInt(arg1)
}

func (v *validator) IsIntLe(arg0, arg1 string) bool {
	if !v.IsInt(arg0) || !v.IsInt(arg1) {
		return false
	}

	return v.toInt(arg0) <= v.toInt(arg1)
}

func (v *validator) IsIntBetween(arg0, arg1 string) bool {
	if !v.IsInt(arg0) {
		return false
	}

	parts := regexpCommaSep.Split(arg1, -1)

	if len(parts) != 2 {
		return false
	}

	p1 := strings.TrimSpace(parts[0])
	p2 := strings.TrimSpace(parts[1])

	if !v.IsInt(p1) || !v.IsInt(p2) {
		return false
	}

	n1 := v.toInt(arg1)
	return n1 >= v.toInt(p1) && n1 <= v.toInt(p2)
}

func (v *validator) IsIntIn(arg0, arg1 string, not ...bool) bool {
	if !v.IsInt(arg0) {
		return false
	}

	n1 := v.toInt(arg0)
	parts := regexpCommaSep.Split(arg1, -1)

	if len(not) > 0 && not[0] {
		flag := true

		for _, p := range parts {
			n2, err := strconv.Atoi(strings.TrimSpace(p))

			if err != nil {
				continue
			}

			if n1 == n2 {
				flag = false
				break
			}
		}

		return flag
	}

	for _, p := range parts {
		n2, err := strconv.Atoi(strings.TrimSpace(p))

		if err != nil {
			continue
		}

		if n1 == n2 {
			return true
		}
	}

	return false
}

func (v *validator) IsNotIntIn(arg0, arg1 string) bool {
	return v.IsIntIn(arg0, arg1, true)
}

func (v *validator) IsFloat(arg0 string) bool {
	if arg0 == "" {
		return false
	}

	if v.IsInt(arg0) {
		return true
	}

	if !strings.Contains(arg0, ".") {
		return false
	}

	parts := strings.Split(arg0, ".")

	if len(parts) < 2 || !v.IsInt(parts[0]) {
		return false
	}

	return regexpNumbers.MatchString(parts[1])
}

func (v *validator) IsFloatEq(arg0, arg1 string, not ...bool) bool {
	if !v.IsFloat(arg0) || !v.IsFloat(arg1) {
		return false
	}

	n1 := v.toBigFloat(arg0)
	n2 := v.toBigFloat(arg1)

	if len(not) > 0 && not[0] {
		return n1.Cmp(n2) != 0
	}

	return n1.Cmp(n2) == 0
}

func (v *validator) IsFloatNe(arg0, arg1 string) bool {
	return v.IsFloatEq(arg0, arg1, true)
}

func (v *validator) IsFloatGt(arg0, arg1 string) bool {
	if !v.IsFloat(arg0) || !v.IsFloat(arg1) {
		return false
	}

	n1 := v.toBigFloat(arg0)
	n2 := v.toBigFloat(arg1)
	return n1.Cmp(n2) == 1
}

func (v *validator) IsFloatGe(arg0, arg1 string) bool {
	if !v.IsFloat(arg0) || !v.IsFloat(arg1) {
		return false
	}

	n1 := v.toBigFloat(arg0)
	n2 := v.toBigFloat(arg1)
	return n1.Cmp(n2) != -1
}

func (v *validator) IsFloatLt(arg0, arg1 string) bool {
	if !v.IsFloat(arg0) || !v.IsFloat(arg1) {
		return false
	}

	n1 := v.toBigFloat(arg0)
	n2 := v.toBigFloat(arg1)
	return n1.Cmp(n2) == -1
}

func (v *validator) IsFloatLe(arg0, arg1 string) bool {
	if !v.IsInt(arg0) || !v.IsInt(arg1) {
		return false
	}

	n1 := v.toBigFloat(arg0)
	n2 := v.toBigFloat(arg1)
	return n1.Cmp(n2) != 1
}

func (v *validator) IsFloatBetween(arg0, arg1 string) bool {
	if !v.IsFloat(arg0) {
		return false
	}

	parts := regexpCommaSep.Split(arg1, -1)

	if len(parts) != 2 {
		return false
	}

	p1 := strings.TrimSpace(parts[0])
	p2 := strings.TrimSpace(parts[1])

	if !v.IsFloat(p1) || !v.IsFloat(p2) {
		return false
	}

	n0 := v.toBigFloat(arg0)
	n1 := v.toBigFloat(p1)
	n2 := v.toBigFloat(p2)
	return n0.Cmp(n1) != -1 && n0.Cmp(n2) != 1
}

func (v *validator) IsStrEq(arg0, arg1 string, ignoreCase ...bool) bool {
	if len(ignoreCase) > 0 && ignoreCase[0] {
		return strings.ToLower(arg0) == strings.ToLower(arg1)
	}

	return arg0 == arg1
}

func (v *validator) IsStrEqI(arg0, arg1 string) bool {
	return v.IsStrEq(arg0, arg1, true)
}

func (v *validator) IsStrNe(arg0, arg1 string, ignoreCase ...bool) bool {
	if len(ignoreCase) > 0 && ignoreCase[0] {
		return strings.ToLower(arg0) != strings.ToLower(arg1)
	}

	return arg0 != arg1
}

func (v *validator) IsStrNeI(arg0, arg1 string) bool {
	return v.IsStrNe(arg0, arg1, true)
}

func (v *validator) IsStrIn(arg0, arg1 string, ignoreCase ...bool) bool {
	parts := regexpCommaSep.Split(arg1, -1)

	if len(ignoreCase) > 0 && ignoreCase[0] {
		arg0 = strings.ToLower(arg0)
		flag := true

		for _, p := range parts {
			if arg0 == strings.ToLower(p) {
				flag = false
				break
			}
		}

		return flag
	}

	for _, p := range parts {
		if arg0 == p {
			return true
		}
	}

	return false
}

func (v *validator) IsStrInI(arg0, arg1 string) bool {
	return v.IsStrIn(arg0, arg1, true)
}

func (v *validator) IsStrNotIn(arg0, arg1 string) bool {
	return !v.IsStrIn(arg0, arg1)
}

func (v *validator) IsStrNotInI(arg0, arg1 string) bool {
	return !v.IsStrIn(arg0, arg1, true)
}

func (v *validator) IsStrLen(arg0, arg1 string) bool {
	if !v.IsInt(arg1) {
		return false
	}

	return len(arg0) == v.toInt(arg1)
}

func (v *validator) IsStrLenGt(arg0, arg1 string) bool {
	if !v.IsInt(arg1) {
		return false
	}

	return len(arg0) > v.toInt(arg1)
}

func (v *validator) IsStrlenGe(arg0, arg1 string) bool {
	if !v.IsInt(arg1) {
		return false
	}

	return len(arg0) >= v.toInt(arg1)
}

func (v *validator) IsStrlenLt(arg0, arg1 string) bool {
	if !v.IsInt(arg1) {
		return false
	}

	return len(arg0) < v.toInt(arg1)
}

func (v *validator) IsStrLenLe(arg0, arg1 string) bool {
	if !v.IsInt(arg1) {
		return false
	}

	return len(arg0) <= v.toInt(arg1)
}

func (v *validator) IsStrLenBetween(arg0, arg1 string) bool {
	parts := regexpCommaSep.Split(arg1, -1)

	if len(parts) != 2 {
		return false
	}

	p1 := strings.TrimSpace(parts[0])
	p2 := strings.TrimSpace(parts[1])

	if !v.IsInt(p1) || !v.IsInt(p2) {
		return false
	}

	n1 := len(arg0)
	return n1 >= v.toInt(p1) && n1 <= v.toInt(p2)
}

func (v *validator) IsAlphas(arg0 string) bool {
	if arg0 == "" {
		return false
	}

	return regexpAlphas.MatchString(arg0)
}

func (v *validator) IsNumbers(arg0 string) bool {
	if arg0 == "" {
		return false
	}

	return regexpNumbers.MatchString(arg0)
}

func (v *validator) IsAlnum(arg0 string) bool {
	if arg0 == "" {
		return false
	}

	return regexpAlnum.MatchString(arg0)
}

func (v *validator) IsMobile(arg0 string) bool {
	if arg0 == "" {
		return false
	}

	n1 := len(arg0)

	if n1 < 11 || n1 > 16 {
		return false
	}

	return regexpMobileNumber.re1.MatchString(arg0)
}

func (v *validator) IsEmail(arg0 string) bool {
	if arg0 == "" || !strings.Contains(arg0, "@") {
		return false
	}

	parts := strings.Split(arg0, "@")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return false
	}

	p2 := parts[1]

	if !strings.Contains(p2, ".") || strings.HasPrefix(p2, ".") || strings.HasSuffix(p2, ".") {
		return false
	}

	return regexpEmail.re1.MatchString(p2) && regexpEmail.re2.MatchString(p2)
}

func (v *validator) IsPasswordTooSimple(arg0 string) bool {
	if arg0 == "" {
		return false
	}

	r1 := 'A'
	r2 := 'Z'
	r3 := 'a'
	r4 := 'z'
	r5 := '0'
	r6 := '9'
	n1 := 0
	n2 := 0
	n3 := 0
	n4 := 0

	strings.Map(func(r rune) rune {
		if r >= r1 && r <= r2 {
			n1 = 1
			return r
		}

		if r >= r3 && r <= r4 {
			n2 = 1
			return r
		}

		if r >= r5 && r <= r6 {
			n3 = 1
			return r
		}

		n4 = 1
		return r
	}, arg0)

	return n1 + n2 + n3 + n4 >= 3
}

func (v *validator) IsIdcard(arg0 string) bool {
	if arg0 == "" {
		return false
	}

	if !idcard.re1.MatchString(arg0) {
		return false
	}

	sum := 0
	n1 := len(arg0) - 1

	for i := 0; i < n1; i++ {
		sum += v.toInt(arg0[i:i+1]) * idcard.a1[i]
	}

	return strings.ToLower(idcard.a2[sum % 11]) == strings.ToLower(arg0[n1:])
}

func (v *validator) IsRegexp(arg0, arg1 string) bool {
	if arg0 == "" || arg1 == "" {
		return false
	}

	re, err := regexp.Compile(arg1)

	if err != nil {
		return false
	}

	return re.MatchString(arg0)
}

func (v *validator) toInt(arg0 string) int {
	n1, _ := strconv.Atoi(arg0)
	return n1
}

func (v *validator) toBigFloat(arg0 string) *big.Float {
	n1, _ := strconv.ParseFloat(arg0, 64)
	return big.NewFloat(n1)
}
