package surveydog_test

import (
	"testing"
	"time"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/nhatthm/surveymock"
	"github.com/nhatthm/surveymock/options"
	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/surveydog"
)

func TestSurvey_ExpectationsWereMet(t *testing.T) {
	t.Parallel()

	s := surveydog.NewSurvey(t, func(s *surveymock.Survey) {
		s.ExpectPassword("Enter password:").Answer("password")
	}).Start("test")

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

	s.Close()
}

func TestSurvey_ExpectationsWereNotMet(t *testing.T) {
	t.Parallel()

	testingT := T()
	s := surveydog.NewSurvey(testingT, func(s *surveymock.Survey) {
		s.ExpectPassword("Enter password:")
	}).Start("test")

	time.Sleep(50 * time.Millisecond)

	s.Close()

	expectedErr := "there are remaining expectations that were not met:\n\nType   : Password\nMessage: \"Enter password:\"\nAnswer : <no answer>\n"

	assert.EqualError(t, s.ExpectationsWereMet(), expectedErr)
}
