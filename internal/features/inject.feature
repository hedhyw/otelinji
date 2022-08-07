Feature: Adding opentracing definition

    Scenario: User wants to add opentelemetry layers
        When the user provides an input file `inject.in.test`
        And runs the application
        Then the output equals to the content of the file `inject.exp.test`
