Feature: Add items to cart

  As a shopper
  I want to add multiple products to my cart
  So that I can review them before purchasing

  Background:
    Given I navigate to the SauceDemo login page
    And I log in as "standard_user" with password "secret_sauce"
    And I am on the inventory page

  Scenario: Add 3 products and verify them in the cart
    When I add "Sauce Labs Backpack" to the cart
    And I add "Sauce Labs Bike Light" to the cart
    And I add "Sauce Labs Bolt T-Shirt" to the cart
    Then the cart badge should show "3"
    When I navigate to the cart page
    Then I should see 3 items in the cart
    And the cart should contain "Sauce Labs Backpack"
    And the cart should contain "Sauce Labs Bike Light"
    And the cart should contain "Sauce Labs Bolt T-Shirt"
    And each item should display a price
