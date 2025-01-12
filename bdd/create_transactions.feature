Feature: Transactions
    As an authenticated customer
            I need to be able to create my transactions

    Scenario: Creating a transaction
        When I add a transaction
        Then I should see the new transaction

    Scenario: Creating a recurring transaction
        When I add a recurring transaction to it on "2022-09-20" until "2023-10-31"
        Then I should be able to see all recurring transactions

    Scenario: Creating a mix of recurring and non recurring transactions
        When I add a recurring transaction to it on "2022-09-20" until "2023-10-31"
            And I add a transaction on "2022-12-13"
            And I add a transaction on "2023-01-18"
        Then I should be able to see all transactions in order
