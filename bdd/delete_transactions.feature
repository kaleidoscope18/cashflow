Feature: Transactions
    As a customer
        I need to be able to delete my transactions

    Scenario: Deleting a non-recurring transaction
        Given there is a chequing account with transactions
        When I remove a transaction
        Then I should see remaining transactions

    Scenario: Deleting a recurring transaction
        When I add a recurring transaction to it on "2030-09-20" until "2031-10-31"
        And I delete the recurring transaction on "2030-12-15"
        Then I should not see the recurring transaction on "2030-12-15" but see the rest

    Scenario: Deleting a transaction from a plan
        When I add a plan to "2030-09-20" until "2031-10-31"
        And I delete the transaction
        Then I should have an error since the transaction is part of a plan
