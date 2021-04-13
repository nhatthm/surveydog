package surveydog

import (
	"sync"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/cucumber/godog"
	"github.com/nhatthm/surveymock"
)

// Manager is a wrapper around *surveymock.Survey to make it run with cucumber/godog.
type Manager struct {
	surveys map[string]*Survey
	current string

	mu sync.Mutex

	options []surveymock.MockOption
}

// RegisterContext register the survey to a *godog.ScenarioContext.
func (m *Manager) RegisterContext(t surveymock.TestingT, ctx *godog.ScenarioContext) {
	// Manage state.
	ctx.BeforeScenario(func(sc *godog.Scenario) {
		m.beforeScenario(t, sc)
	})
	ctx.AfterScenario(m.afterScenario)

	// Confirm prompt
	ctx.Step(`gets? a(?:nother)? confirm prompt "([^"]*)".* answers? yes`, m.expectConfirmYes)
	ctx.Step(`gets? a(?:nother)? confirm prompt "([^"]*)".* answers? no`, m.expectConfirmNo)
	ctx.Step(`gets? a(?:nother)? confirm prompt "([^"]*)".* answers? "([^"]*)"`, m.expectConfirmAnswer)
	ctx.Step(`gets? a(?:nother)? confirm prompt "([^"]*)".* interrupts?`, m.expectConfirmInterrupt)
	ctx.Step(`gets? a(?:nother)? confirm prompt "([^"]*)".* asks? for help and sees? "([^"]*)"`, m.expectConfirmHelp)

	// Password prompt.
	ctx.Step(`gets? a(?:nother)? password prompt "([^"]*)".* answers? "([^"]*)"`, m.expectPasswordAnswer)
	ctx.Step(`gets? a(?:nother)? password prompt "([^"]*)".* interrupts?`, m.expectPasswordInterrupt)
	ctx.Step(`gets? a(?:nother)? password prompt "([^"]*)".* asks? for help and sees? "([^"]*)"`, m.expectPasswordHelp)
}

// Stdio returns terminal.Stdio of the current survey.
func (m *Manager) Stdio() terminal.Stdio {
	return m.survey().Stdio()
}

func (m *Manager) beforeScenario(t surveymock.TestingT, sc *godog.Scenario) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.current = sc.Id
	m.surveys[m.current] = NewSurvey(t, m.options...).Start()
}

func (m *Manager) afterScenario(sc *godog.Scenario, _ error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if survey := m.surveys[sc.Id]; survey != nil {
		survey.Close()
		delete(m.surveys, sc.Id)
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

func (m *Manager) expectConfirmAnswer(message string, answer string) error {
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

func (m *Manager) expectPasswordAnswer(message string, answer string) error {
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

// New initiates a new *surveydog.Manager.
func New(options ...surveymock.MockOption) *Manager {
	return &Manager{
		surveys: make(map[string]*Survey),
		options: options,
	}
}
