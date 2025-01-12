Feature: Edit transactions
    As a customer
            I need to be able to edit my transactions

    Scenario: Editing a transaction
        When I add a transaction
            And I edit the transaction
        Then I should see edited transaction

    Scenario: Editing a recurring transaction
        When I add a recurring transaction between "2022-09-20" until "2023-10-31"
            And I edit the transaction
        Then I should see edited transactions

    Scenario: Editing a recurring transaction from given date
        When I add a recurring transaction to it on "2022-09-20" until "2023-10-31"
            And I edit the transaction from "2022-12-15"
        Then I should see unchanged transactions before "2022-12-15" and the rest should change

    Scenario: Editing a recurring transaction from given date but no date is provided
        When I add a recurring transaction to it on "2022-09-20" until "2023-10-31"
            And I edit the transaction info from date without a date
        Then I should have an expected date error

    Scenario: Editing a recurring transaction only on date but no date is provided
        When I add a recurring transaction to it on "2022-09-20" until "2023-10-31"
            And I edit the transaction info on date only without a date
        Then I should have an expected date error

    Scenario: Editing a recurring transaction only on date
        When I add a recurring transaction to it on "2022-09-20" until "2023-10-31"
            And I edit the transaction on "2022-12-15"
        Then I should see unchanged transactions except transaction on "2022-12-15"
