package jsonx

type option struct {
	handleTimeField   bool
	stripZeroTimePart bool
	replacements      [][2]string
}

func NewToJsonOption() *option {
	return &option{}
}

func (opt *option) HandleTimeField() *option {
	opt.handleTimeField = true
	return opt
}

func (opt *option) StripZeroTimePart() *option {
	opt.stripZeroTimePart = true
	return opt
}

func (opt *option) WithReplacements(replacements [][2]string) *option {
	if len(replacements) > 0 {
		opt.replacements = append(opt.replacements, replacements...)
	}

	return opt
}
