Feature: Remove item from cart

  As a shopper
  I want to remove an item from my cart
  So that I can adjust my order before checking out

  Background:
    Given I navigate to the SauceDemo login page
    And I log in as "standard_user" with password "secret_sauce"
    And I am on the inventory page

  Scenario: Remove one item from a cart with 3 items
    When I add "Sauce Labs Backpack" to the cart
    And I add "Sauce Labs Bike Light" to the cart
    And I add "Sauce Labs Bolt T-Shirt" to the cart
    And I navigate to the cart page
    Then I should see 3 items in the cart
    When I remove "Sauce Labs Bike Light" from the cart
    Then I should see 2 items in the cart
    And the cart should not contain "Sauce Labs Bike Light"
