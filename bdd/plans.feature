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
            And I start this plan on 2076-01-01
            And the money comes from my chequing account
            And the money goes to my savings account
            And I want to transfer money every 3 days, two times
        When I create a plan
        Then I should be able to list the transactions for this plan

    Scenario: Create a plan to an external account
        Given I named my plan "popcorn maker"
            And I need a total amount of 300 dollars
            And I start this plan on 2084-01-01
            And the money comes from my chequing account
            And the money goes to an external account named "piggy bank"
            And I want to transfer money every 2 days, ten times
        When I create a plan
        Then I should be able to list the transactions for this plan

    Scenario: Create a plan from an external account
        Given I named my plan "smoothie machine"
            And I need a total amount of 2 dollars
            And I start this plan on 2081-01-01
            And the money comes from an external account named "piggy bank"
            And the money goes to my chequing account
            And I want to transfer money every 2 days, two times
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

    Scenario: Delete a plan
        Given I created a plan
        When I delete the plan
        Then I should not see the plan or its transactions anymore
