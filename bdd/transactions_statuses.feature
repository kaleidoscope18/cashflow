Feature: Transactions statuses
    As a customer
            I need to be able to see the statuses of my transactions

    Scenario: Getting transactions statuses
        Given I add a transaction today and another later
        Then I should be able to see the correct statuses for each transaction
