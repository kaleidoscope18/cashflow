Feature: Balances
    As a customer
            I need to be able to set an end-of-day balance for any given day
            And it will adjust the account's balances for each transaction

    Scenario: Adding a balance
        When I add a balance to chequing account
        Then it should be in balances list

    Scenario: Adding a balance without a date
        When I add a balance without a date to chequing account
        Then the new balance should have today's date

    Scenario: Listing transactions with balances
        Given there is a chequing account with a transaction
        When I add a balance to it
        Then I should be able to see the transactions with the right balance
