Feature: Transactions
    As a customer
            I need to be able to delete my transactions

    Scenario: Deleting a non-recurring transaction
        Given there is a chequing account with transactions
        When I remove a transaction
        Then I should see remaining transactions

    Scenario: Deleting a recurring transaction
        When I add a recurring transaction to it on "2028-09-20" until "2029-10-31"
            And I delete the recurring transaction on "2028-12-15"
        Then I should not see the recurring transaction on "2028-12-15" but see the rest
