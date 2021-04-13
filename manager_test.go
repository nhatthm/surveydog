package surveydog

import (
	"fmt"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/nhatthm/surveymock"
	"github.com/stretchr/testify/assert"
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

func TestManager_ExpectationsWereNotMet(t *testing.T) {
	t.Parallel()

	testingT := T()
	s := New()
	sc := &godog.Scenario{Id: "42", Name: "ExpectationsWereNotMet"}

	s.beforeScenario(testingT, sc)

	assert.Nil(t, s.expectPasswordAnswer("Enter password:", "password"))

	<-time.After(50 * time.Millisecond)

	s.afterScenario(testingT, sc)

	expectedError := `in scenario "ExpectationsWereNotMet", there are remaining expectations that were not met:\
[\t\s]*Type   : Password\
[\t\s]*Message: "Enter password:"\
[\t\s]*Answer : "password"\
`

	assert.Regexp(t, expectedError, testingT.ErrorString())
}
