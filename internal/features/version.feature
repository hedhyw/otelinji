Feature: Version
    As a user I want to know the application version.

    Scenario: User wants to know the version
        Given the user provides the flag "-version"
        When they run the application
        Then the output contains the version of the server
