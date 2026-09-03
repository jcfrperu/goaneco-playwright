Feature: Login and verify inventory

  As a shopper
  I want to log in with valid credentials
  So that I can see the full product catalog

  Background:
    Given I navigate to the SauceDemo login page
    And I log in as "standard_user" with password "secret_sauce"

  Scenario: Inventory page displays 6 products after login
    When I am on the inventory page
    Then I should see 6 products listed
    And the sort dropdown should be visible
