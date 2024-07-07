Feature: Transactions
    As a customer
    I need to be able to manage my transactions

    Scenario: Adding a transaction
        Given there is an account
        When I add a transaction to it
        And I list the transactions
        Then I should see the new transaction
    
    Scenario: Removing a transaction
        Given there is an existing transaction in chequing account
        When I remove it
        And I list the transactions
        Then I should not see the new transaction

    Scenario: Adding a recurring transaction
        Given there is an account
        When I add a recurring transaction to it
        And I list the transactions
        Then I should be able to see all recurring transactions
