package errorx

import (
	"github.com/go-errors/errors"
	"strings"
)

func StacktraceAsArray(arg0 interface{}) []string {
	traceLines := make([]string, 0)

	if arg0 == nil {
		return traceLines
	}

	var stacktrace string

	switch t := arg0.(type) {
	case *errors.Error:
		stacktrace = t.ErrorStack()
	case error:
		stacktrace = errors.New(t).ErrorStack()
	}

	if stacktrace == "" {
		return traceLines
	}

	stacktrace = strings.ReplaceAll(stacktrace, "\r", "")
	lines := strings.Split(stacktrace, "\n")

	if strings.Contains(stacktrace, "src/runtime/panic.go") {
		n1 := -1

		for i := 0; i < len(lines); i++ {
			if i == 0 {
				traceLines = append(traceLines, lines[i])
				continue
			}

			if strings.Contains(lines[i], "src/runtime/panic.go") {
				n1 = i
				continue
			}

			if strings.Contains(lines[i], "src/runtime/proc.go") ||
				strings.Contains(lines[i], "src/runtime/asm_amd64") {
				break
			}

			if n1 < 0 || i <= n1 + 1 {
				continue
			}

			traceLines = append(traceLines, lines[i])
		}
	} else {
		for i := 0; i < len(lines); i++ {
			if strings.Contains(lines[i], "src/runtime/proc.go") ||
				strings.Contains(lines[i], "src/runtime/asm_amd64") {
				break
			}

			traceLines = append(traceLines, lines[i])
		}
	}

	return traceLines
}

func Stacktrace(arg0 interface{}, separator ...string) string {
	sep := "\n"

	if len(separator) > 0 && separator[0] != "" {
		sep = separator[0]
	}

	lines := StacktraceAsArray(arg0)

	if len(lines) < 1 {
		return ""
	}

	return strings.Join(lines, sep)
}
