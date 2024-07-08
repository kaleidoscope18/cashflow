Feature: Transactions
    As a customer
    I need to be able to manage my transactions

    Scenario: Adding a transaction
        Given there is an account
        When I add a transaction on "2022/10/20" to it
        And I list the transactions
        Then I should see the new transaction
    
    Scenario: Removing a transaction
        Given there is an account with transactions
        And I list the transactions
        When I remove a transaction
        And I list the transactions
        Then I should see remaining transactions

    Scenario: Adding a recurring transaction
        Given there is an account
        When I add a recurring transaction to it
        And I list the transactions
        Then I should be able to see all recurring transactions

    Scenario: Getting transactions statuses
        Given there is an account with transactions
        And I add a transaction today and another later
        And I list the transactions
        Then I should be able to see the correct statuses for each transaction

    Scenario: Editing a recurring transaction
        Given there is an account
        When I add a recurring transaction to it
        And I list the transactions
        And I edit the transaction info
        And I list the transactions
        Then I should see changed transactions

    Scenario: Editing a recurring transaction from given date
        Given there is an account
        When I add a recurring transaction to it
        And I edit the transaction info from "2022/12/15"
        And I list the transactions
        Then I should see unchanged transactions before and the rest should change

    Scenario: Editing a recurring transaction only on date
        Given there is an account
        When I add a recurring transaction to it
        And I edit the transaction info on "2022/12/15"
        And I list the transactions
        Then I should see unchanged transactions except transaction on "2022/12/15"

    Scenario: Deleting a recurring transaction
        Given there is an account
        When I add a recurring transaction to it
        Then I should see unchanged transactions except transaction on "2022/12/15"

