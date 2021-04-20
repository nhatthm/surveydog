Feature: Password

    Scenario: With help and receive an answer
        Given I see a password prompt "Enter password: [? for help]", I ask for help and see "It is a secret"
        And then I see another password prompt "Enter password:", I answer "123456 with help"

        Then ask for password "Enter password:" with help "It is a secret", receive "123456 with help"

    Scenario: Without help and receive an answer
        Given I see a password prompt "Enter password:", I answer "123456"

        Then ask for password "Enter password:", receive "123456"

    Scenario: Interrupted
        Given I see a password prompt "Enter password:", I interrupt

        Then ask for password "Enter password:", get interrupted

    Scenario: Multiple Password prompts
        Given I see a password prompt "Enter username:", I answer "user@example.org"
        And I see a password prompt "Enter password:", I answer "123456"

        Then ask for password "Enter username:", receive "user@example.org"
        And ask for password "Enter password:", receive "123456"

        And console output is:
        """
        ? Enter username: ****************
        ? Enter password: ******
        """
