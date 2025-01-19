Feature: Transactions details
    As a customer
            I need to be able to see the details of my transactions

    Scenario: Getting transactions statuses
        Given I add a transaction today and another later
        Then I should be able to see the correct statuses for each transaction

    Scenario: Marking a transaction as automatic or manual
        Given I create two transactions
            And I mark one as automatic
            And I mark the other as manual
        Then I should be able to see the correct details for each transaction
