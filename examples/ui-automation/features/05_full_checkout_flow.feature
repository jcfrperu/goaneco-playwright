Feature: Full checkout flow

  As a shopper
  I want to complete a purchase end-to-end
  So that I receive an order confirmation

  Background:
    Given I navigate to the SauceDemo login page
    And I log in as "standard_user" with password "secret_sauce"
    And I am on the inventory page

  Scenario: Complete a purchase and receive order confirmation
    When I add "Sauce Labs Backpack" to the cart
    And I navigate to the cart page
    And I click the Checkout button
    And I fill in my first name "John", last name "Doe", and postal code "10001"
    And I click the Continue button
    Then the order overview should show 1 item
    And the subtotal label should start with "Item total:"
    And the total label should start with "Total:"
    When I click the Finish button
    Then I should see the confirmation message "Thank you for your order!"
