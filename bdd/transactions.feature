Feature: Transactions
    As a customer
    I need to be able to manage my transactions

  Scenario: Adding a transaction
    Given login as "authenticated"
    Given there is an account
    When I add a transaction on "2022/10/20"
    And I list the transactions
    Then I should see the new transaction

  Scenario: Removing a transaction
    Given login as "authenticated"
    Given there is an account with transactions
    And I list the transactions
    When I remove a transaction
    And I list the transactions
    Then I should see remaining transactions

  Scenario: Getting transactions statuses
    Given login as "authenticated"
    Given there is an account with transactions
    And I add a transaction today and another later
    And I list the transactions
    Then I should be able to see the correct statuses for each transaction

  Scenario: Adding a recurring transaction
    Given login as "authenticated"
    Given there is an account
    When I add a recurring transaction to it on "2022/09/20" until "2023/10/31"
    And I list the transactions
    Then I should be able to see all recurring transactions

  Scenario: Editing a transaction
    Given login as "authenticated"
    Given there is an account
    When I add a transaction on "2022/10/15"
    And I list the transactions
    And I edit the transaction info
    And I list the transactions
    Then I should see changed transactions

  Scenario: Editing a recurring transaction
    Given login as "authenticated"
    Given there is an account
    When I add a recurring transaction to it on "2022/09/20" until "2023/10/31"
    And I list the transactions
    And I edit the transaction info
    And I list the transactions
    Then I should see changed transactions

  Scenario: Editing a recurring transaction from given date
    Given login as "authenticated"
    Given there is an account
    When I add a recurring transaction to it on "2022/09/20" until "2023/10/31"
    And I list the transactions
    And I edit the transaction from "2022/12/15"
    And I list the transactions
    Then I should see unchanged transactions before "2022/12/15" and the rest should change

  Scenario: Editing a recurring transaction from given date but no date is provided
    Given login as "authenticated"
    Given there is an account
    When I add a recurring transaction to it on "2022/09/20" until "2023/10/31"
    And I list the transactions
    And I edit the transaction info from date without a date
    Then I should have an expected date error

  Scenario: Editing a recurring transaction only on date but no date is provided
    Given login as "authenticated"
    Given there is an account
    When I add a recurring transaction to it on "2022/09/20" until "2023/10/31"
    And I list the transactions
    And I edit the transaction info on date only without a date
    Then I should have an expected date error

  Scenario: Editing a recurring transaction only on date
    Given login as "authenticated"
    Given there is an account
    When I add a recurring transaction to it on "2022/09/20" until "2023/10/31"
    And I list the transactions
    And I edit the transaction on "2022/12/15"
    And I list the transactions
    Then I should see unchanged transactions except transaction on "2022/12/15"

  @wip
  Scenario: Deleting a recurring transaction
    Given login as "authenticated"
    Given there is an account
    When I add a recurring transaction to it on "2022/09/20" until "2023/10/31"
    Then I should see unchanged transactions except transaction on "2022/12/15"
