Feature: Password

    Scenario: With help and receive an answer
        Given the app gets a password prompt "Enter password: [? for help]", it will ask for help and sees "It is a secret"
        And then the app gets another password prompt "Enter password:", it will answer "123456 with help"

        Then ask for password "Enter password:" with help "It is a secret", receive "123456 with help"

    Scenario: Without help and receive an answer
        Given the app gets a password prompt "Enter password:", it will answer "123456"

        Then ask for password "Enter password:", receive "123456"

    Scenario: Interrupted
        Given the app gets a password prompt "Enter password:", it will interrupt

        Then ask for password "Enter password:", get interrupted

