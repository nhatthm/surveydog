package surveydog_test

import (
	"fmt"

	"github.com/nhatthm/surveymock"
)

type TestingT struct {
	error *surveymock.Buffer
	log   *surveymock.Buffer

	clean func()
}

func (t *TestingT) Errorf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(t.error, format, args...)
}

func (t *TestingT) Log(args ...interface{}) {
	_, _ = fmt.Fprintln(t.log, args...)
}

func (t *TestingT) Logf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(t.log, format, args...)
}

func (t *TestingT) FailNow() {
	panic("failed")
}

func (t *TestingT) Cleanup(clean func()) {
	t.clean = clean
}

func (t *TestingT) ErrorString() string {
	return t.error.String()
}

func (t *TestingT) LogString() string {
	return t.log.String()
}

func T() *TestingT {
	return &TestingT{
		error: new(surveymock.Buffer),
		log:   new(surveymock.Buffer),
		clean: func() {},
	}
}
