package surveydog

import (
	"fmt"
	"sync"
	"time"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/cucumber/godog"
	"github.com/nhatthm/surveyexpect"
	"github.com/stretchr/testify/assert"
)

// Manager is a wrapper around *surveyexpect.Survey to make it run with cucumber/godog.
type Manager struct {
	surveys map[string]*Survey
	current string

	mu sync.Mutex

	options []surveyexpect.ExpectOption
}

// RegisterContext register the survey to a *godog.ScenarioContext.
func (m *Manager) RegisterContext(t surveyexpect.TestingT, ctx *godog.ScenarioContext, listeners ...func(sc *godog.Scenario, stdio terminal.Stdio)) {
	// Manage state.
	ctx.BeforeScenario(func(sc *godog.Scenario) {
		m.beforeScenario(t, sc, listeners...)
	})

	ctx.AfterScenario(func(sc *godog.Scenario, _ error) {
		m.afterScenario(t, sc)
	})

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

// Stdio returns terminal.Stdio of the current survey.
func (m *Manager) Stdio() terminal.Stdio {
	return m.survey().Stdio()
}

func (m *Manager) beforeScenario(t surveyexpect.TestingT, sc *godog.Scenario, listeners ...func(sc *godog.Scenario, stdio terminal.Stdio)) {
	m.mu.Lock()
	defer m.mu.Unlock()

	s := NewSurvey(t, m.options...).Start(sc.Name)

	m.current = sc.Id
	m.surveys[m.current] = s

	for _, l := range listeners {
		l(sc, s.Stdio())
	}
}

func (m *Manager) afterScenario(t surveyexpect.TestingT, sc *godog.Scenario) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if s := m.surveys[sc.Id]; s != nil {
		s.Close()
		delete(m.surveys, sc.Id)

		assert.NoError(t, m.expectationsWereMet(sc.Name, s))
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

// New initiates a new *surveydog.Manager.
func New(options ...surveyexpect.ExpectOption) *Manager {
	return &Manager{
		surveys: make(map[string]*Survey),
		options: options,
	}
}
