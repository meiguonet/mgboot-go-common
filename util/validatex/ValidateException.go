package validatex

type ValidateException struct {
	failfast  bool
	errorTips string
	errors    map[string]string
}

func NewValidateExceptionWithErrorTips(errorTips string) ValidateException {
	return ValidateException{
		failfast:  true,
		errorTips: errorTips,
	}
}

func NewValidateExceptionWithErrors(errors map[string]string) ValidateException {
	return ValidateException{errors: errors}
}

func (ex ValidateException) Error() string {
	if ex.failfast {
		return ex.errorTips
	}

	return "data validate exception"
}

func (ex ValidateException) IsFailfast() bool {
	return ex.failfast
}

func (ex ValidateException) GetErrors() map[string]string {
	if len(ex.errors) < 1 {
		return map[string]string{}
	}

	return ex.errors
}
