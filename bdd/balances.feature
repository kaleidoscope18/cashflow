Feature: Balances
    As a customer
    I need to be able to set an account balance at the end of any given day

    Scenario: Setting the balance
        Given there is an account
        When I add a balance to it
        Then I should be able to list the balances

    Scenario: Listing transactions
        Given there is an account
        When I add a balance to it
        And I list the transactions
        Then I should be able to see the transactions with the right balances