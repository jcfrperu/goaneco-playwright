Feature: Sort products by price low to high

  As a budget-conscious shopper
  I want to sort products from cheapest to most expensive
  So that I can find the best deal quickly

  Background:
    Given I navigate to the SauceDemo login page
    And I log in as "standard_user" with password "secret_sauce"
    And I am on the inventory page

  Scenario: First product after sorting is the cheapest
    When I select the sort option "Price (low to high)"
    Then the first product displayed should be "Sauce Labs Onesie"
    And the first product price should be "$7.99"
