package surveydog

import (
	"fmt"
	"testing"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/cucumber/godog"
	"github.com/nhatthm/surveymock"
	"github.com/nhatthm/surveymock/options"
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

func TestManager_expectationsWereMet(t *testing.T) {
	t.Parallel()

	testingT := T()
	s := New()
	sc := &godog.Scenario{Id: "42", Name: "ExpectationsWereMet"}

	s.beforeScenario(testingT, sc)

	assert.Nil(t, s.expectPasswordAnswer("Enter password:", "password"))

	doneCh := make(chan struct{}, 1)

	go func() {
		defer close(doneCh)

		var answer string
		err := survey.AskOne(&survey.Password{Message: "Enter password:"}, &answer, options.WithStdio(s.Stdio()))

		assert.Equal(t, "password", answer)
		assert.NoError(t, err)
	}()

	select {
	case <-time.After(100 * time.Millisecond):
		t.Error("ask timeout")

	case <-doneCh:
	}

	s.afterScenario(testingT, sc)

	assert.Empty(t, testingT.ErrorString())
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

	expectedError := "in scenario \"ExpectationsWereNotMet\", there are remaining expectations that were not met:\n\t            \t\n\t            \tType   : Password\n\t            \tMessage: \"Enter password:\"\n\t            \tAnswer : \"password\"\n"

	assert.Contains(t, testingT.ErrorString(), expectedError)
}
