Feature: Product detail page

  As a shopper
  I want to view the detail page of a product
  So that I can read its full description and confirm its price before buying

  Background:
    Given I navigate to the SauceDemo login page
    And I log in as "standard_user" with password "secret_sauce"
    And I am on the inventory page

  Scenario: View product details and navigate back to inventory
    When I click on the product "Sauce Labs Backpack"
    Then I should be on the product detail page
    And the product name should be "Sauce Labs Backpack"
    And the product price should be "$29.99"
    And the product description should not be empty
    When I click the back-to-products button
    Then I should see 6 products on the inventory page
