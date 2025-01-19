Feature: Transactions

    Scenario: List transactions
        Given I created a plan starting on 1111-01-11
            And I created a transaction on 1111-01-11
            And I created a recurring transaction on 1111-01-11
        When I list the transactions for this plan
        Then I should see the transactions correctly
