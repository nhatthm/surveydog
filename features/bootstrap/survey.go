package bootstrap

import (
	"errors"
	"fmt"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/cucumber/godog"
	"github.com/nhatthm/surveyexpect/options"
	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/surveydog"
)

// Prompt produces prompts and validate the answers.
type Prompt struct {
	m *surveydog.Manager
}

// RegisterContext register Prompt to a *godog.ScenarioContext.
func (p *Prompt) RegisterContext(ctx *godog.ScenarioContext) {
	ctx.Step(`ask for confirm "([^"]*)" with help "([^"]*)", receive yes`, p.askConfirmWithHelpYes)
	ctx.Step(`ask for confirm "([^"]*)" with help "([^"]*)", receive no`, p.askConfirmWithHelpNo)
	ctx.Step(`ask for confirm "([^"]*)", receive yes`, p.askConfirmWithoutHelpYes)
	ctx.Step(`ask for confirm "([^"]*)", receive no`, p.askConfirmWithoutHelpNo)
	ctx.Step(`ask for confirm "([^"]*)", get interrupted`, p.askConfirmInterrupted)

	ctx.Step(`ask for password "([^"]*)" with help "([^"]*)", receive "([^"]*)"`, p.askPasswordWithHelp)
	ctx.Step(`ask for password "([^"]*)", receive "([^"]*)"`, p.askPasswordWithoutHelp)
	ctx.Step(`ask for password "([^"]*)", get interrupted`, p.askPasswordInterrupted)
}

func (p *Prompt) ask(prompt survey.Prompt, response interface{}) (err error) {
	doneCh := make(chan struct{})

	go func() {
		defer close(doneCh)

		err = survey.AskOne(prompt, response, options.WithStdio(p.m.Stdio()))
	}()

	select {
	case <-time.After(time.Second):
		return errors.New("ask timed out") // nolint: goerr113

	case <-doneCh:
		return nil
	}
}

func (p *Prompt) askConfirm(message, help string) (bool, error) {
	prompt := &survey.Confirm{
		Message: message,
		Help:    help,
	}

	var answer bool
	err := p.ask(prompt, &answer)

	return answer, err
}

func (p *Prompt) askConfirmWithHelp(message, help string, expectedAnswer bool) error {
	answer, err := p.askConfirm(message, help)
	if err != nil {
		return err
	}

	if !assert.ObjectsAreEqual(expectedAnswer, answer) {
		return fmt.Errorf("expected answer: %t, got %t", expectedAnswer, answer) // nolint: goerr113
	}

	return nil
}

func (p *Prompt) askConfirmWithHelpYes(message, help string) error {
	return p.askConfirmWithHelp(message, help, true)
}

func (p *Prompt) askConfirmWithHelpNo(message, help string) error {
	return p.askConfirmWithHelp(message, help, false)
}

func (p *Prompt) askConfirmWithoutHelpYes(message string) error {
	return p.askConfirmWithHelp(message, "", true)
}

func (p *Prompt) askConfirmWithoutHelpNo(message string) error {
	return p.askConfirmWithHelp(message, "", false)
}

func (p *Prompt) askConfirmInterrupted(message string) error {
	answer, err := p.askConfirm(message, "")

	if answer {
		return fmt.Errorf("unexpected answer: %t", answer) // nolint: goerr113
	}

	if !errors.Is(err, terminal.InterruptErr) {
		return err
	}

	return nil
}

func (p *Prompt) askPassword(message, help string) (string, error) {
	prompt := &survey.Password{
		Message: message,
		Help:    help,
	}

	var answer string
	err := p.ask(prompt, &answer)

	return answer, err
}

func (p *Prompt) askPasswordWithHelp(message, help, expectedAnswer string) error {
	answer, err := p.askPassword(message, help)
	if err != nil {
		return err
	}

	if !assert.ObjectsAreEqual(expectedAnswer, answer) {
		return fmt.Errorf("expected answer: %s, got %s", expectedAnswer, answer) // nolint: goerr113
	}

	return nil
}

func (p *Prompt) askPasswordWithoutHelp(message, expectedAnswer string) error {
	return p.askPasswordWithHelp(message, "", expectedAnswer)
}

func (p *Prompt) askPasswordInterrupted(message string) error {
	answer, err := p.askPassword(message, "")

	if answer != "" {
		return fmt.Errorf("unexpected answer: %s", answer) // nolint: goerr113
	}

	if !errors.Is(err, terminal.InterruptErr) {
		return err
	}

	return nil
}

// NewPrompt initiates a new *Prompt.
func NewPrompt(m *surveydog.Manager) *Prompt {
	return &Prompt{
		m: m,
	}
}
