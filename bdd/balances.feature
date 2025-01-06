Feature: Balances
    As a customer
    I need to be able to set an end-of-day balance for any given day
    And it will adjust the account's balances for each transaction

  Scenario: Adding a balance
    Given login as "authenticated"
    Given there is a chequing account
    When I add a balance to it
    And I list the balances
    Then it should be in balances list

  Scenario: Adding a balance without a date
    Given login as "authenticated"
    Given there is a chequing account
    When I add a balance without a date to it
    And I list the balances
    Then the new balance should have today's date

  Scenario: Listing transactions with balances
    Given login as "authenticated"
    Given there is a chequing account with transactions
    When I add a balance to it
    And I list the transactions
    Then I should be able to see the transactions with the right balances
