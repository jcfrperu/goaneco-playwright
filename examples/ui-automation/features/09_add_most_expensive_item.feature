Feature: Add the most expensive item to cart

  As a high-end shopper
  I want to sort by price descending and add the first result
  So that I can purchase the most premium product available

  Background:
    Given I navigate to the SauceDemo login page
    And I log in as "standard_user" with password "secret_sauce"
    And I am on the inventory page

  Scenario: Add the most expensive product to the cart
    When I select the sort option "Price (high to low)"
    Then the first product displayed should be "Sauce Labs Fleece Jacket"
    And the first product price should be "$49.99"
    When I add "Sauce Labs Fleece Jacket" to the cart
    And I navigate to the cart page
    Then the cart should contain "Sauce Labs Fleece Jacket"
    And the cart should display the price "$49.99"
