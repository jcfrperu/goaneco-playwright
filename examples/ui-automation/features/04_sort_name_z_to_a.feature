Feature: Sort products by name Z to A

  As a shopper
  I want to sort products in reverse alphabetical order
  So that I can browse starting from Z

  Background:
    Given I navigate to the SauceDemo login page
    And I log in as "standard_user" with password "secret_sauce"
    And I am on the inventory page

  Scenario: Products are sorted in reverse alphabetical order
    When I select the sort option "Name (Z to A)"
    Then the first product displayed should be "Test.allTheThings() T-Shirt (Red)"
    And the last product displayed should be "Sauce Labs Backpack"
