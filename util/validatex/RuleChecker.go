package validatex

type RuleChecker interface {
	GetRuleName() string
	Check(value string, checkValue ...string) bool
}
