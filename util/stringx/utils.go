package stringx

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/meiguonet/mgboot-go-common/enum/DatetimeFormat"
	"github.com/meiguonet/mgboot-go-common/enum/RandomStringType"
	"github.com/meiguonet/mgboot-go-common/util/castx"
	"github.com/meiguonet/mgboot-go-common/util/slicex"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io/ioutil"
	mrand "math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	format1 = "%s-%s-%s"
	regex1  = `^\d{4}-\d{1,2}-\d{1,2}$`
	regex2  = `^\d{4}/\d{1,2}/\d{1,2}$`
	regex3  = `^\d{2}:\d{2}:\d{2}$`
)

var daysOfMonth = [13]int{0, 31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

type customRandSource struct{}

func (s *customRandSource) Seed(seed int64) {
	_ = seed
}

func (s *customRandSource) Uint64() uint64 {
	var n1 uint64
	binary.Read(crand.Reader, binary.BigEndian, &n1)
	return n1
}

func (s *customRandSource) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

func RegexMatch(str, pattern string) bool {
	re, err := regexp.Compile(pattern)

	if err != nil {
		return false
	}

	return re.MatchString(str)
}

func RegexReplace(str, pattern, repl string) string {
	re, err := regexp.Compile(pattern)

	if err != nil {
		return str
	}

	return re.ReplaceAllString(str, repl)
}

func IsNumeric(str string) bool {
	return RegexMatch(str, "^[0-9]+$")
}

func IsLetteric(str string) bool {
	return RegexMatch(str, "^[A-Za-z]+$")
}

func IsAlnum(str string) bool {
	return RegexMatch(str, "^[A-Za-z0-9]+$")
}

func IsInt(str string) bool {
	return RegexMatch(str, "^-?[0-9]+$")
}

func IsFloat(str string) bool {
	if IsInt(str) {
		return true
	}

	if !strings.Contains(str, ".") {
		return false
	}

	p1 := SubstringBefore(str, ".")

	if !IsInt(p1) {
		return false
	}

	return RegexMatch(SubstringAfter(str, "."), "^[0-9]+$")
}

func IsEmail(str string) bool {
	if !strings.Contains(str, "@") {
		return false
	}

	p1 := SubstringBefore(str, "@")

	if p1 == "" {
		return false
	}

	p2 := SubstringAfter(str, "@")

	if p2 == "" ||
		strings.Contains(p2, "@") ||
		!strings.Contains(p2, ".") ||
		strings.HasPrefix(p2, ".") ||
		strings.HasSuffix(p2, ".") {
		return false
	}

	return RegexMatch(SubstringAfterLast(p2, "."), "^[A-Za-z]{2,}$")
}

func IsNatiioalPhoneNumber(str string) bool {
	str = strings.ReplaceAll(str, "（", "(")
	str = strings.ReplaceAll(str, "）", ")")
	str = strings.ReplaceAll(str, "－", "-")
	str = strings.ReplaceAll(str, "—", "-")

	if RegexMatch(str, `^\d{7,}$`) {
		return true
	}

	if RegexMatch(str, `^\d{3,4}[\x20\t]*-[\x20\t]*\d{7,}$`) {
		return true
	}

	if RegexMatch(str, `^\d{3,4}[\x20\t]*-[\x20]\t*\d{7,}[\x20\t]*-[\x20\t]*\d{1,5}$`) {
		return true
	}

	if RegexMatch(str, `^\([\x20\t]*\d{3,4}[\x20\t]*\)[\x20\t]*\d{7,}$`) {
		return true
	}

	return RegexMatch(str, `^([\x20\t]*\d{3,4}[\x20\t]*\)[\x20\t]*\d{7,}[\x20\t]*-[\x20\t]*\d{1,5}$`)
}

func IsNationalMobileNumber(str string) bool {
	n1 := len(str)

	if n1 < 11 || n1 > 16 {
		return false
	}

	return RegexMatch(str, "^[1-9][0-9]+$")
}

func IsDate(str string) bool {
	if !RegexMatch(str, regex1) && !RegexMatch(str, regex2) {
		return false
	}

	a1 := strings.Split(strings.ReplaceAll(str, "/", "-"), "-")
	year := -1
	month := -1
	day := -1

	for _, s1 := range a1 {
		n1, err := castx.ToIntE(s1)

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

func IsTime(str string) bool {
	if !RegexMatch(str, regex3) {
		return false
	}

	a1 := strings.Split(str, ":")
	hours := -1
	minutes := -1
	seconds := -1

	for _, s1 := range a1 {
		n1, err := castx.ToIntE(s1)

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

func IsDatetime(str string) bool {
	parts := strings.Split(str,  " ")
	return len(parts) == 2 && IsDate(parts[0]) && IsTime(parts[1])
}

func Ucfirst(str string) string {
	if str == "" {
		return ""
	}

	if len(str) < 2 {
		return strings.ToUpper(str)
	}

	return strings.ToUpper(str[:1]) + str[1:]
}

func Lcfirst(str string) string {
	if str == "" {
		return ""
	}

	if len(str) < 2 {
		return strings.ToLower(str)
	}

	return strings.ToLower(str[:1]) + str[1:]
}

func Ucwords(str, joinBy string, delimiters ...string) string {
	if str == "" {
		return str
	}

	chars := make([]rune, 0)

	for _, s1 := range delimiters {
		if len(s1) != 1 {
			continue
		}

		chars = append(chars, []rune(s1)[0])
	}

	if len(chars) < 1 {
		chars = []rune{' '}
	}

	parts := strings.FieldsFunc(str, func(r rune) bool {
		var matched bool

		for _, ch := range chars {
			if r == ch {
				matched = true
				break
			}
		}

		return matched
	})

	if len(parts) < 1 {
		return str
	}

	for i, p := range parts {
		parts[i] = Ucfirst(p)
	}

	return strings.Join(parts, joinBy)
}

func Lcwords(str, joinBy string, delimiters ...string) string {
	if str == "" {
		return str
	}

	chars := make([]rune, 0)

	for _, s1 := range delimiters {
		if len(s1) != 1 {
			continue
		}

		chars = append(chars, []rune(s1)[0])
	}

	if len(chars) < 1 {
		chars = []rune{' '}
	}

	parts := strings.FieldsFunc(str, func(r rune) bool {
		var matched bool

		for _, ch := range chars {
			if r == ch {
				matched = true
				break
			}
		}

		return matched
	})

	if len(parts) < 1 {
		return str
	}

	for i, p := range parts {
		parts[i] = Ucfirst(p)
	}

	return strings.Join(parts, joinBy)
}

func ToTime(str string) time.Time {
	t1 := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.Local)

	if str == "" {
		return t1
	}

	s1 := strings.TrimSpace(str)
	s1 = strings.ReplaceAll(s1, "/", "-")
	s1 = ReplaceWithRegexp(s1, `[\x20\t]+`, " ")
	var a1, a2 []string

	if strings.Contains(s1, " ") {
		a1 = strings.Split(SubstringBefore(s1, " "), "-")
		a2 = strings.Split(SubstringAfter(s1, " "), ":")
	} else {
		a1 = strings.Split(s1, "-")
	}

	year := -1
	month := -1
	day := -1

	for _, s2 := range a1 {
		s2 = strings.TrimSpace(s2)
		n1, err := castx.ToIntE(s2)

		if err != nil || n1 < 0 {
			continue
		}

		if year < 0 {
			year = n1
			continue
		}

		if month < 0 {
			month = n1
			continue
		}

		if day < 0 {
			day = n1
		}
	}

	if year < 1000 || month < 1 || month > 12 || day < 1 || day > 31 {
		return t1
	}

	isLeapYear := (year % 4 == 0 && year % 100 != 0) || year % 400 == 0

	if month == 2 && !isLeapYear && day > 28 {
		return t1
	}

	hour := -1
	minute := -1
	second := -1

	for _, s2 := range a2 {
		s2 = strings.TrimSpace(s2)
		n1, err := castx.ToIntE(s2)

		if err != nil || n1 < 0 {
			continue
		}

		if hour < 0 {
			hour = n1
			continue
		}

		if minute < 0 {
			minute = n1
			continue
		}

		if second < 0 {
			second = n1
		}
	}

	if hour < 0 {
		hour = 0
	}

	if minute < 0 {
		minute = 0
	}

	if second < 0 {
		second = 0
	}

	if hour > 23 || minute > 59 || second > 59 {
		return t1
	}

	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)
}

func ToShortDate(str string) time.Time {
	year, month, day := ToTime(str).Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

func ToTimestamp(str string) int64 {
	if str == "" {
		return 0
	}

	return ToTime(str).Unix()
}

func ToCharArray(str string) []string {
	if str == "" {
		return []string{}
	}

	runes := []rune(str)
	chars := make([]string, len(runes))

	for i, r := range runes {
		chars[i] = string(r)
	}

	return chars
}

func SplitByChars(str string, chars ...rune) []string {
	if str == "" {
		return []string{}
	}

	if len(chars) < 1 {
		return []string{str}
	}

	return strings.FieldsFunc(str, func(r rune) bool {
		for _, ch := range chars {
			if ch == r {
				return true
			}
		}

		return false
	})
}

func SplitWithRegexp(str, pattern string) []string {
	if str == "" {
		return []string{}
	}

	re, err := regexp.Compile(pattern)

	if err != nil {
		return []string{}
	}

	return re.Split(str, -1)
}

func ReplaceWithRegexp(str, pattern, repl string) string {
	if str == "" {
		return ""
	}

	re, err := regexp.Compile(pattern)

	if err != nil {
		return str
	}

	return re.ReplaceAllString(str, repl)
}

func SubstringBefore(str, delimiter string, last ...bool) string {
	var idx int

	if len(last) > 0 && last[0] {
		idx = strings.LastIndex(str, delimiter)
	} else {
		idx = strings.Index(str, delimiter)
	}

	if idx < 1 {
		return ""
	}

	return str[:idx]
}

func SubstringBeforeLast(str, delimiter string) string {
	return SubstringBefore(str, delimiter, true)
}

func SubstringAfter(str, delimiter string, last ...bool) string {
	var idx int

	if len(last) > 0 && last[0] {
		idx = strings.LastIndex(str, delimiter)
	} else {
		idx = strings.Index(str, delimiter)
	}

	if idx < 0 {
		return ""
	}

	idx += len(delimiter)

	if idx >= len(str) {
		return ""
	}

	return str[idx:]
}

func SubstringAfterLast(str, delimiter string) string {
	return SubstringAfter(str, delimiter, true)
}

func Mask(str string, prefixLen, suffixLen int) string {
	if str == "" {
		return ""
	}

	if prefixLen < 1 {
		return str
	}

	chars := ToCharArray(str)
	n1 := len(chars)

	if n1 <= prefixLen {
		return str
	}

	p1 := strings.Join(chars[0:prefixLen], "")
	n2 := 0

	if suffixLen > 0 {
		n2 = suffixLen
	}

	if n2 < 1 {
		return p1 + strings.Repeat("*", n1 - prefixLen)
	}

	if suffixLen + n2 == n1 {
		return str
	}

	n3 := n1 - n2
	p2 := strings.Join(chars[n3:], "")
	return p1 + strings.Repeat("*", n1 - prefixLen - n2) + p2
}

func MaskEmail(str string) string {
	if str == "" || !IsEmail(str) {
		return str
	}

	p1 := SubstringBefore(str, "@")
	p2 := SubstringAfter(str, "@")

	switch utf8.RuneCountInString(p1) {
	case 1, 2:
		return Mask(p1, 1, 0) + "@" + p2
	case 3:
		return Mask(p1, 2, 0) + "@" + p2
	case 4:
		return Mask(p1, 2, 1) + "@" + p2
	case 5:
		return Mask(p1, 2, 2) + "@" + p2
	default:
		return Mask(p1, 3, 2) + "@" + p2
	}
}

func RemoveSqlSpecialChars(str string) string {
	if str == "" {
		return ""
	}

	str = strings.ReplaceAll(str, "\r", "")
	str = strings.ReplaceAll(str, "\n", "")
	str = strings.ReplaceAll(str, "\\\"", "")
	str = strings.ReplaceAll(str, "'", "")
	str = strings.ReplaceAll(str, "\"", "")
	return str
}

func EnsureLeft(str, prefix string) string {
	if str == "" || strings.HasPrefix(str, prefix) {
		return str
	}

	return prefix + str
}

func EnsureRight(str, suffix string) string {
	if str == "" || strings.HasSuffix(str, suffix) {
		return str
	}

	return str + suffix
}

func RegexpCatch(str, regex string) []string {
	re, err := regexp.Compile(regex)

	if err != nil {
		return []string{}
	}

	matches := re.FindStringSubmatch(str)

	if matches == nil {
		return []string{}
	}

	return matches
}

func RegexpCatchAll(str, regex string) [][]string {
	re, err := regexp.Compile(regex)

	if err != nil {
		return [][]string{}
	}

	matches := re.FindAllStringSubmatch(str, -1)

	if matches == nil {
		return [][]string{}
	}

	return matches
}

func GetFileExt(str string) string {
	if str == "" || !strings.Contains(str, ".") {
		return ""
	}

	return strings.ToLower(SubstringAfterLast(str, "."))
}

func ToDataSize(str string) int64 {
	if str == "" {
		return 0
	}

	regex1 := regexp.MustCompile("[^0-9]")
	digits := regex1.ReplaceAllString(str, "")

	if digits == "" {
		return 0
	}

	num := castx.ToInt64(digits)
	regex2 := regexp.MustCompile("[^A-Za-z]")
	unit := strings.ToLower(regex2.ReplaceAllString(str, ""))

	switch unit {
	case "":
		return num
	case "k":
		return num * 1024
	case "m":
		return num * 1024 * 1024
	case "g":
		return num * 1024 * 1024 * 1024
	default:
		return 0
	}
}

func ToDuration(str string) time.Duration {
	if str == "" {
		return 0
	}

	if d1, err := time.ParseDuration(str); err == nil {
		return d1
	}

	re := regexp.MustCompile(`[\x20\t]+`)
	str = re.ReplaceAllString(str, "")
	str = strings.TrimSpace(strings.ToLower(str))
	isNegative := strings.HasPrefix(str, "-")
	str = strings.TrimPrefix(str, "-")
	var d1 time.Duration
	re = regexp.MustCompile("^[1-9][0-9]*$")

	if re.MatchString(str) {
		n1 := castx.ToInt64(str)

		if isNegative {
			return d1 - time.Duration(n1) * time.Millisecond
		}

		return d1 + time.Duration(n1) * time.Millisecond
	}

	re = regexp.MustCompile("([1-9][0-9]*)d")
	groups := re.FindStringSubmatch(str)

	if len(groups) > 1 {
		n1 := castx.ToInt64(groups[1]) * 24

		if isNegative {
			d1 -= time.Duration(n1) * time.Hour
		} else {
			d1 += time.Duration(n1) * time.Hour
		}
	}

	re = regexp.MustCompile("([1-9][0-9]*)h")
	groups = re.FindStringSubmatch(str)

	if len(groups) > 1 {
		n1 := castx.ToInt64(groups[1])

		if isNegative {
			d1 -= time.Duration(n1) * time.Hour
		} else {
			d1 += time.Duration(n1) * time.Hour
		}
	}

	re = regexp.MustCompile("([1-9][0-9]*)m")
	groups = re.FindStringSubmatch(str)

	if len(groups) > 1 {
		n1 := castx.ToInt64(groups[1])

		if isNegative {
			d1 -= time.Duration(n1) * time.Minute
		} else {
			d1 += time.Duration(n1) * time.Minute
		}
	}

	re = regexp.MustCompile("([1-9][0-9]*)s")
	groups = re.FindStringSubmatch(str)

	if len(groups) > 1 {
		n1 := castx.ToInt64(groups[1])

		if isNegative {
			d1 -= time.Duration(n1) * time.Second
		} else {
			d1 += time.Duration(n1) * time.Second
		}
	}

	re = regexp.MustCompile("([1-9][0-9]*)ms")
	groups = re.FindStringSubmatch(str)

	if len(groups) > 1 {
		n1 := castx.ToInt64(groups[1])

		if isNegative {
			d1 -= time.Duration(n1) * time.Millisecond
		} else {
			d1 += time.Duration(n1) * time.Millisecond
		}
	}

	re = regexp.MustCompile("([1-9][0-9]*)us|µs")
	groups = re.FindStringSubmatch(str)

	if len(groups) > 1 {
		n1 := castx.ToInt64(groups[1])

		if isNegative {
			d1 -= time.Duration(n1) * time.Microsecond
		} else {
			d1 += time.Duration(n1) * time.Microsecond
		}
	}

	re = regexp.MustCompile("([1-9][0-9]*)ns")
	groups = re.FindStringSubmatch(str)

	if len(groups) > 1 {
		n1 := castx.ToInt64(groups[1])

		if isNegative {
			d1 -= time.Duration(n1) * time.Nanosecond
		} else {
			d1 += time.Duration(n1) * time.Nanosecond
		}
	}

	return d1
}

func TrimStart(str string, chars ...rune) string {
	if str == "" || len(chars) < 1 {
		return str
	}

	return strings.TrimLeftFunc(str, func(r rune) bool {
		for _, ch := range chars {
			if ch == r {
				return true
			}
		}

		return false
	})
}

func TrimEnd(str string, chars ...rune) string {
	if str == "" || len(chars) < 1 {
		return str
	}

	return strings.TrimRightFunc(str, func(r rune) bool {
		for _, ch := range chars {
			if ch == r {
				return true
			}
		}

		return false
	})
}

func Trim(str string, chars ...rune) string {
	if str == "" || len(chars) < 1 {
		return str
	}

	return TrimEnd(TrimStart(str, chars...), chars...)
}

func PadStart(str string, length int, padChars ...string) string {
	if str == "" {
		return ""
	}

	n1 := len(str)

	if n1 >= length {
		return str
	}

	padChar := " "

	if len(padChars) > 0 {
		padChar = padChars[0]
	}

	return strings.Repeat(padChar, length - n1) + str
}

func PadEnd(str string, length int, padChars ...string) string {
	if str == "" {
		return ""
	}

	n1 := len(str)

	if n1 >= length {
		return str
	}

	padChar := " "

	if len(padChars) > 0 {
		padChar = padChars[0]
	}

	return str + strings.Repeat(padChar, length - n1)
}

func GetRandomString(length int, randomStringType ...int) string {
	if length < 1 {
		return ""
	}

	_randomStringType := RandomStringType.Default
	supportedTypes := []int{RandomStringType.Default, RandomStringType.Alpha, RandomStringType.Alnum}

	if len(randomStringType) > 0 && slicex.InIntSlice(randomStringType[0], supportedTypes) {
		_randomStringType = randomStringType[0]
	}

	chars := make([]rune, 0)

	switch _randomStringType {
	case RandomStringType.Alpha:
		for ch := 'A'; ch <= 'Z'; ch++ {
			chars = append(chars, ch)
		}

		for ch := 'a'; ch <= 'z'; ch++ {
			chars = append(chars, ch)
		}
	case RandomStringType.Alnum:
		for ch := '0'; ch <= '9'; ch++ {
			chars = append(chars, ch)
		}
	default:
		for ch := 'A'; ch <= 'Z'; ch++ {
			chars = append(chars, ch)
		}

		for ch := '0'; ch <= '9'; ch++ {
			chars = append(chars, ch)
		}

		for ch := 'a'; ch <= 'z'; ch++ {
			chars = append(chars, ch)
		}
	}

	max := len(chars) - 1
	sb := strings.Builder{}
	_rand := mrand.New(&customRandSource{})

	for i := 1; i <= length; i++ {
		idx := _rand.Intn(max)
		sb.WriteRune(chars[idx])
	}

	return sb.String()
}

func ThousandSep(str string) string {
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

func ReadFromFile(fpath string) string {
	buf, _ := ioutil.ReadFile(fpath)

	if len(buf) < 1 {
		return ""
	}

	return string(buf)
}

func IsPasswordTooSimple(str string) bool {
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
	}, str)

	return n1 + n2 + n3 + n4 < 3
}

func IsIdcard(str string) bool {
	pattern := `^[1-9]\d{5}(18|19|20)\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$`

	if !RegexMatch(str, pattern) {
		return false
	}

	a1 := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}

	a2 := []string{
		"1",
		"0",
		"x",
		"9",
		"8",
		"7",
		"6",
		"5",
		"4",
		"3",
		"2",
	}

	sum := 0
	n1 := len(str) - 1

	for i := 0; i < n1; i++ {
		sum += castx.ToInt(str[i:i+1]) * a1[i]
	}

	return strings.ToLower(a2[sum % 11]) == strings.ToLower(str[n1:])
}

func GetAgeByBirthday(str string) int {
	t1, err := time.Parse(DatetimeFormat.DateOnly, str)

	if err != nil {
		return -1
	}

	year1 := t1.Year()
	month1 := int(t1.Month())
	day1 := t1.Day()
	t2 := time.Now()
	year2 := t2.Year()
	month2 := int(t2.Month())
	day2 := t2.Day()
	age := year2 - year1

	if age < 0 {
		return -1
	}

	if month2 > month1 || (month2 == month1 && day2 >= day1) {
		age += 1
	}

	return age
}

func GetBirthdayByIdcardNo(str string) string {
	if len(str) < 14 {
		return ""
	}

	year := str[6:10]
	month := str[10:12]
	day := str[12:14]
	birthday := fmt.Sprintf(format1, year, month, day)

	if !IsDate(birthday) {
		return ""
	}

	return birthday
}

func HandleTextareaInput(str string) string {
	if str == "" {
		return ""
	}

	str = strings.ReplaceAll(str, "\r", "")
	str = ReplaceWithRegexp(str, "[\\x20\\t]*\n", "")
	str = strings.ReplaceAll(str, "\n", "")
	str = strings.TrimSpace(str)
	return str
}

func ToGb18030(str string) string {
	if str == "" {
		return ""
	}

	buf, err := simplifiedchinese.GB18030.NewEncoder().Bytes([]byte(str))

	if err != nil {
		return ""
	}

	return string(buf)
}

func UnescapeUnicode(arg0 interface{}) string {
	var buf []byte

	if _buf, ok := arg0.([]byte); ok {
		buf = _buf
	} else if s1, ok := arg0.(string); ok && s1 != "" {
		buf = []byte(s1)
	}

	if len(buf) < 1 {
		return ""
	}

	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(buf)), `\\u`, `\u`, -1))

	if err != nil {
		return ""
	}

	return str
}

func StripTags(str string) string {
	if str == "" || IsNumeric(str) || IsAlnum(str) {
		return str
	}

	return strip.StripTags(str)
}
