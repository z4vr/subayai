package errorarray

import (
	"fmt"
	"strings"
)

type FormatFunc func(errors []error) string

var (
	defaultFormatFunc = func(errors []error) string {
		ln := len(errors)
		lines := make([]string, ln+1)
		lines[0] = fmt.Sprintf("ErrorArray stack [%d inner errors]:", ln)
		for i := 1; i < ln+1; i++ {
			lines[i] = fmt.Sprintf("  [%02d] %s", i-1, errors[i-1].Error())
		}
		return strings.Join(lines, "\n")
	}
)

type ErrorArray struct {
	errors     []error
	formatFunc FormatFunc
}

func New(formatFunc ...FormatFunc) (e *ErrorArray) {
	e = new(ErrorArray)

	if formatFunc != nil && len(formatFunc) > 0 && formatFunc[0] != nil {
		e.formatFunc = formatFunc[0]
	} else {
		e.formatFunc = defaultFormatFunc
	}

	return
}

func (e *ErrorArray) Error() string {
	if e.Len() == 0 {
		return ""
	}
	return e.formatFunc(e.errors)
}

// Errors returns the internal list of errors.
func (e *ErrorArray) Errors() []error {
	return e.errors
}

// ForEach iterates over the ErrorArray and
// executes the given function for each error
func (e *ErrorArray) ForEach(f func(err error, i int)) {
	for i, err := range e.errors {
		f(err, i)
	}
}

// Append adds the given error to the ErrorArray
// and returns the ErrorArray instance.
func (e *ErrorArray) Append(err ...error) {
	for _, err := range err {
		if err != nil {
			e.errors = append(e.errors, err)
		}
	}
}

// Len returns the number of errors in the ErrorArray
func (e *ErrorArray) Len() int {
	return len(e.errors)
}

// Nillify returns the ErrorArray instance
// if it has errors, otherwise it returns nil.
func (e *ErrorArray) Nillify() error {
	if e.Len() > 0 {
		return e
	}
	return nil
}
