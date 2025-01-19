Feature: Transactions
    As an authenticated customer
            I need to be able to list all my transactions

    Scenario: List transactions
        Given I created a plan starting on 1111-01-11
            And I created a transaction on 1111-01-12
            And I created a recurring transaction on 1111-01-13
        When I list the transactions for this plan
        Then I should see the transactions correctly
