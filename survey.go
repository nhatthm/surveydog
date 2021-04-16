package surveydog

import (
	"errors"
	"sync"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/nhatthm/surveyexpect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Survey is a wrapper around *surveyexpect.Survey to make it run with cucumber/godog.
type Survey struct {
	*surveyexpect.Survey
	console surveyexpect.Console
	output  *surveyexpect.Buffer
	state   *vt10x.State

	test surveyexpect.TestingT
	mu   sync.Mutex

	doneChan chan struct{}
}

func (s *Survey) getDoneChan() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.getDoneChanLocked()
}

func (s *Survey) getDoneChanLocked() chan struct{} {
	if s.doneChan == nil {
		s.doneChan = make(chan struct{})
	}

	return s.doneChan
}

func (s *Survey) closeDoneChan() {
	s.mu.Lock()
	defer s.mu.Unlock()

	ch := s.getDoneChanLocked()

	select {
	case <-ch:
		// Already closed. Don't close again.

	default:
		// Safe to close here. We're the only closer, guarded
		// by s.mu.
		close(ch)
	}
}

// Stdio returns terminal.Stdio from surveyexpect.Console.
func (s *Survey) Stdio() terminal.Stdio {
	s.mu.Lock()
	defer s.mu.Unlock()

	return terminal.Stdio{
		In:  s.console.Tty(),
		Out: s.console.Tty(),
		Err: s.console.Tty(),
	}
}

// Expect runs an expectation against a given console.
func (s *Survey) Expect(c surveyexpect.Console) error {
	for {
		select {
		case <-s.getDoneChan():
			return nil

		default:
			err := s.Survey.Expect(c)
			if err != nil && !errors.Is(err, surveyexpect.ErrNoExpectation) {
				return err
			}
		}
	}
}

// Start starts a new survey.
func (s *Survey) Start(scenario string) *Survey {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.test.Logf("Scenario: %s\n", scenario)

	s.output = new(surveyexpect.Buffer)

	console, state, err := vt10x.NewVT10XConsole(expect.WithStdout(s.output))
	require.NoError(s.test, err)

	s.console = console
	s.state = state

	go func() {
		assert.NoError(s.test, s.Expect(s.console))
	}()

	return s
}

// Close notifies other parties and close the survey.
func (s *Survey) Close() {
	s.closeDoneChan()

	s.mu.Lock()
	defer s.mu.Unlock()

	s.test.Logf("Raw output: %q\n", s.output.String())

	// Dump the terminal's screen.
	s.test.Logf("State: \n%s\n", expect.StripTrailingEmptyLines(s.state.String()))
	s.test.Log()

	s.console = nil
	s.state = nil
	s.output = nil
}

// NewSurvey creates a new survey.
func NewSurvey(t surveyexpect.TestingT, options ...surveyexpect.ExpectOption) *Survey {
	return &Survey{
		Survey: surveyexpect.New(t, options...),
		test:   t,
	}
}
