package surveydog

import (
	"fmt"
	"sync"
	"time"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/Netflix/go-expect"
	"github.com/cucumber/godog"
	"github.com/nhatthm/consoledog"
	"github.com/nhatthm/surveyexpect"
	"github.com/stretchr/testify/assert"
)

// Starter is a callback when survey starts.
type Starter func(sc *godog.Scenario, stdio terminal.Stdio)

// Manager is a wrapper around *surveyexpect.Survey to make it run with cucumber/godog.
type Manager struct {
	console *consoledog.Manager
	surveys map[string]*Survey
	current string

	starters []Starter

	test surveyexpect.TestingT
	mu   sync.Mutex

	options []surveyexpect.ExpectOption
}

func (m *Manager) registerConsole(ctx *godog.ScenarioContext) {
	if m.console != nil {
		return
	}

	console := consoledog.New(m.test)
	m.attach(console)

	// Manage state.
	ctx.BeforeScenario(func(sc *godog.Scenario) {
		console.NewConsole(sc)
	})

	ctx.AfterScenario(func(sc *godog.Scenario, _ error) {
		console.CloseConsole(sc)
	})
}

// RegisterContext register the survey to a *godog.ScenarioContext.
func (m *Manager) RegisterContext(ctx *godog.ScenarioContext) {
	m.registerConsole(ctx)

	// Confirm prompt
	ctx.Step(`(?:(?:get)|(?:see))s? a(?:nother)? confirm prompt "([^"]*)".* answers? yes`, m.expectConfirmYes)
	ctx.Step(`(?:(?:get)|(?:see))s? a(?:nother)? confirm prompt "([^"]*)".* answers? no`, m.expectConfirmNo)
	ctx.Step(`(?:(?:get)|(?:see))s? a(?:nother)? confirm prompt "([^"]*)".* answers? "([^"]*)"`, m.expectConfirmAnswer)
	ctx.Step(`(?:(?:get)|(?:see))s? a(?:nother)? confirm prompt "([^"]*)".* interrupts?`, m.expectConfirmInterrupt)
	ctx.Step(`(?:(?:get)|(?:see))s? a(?:nother)? confirm prompt "([^"]*)".* asks? for help and sees? "([^"]*)"`, m.expectConfirmHelp)

	// Password prompt.
	ctx.Step(`(?:(?:get)|(?:see))s? a(?:nother)? password prompt "([^"]*)".* answers? "([^"]*)"`, m.expectPasswordAnswer)
	ctx.Step(`(?:(?:get)|(?:see))s? a(?:nother)? password prompt "([^"]*)".* interrupts?`, m.expectPasswordInterrupt)
	ctx.Step(`(?:(?:get)|(?:see))s? a(?:nother)? password prompt "([^"]*)".* asks? for help and sees? "([^"]*)"`, m.expectPasswordHelp)
}

func (m *Manager) start(sc *godog.Scenario, console *expect.Console) {
	m.mu.Lock()
	defer m.mu.Unlock()

	s := NewSurvey(m.test, m.options...).Start(console)

	m.current = sc.Id
	m.surveys[m.current] = s

	for _, start := range m.starters {
		start(sc, terminal.Stdio{
			In:  console.Tty(),
			Out: console.Tty(),
			Err: console.Tty(),
		})
	}
}

func (m *Manager) close(sc *godog.Scenario) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if s := m.surveys[sc.Id]; s != nil {
		s.Close()
		delete(m.surveys, sc.Id)

		assert.NoError(m.test, m.expectationsWereMet(sc.Name, s))
	}

	m.current = ""
}

func (m *Manager) survey() *Survey {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.surveys[m.current]
}

func (m *Manager) expectConfirmYes(message string) error {
	m.survey().ExpectConfirm(message).Yes()

	return nil
}

func (m *Manager) expectConfirmNo(message string) error {
	m.survey().ExpectConfirm(message).No()

	return nil
}

func (m *Manager) expectConfirmAnswer(message, answer string) error {
	m.survey().ExpectConfirm(message).Answer(answer)

	return nil
}

func (m *Manager) expectConfirmInterrupt(message string) error {
	m.survey().ExpectConfirm(message).Interrupt()

	return nil
}

func (m *Manager) expectConfirmHelp(message, help string) error {
	m.survey().ExpectConfirm(message).ShowHelp(help)

	return nil
}

func (m *Manager) expectPasswordAnswer(message, answer string) error {
	m.survey().ExpectPassword(message).Answer(answer)

	return nil
}

func (m *Manager) expectPasswordInterrupt(message string) error {
	m.survey().ExpectPassword(message).Interrupt()

	return nil
}

func (m *Manager) expectPasswordHelp(message, help string) error {
	m.survey().ExpectPassword(message).ShowHelp(help)

	return nil
}

// expectationsWereMet checks whether all queued expectations were met in order.
// If any of them was not met - an error is returned.
func (m *Manager) expectationsWereMet(scenario string, s *Survey) error {
	<-time.After(surveyexpect.ReactionTime)

	err := s.ExpectationsWereMet()
	if err == nil {
		return nil
	}

	return fmt.Errorf("in scenario %q, %w", scenario, err)
}

func (m *Manager) attach(console *consoledog.Manager) *consoledog.Manager {
	return console.
		WithStarter(m.start).
		WithCloser(m.close)
}

// WithConsole sets console manager.
func (m *Manager) WithConsole(console *consoledog.Manager) *Manager {
	m.console = m.attach(console)

	return m
}

// WithStarter adds a mew Starter to Manager.
func (m *Manager) WithStarter(s Starter) *Manager {
	m.starters = append(m.starters, s)

	return m
}

// New initiates a new *surveydog.Manager.
func New(t surveyexpect.TestingT, options ...surveyexpect.ExpectOption) *Manager {
	return &Manager{
		surveys: make(map[string]*Survey),
		options: options,
		test:    t,
	}
}
