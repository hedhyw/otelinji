Feature: Inject
    Inject OpenTelemetry code to the functions of the file.

    Scenario: User wants to add opentelemetry layers
        Given the user provides an input file `inject.exp.go.txt`
        When he runs the application
        Then the output equals to the content of the file `inject.in.go.txt`
