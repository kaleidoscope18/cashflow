Feature: Accounts
    As an authenticated customer
            I need to be able to manage my accounts

    Scenario: Filtering transactions from specific account
        When I add a transaction on "2022-10-21" to savings account
        When I add a transaction on "2022-10-22" to chequing account
        Then I should see the new transaction only on specified account
