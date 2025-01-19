Feature: Transactions
    As an authenticated customer
            I need to be able to list all my transactions

    Scenario: List transactions
        Given I created a plan starting on 1111-01-11
            And I created a transaction on 1111-01-12
            And I created a recurring transaction on 1111-01-13
        When I list the transactions for this plan
        Then I should see the transactions correctly

    Scenario: Transactions with recurrency on weekdays only
        Given I created a recurring transaction on the 11th once a month on weekdays only
            And I want transactions to be postponed to the next monday if it falls on a weekend
        Then I should see the transactions correctly

    Scenario: Transactions and balances together
        Given I create a realistic budget starting on 2025-01-01
            And I set the balance on 2025-01-19 to 318$
            And I set the balance on 2025-04-16 to 1$
        Then I should see the transactions correctly

    Scenario: Automatically reordering manual transactions
        Given I created a manual withdrawal on 2025-01-01
            And I created another smaller manual withdrawal on the same day
            And I created an automatic transaction on the same day
            And I created another smaller automatic transaction on the same day
            And I created a manual deposit on the same day
        When I list the transactions for this plan
        Then I should see the transactions correctly ordered

