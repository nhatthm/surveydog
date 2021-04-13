Feature: Confirm

    Scenario: With help and receive a yes
        Given the app gets a confirm prompt "Confirm? [? for help] (y/N)", it will ask for help and sees "This action cannot be undone"
        And then the app gets another confirm prompt "Confirm? (y/N)", it will answer yes

        Then ask for confirm "Confirm?" with help "This action cannot be undone", receive yes

    Scenario: With help and receive a no
        Given the app gets a confirm prompt "Confirm? [? for help] (y/N)", it will ask for help and sees "This action cannot be undone"
        And then the app gets another confirm prompt "Confirm? (y/N)", it will answer no

        Then ask for confirm "Confirm?" with help "This action cannot be undone", receive no

    Scenario: Without help and receive a yes
        Given the app gets a confirm prompt "Confirm? (y/N)", it will answer yes

        Then ask for confirm "Confirm?", receive yes

    Scenario: Without help and receive a no
        Given the app gets a confirm prompt "Confirm? (y/N)", it will answer no

        Then ask for confirm "Confirm?", receive no

    Scenario: Interrupted
        Given the app gets a confirm prompt "Confirm? (y/N)", it will interrupt

        Then ask for confirm "Confirm? (y/N)", get interrupted
