Feature: Transactions
    As a customer
    I need to be able to manage my transactions

    Scenario: Adding a transaction
        Given there is a chequing account
        When I add a transaction to it
        Then I should be able to see the transactions
    
    Scenario: Removing a transaction
        Given there is an existing transaction in chequing account
        When I remove it
        Then it should be removed

    Scenario: Listing transactions
        Given there is a chequing account
        When I list the transactions
        Then I should be able to see the transactions

    Scenario: Adding a recurring transaction
        Given there is a chequing account
        When I add a recurring transaction to it
        And I list the transactions between two dates
        Then I should be able to see all recurring transactions