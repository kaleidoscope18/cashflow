Feature: Plans
    As an authenticated customer
            I need to be able to create financial plans

    Scenario: List plans
        Given I created plans
        When I request the list of plans
        Then I should receive a list of plans

    Scenario: Create a plan between accounts
        Given I named my plan "ping-pong table"
            And I need a total amount of 934.34 dollars
            And I start this plan on 2025-01-01
            And the money comes from my chequing account
            And the money goes to my savings account
            And I want to transfer money every 3 days
        When I create a plan
        Then I should be able to list the transactions for this plan

    Scenario: Create a plan with external account
        Given I named my plan "ping-pong table"
            And I need a total amount of 934.34 dollars
            And I start this plan on 2025-01-01
            And the money comes from my chequing account
            And the money goes to an external account named "piggy bank"
            And I want to transfer money every 3 days
        When I create a plan
        Then I should be able to list the transactions for this plan

    Scenario: Create a one-time plan
        Given I named my plan "jelly beans"
            And I need a total amount of 1 dollars
            And I start this plan on 2025-01-01
            And the money comes from my savings account
            And the money goes to my chequing account
            And I want to transfer money only once
        When I create a plan
        Then I should be able to list the single transaction for this plan
