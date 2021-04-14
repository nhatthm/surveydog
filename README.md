# Cucumber Survey steps for Go

[![Build Status](https://github.com/nhatthm/surveydog/actions/workflows/test.yaml/badge.svg)](https://github.com/nhatthm/surveydog/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/nhatthm/surveydog/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/nhatthm/surveydog)
[![Go Report Card](https://goreportcard.com/badge/github.com/nhatthm/surveydog)](https://goreportcard.com/report/github.com/nhatthm/surveydog)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/github.com/nhatthm/surveydog)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

Tests with `AlecAivazis/survey` and `cucumber/godog`

## Prerequisites

- `Go >= 1.14`

## Install

```bash
go get github.com/nhatthm/surveydog
```

## Usage

### Supported Types

For now, it only supports [`Confirm`](https://github.com/AlecAivazis/survey#confirm) and [`Password`](https://github.com/AlecAivazis/survey#password)

## Setup

Step 1: Register to `godog`

Initialize a `surveydog.Manager` with `surveydog.New()` then add it into the `ScenarioInitializer`

Step 2: Pass `stdio` to the prompts

Same as [`surveyexpect`](https://github.com/nhatthm/surveyexpect#expect), you have to define a way to inject `terminal.Stdio` into the prompts in your code. For
every scenario, the manager will start a new terminal emulator. Without the injection, there is no way to capture and response to the prompts.

For example:

```go
package mypackage

import (
    "math/rand"
    "testing"

    survey "github.com/AlecAivazis/survey/v2"
    "github.com/AlecAivazis/survey/v2/terminal"
    "github.com/cucumber/godog"
    "github.com/nhatthm/surveydog"
    "github.com/nhatthm/surveyexpect/options"
)

type Wizard struct {
    stdio terminal.Stdio
}

func (w *Wizard) withStdio(stdio terminal.Stdio) {
    w.stdio = stdio
}

func (w *Wizard) ask() (bool, error) {
    var response bool

    p := &survey.Confirm{Message: "Confirm?"}
    err := survey.AskOne(p, &response, options.WithStdio(w.stdio))

    return response, err
}

func TestIntegration(t *testing.T) {
    m := surveydog.New()
    wizard := &Wizard{}

    suite := godog.TestSuite{
        Name: "Integration",
        ScenarioInitializer: func(ctx *godog.ScenarioContext) {
            // Register `surveydog.Manager` with listeners.
            m.RegisterContext(t, ctx, func(_ *godog.Scenario, stdio terminal.Stdio) {
                // Update stdio of the wizard.
                wizard.withStdio(stdio)
            })
        },
        Options: &godog.Options{
            Strict:    true,
            Output:    out,
            Randomize: rand.Int63(),
        },
    }

    // Run the suite that triggers wizard.ask()
}
```

See more: [#Examples](#Examples)

## Steps

### Confirm

#### Yes

Expect to see a Confirm prompt and answer `yes`.

Pattern: `(?:(?:get)|(?:see))s? a(?:nother)? confirm prompt "([^"]*)".* answers? yes`

Example:

```gherkin
    Scenario: Receive a yes
        Given I see a confirm prompt "Confirm? (y/N)", I answer yes

        Then ask for confirm "Confirm?", receive yes
```

#### No

Expect to see a Confirm prompt and answer `no`.

Pattern: `(?:(?:get)|(?:see))s? a(?:nother)? confirm prompt "([^"]*)".* answers? no`

Example:

```gherkin
    Scenario: Receive a no
        Given I see a confirm prompt "Confirm? (y/N)", I answer no

        Then ask for confirm "Confirm?", receive no
```

#### Invalid answer

Expect to see a Confirm prompt and answer an invalid response (not a `yes` or `no`).

Pattern: `(?:(?:get)|(?:see))s? a(?:nother)? confirm prompt "([^"]*)".* answers? "([^"]*)"`

Example:

```gherkin
    Scenario: Invalid answer
        Given I see a confirm prompt "Confirm? (y/N)", I answer "nahhh"
        # Because the answer is invalid, survey will prompt again.
        And then I see another confirm prompt "Confirm? (y/N)", I answer no

        Then ask for confirm "Confirm?", receive no
```

#### Interrupt

Expect to see a Confirm prompt and interrupt (^C).

Pattern: `(?:(?:get)|(?:see))s? a(?:nother)? confirm prompt "([^"]*)".* interrupts?`

Example:

```gherkin
    Scenario: Interrupted
        Given I see a confirm prompt "Confirm? (y/N)", I interrupt

        Then ask for confirm "Confirm?", get interrupted
```

#### With Help

Expect to see a Confirm prompt, ask for help and then expect to see a Help message.

Pattern: `(?:(?:get)|(?:see))s? a(?:nother)? confirm prompt "([^"]*)".* asks? for help and sees? "([^"]*)"`

Example:

```gherkin
    Scenario: With help and receive a yes
        Given I see a confirm prompt "Confirm? [? for help] (y/N)", I ask for help and see "This action cannot be undone"
        And then I see another confirm prompt "Confirm? (y/N)", I answer yes

        Then ask for confirm "Confirm?" with help "This action cannot be undone", receive yes
```

### Password

#### Answer

Expect to see a Password prompt and answer it.

Pattern: `(?:(?:get)|(?:see))s? a(?:nother)? password prompt "([^"]*)".* answers? "([^"]*)"`

Example:

```gherkin
    Scenario: Receive an answer
        Given I see a password prompt "Enter password:", I answer "123456"

        Then ask for password "Enter password:", receive "123456"
```

#### Interrupt

Expect to see a Password prompt and interrupt (^C).

Pattern: `(?:(?:get)|(?:see))s? a(?:nother)? password prompt "([^"]*)".* interrupts?`

Example:

```gherkin
    Scenario: Interrupted
        Given I see a password prompt "Enter password:", I interrupt

        Then ask for password "Enter password:", get interrupted
```

#### With Help

Expect to see a Password prompt, ask for help and then expect to see a Help message.

Pattern: `(?:(?:get)|(?:see))s? a(?:nother)? password prompt "([^"]*)".* asks? for help and sees? "([^"]*)"`

Example:

```gherkin
    Scenario: With help and receive an answer
        Given I see a password prompt "Enter password: [? for help]", I ask for help and see "It is a secret"
        And then I see another password prompt "Enter password:", I answer "123456"

        Then ask for password "Enter password:" with help "It is a secret", receive "123456"
```

## Examples

- Register for
  injection: https://github.com/nhatthm/surveydog/blob/33d3788186947856b8fbf11e94604c450f6af6e4/features/bootstrap/godog_test.go#L49
- Inject: https://github.com/nhatthm/surveydog/blob/33d3788186947856b8fbf11e94604c450f6af6e4/features/bootstrap/survey.go#L36-L41
- Use: https://github.com/nhatthm/surveydog/blob/33d3788186947856b8fbf11e94604c450f6af6e4/features/bootstrap/survey.go#L57

Full suite: https://github.com/nhatthm/surveydog/tree/master/features

## Donation

If this project help you reduce time to develop, you can give me a cup of coffee :)

### Paypal donation

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;or scan this

<img src="https://user-images.githubusercontent.com/1154587/113494222-ad8cb200-94e6-11eb-9ef3-eb883ada222a.png" width="147px" />
