Feature: Logout via burger menu

  As an authenticated shopper
  I want to log out through the navigation menu
  So that my session is terminated and the login page is shown

  Background:
    Given I navigate to the SauceDemo login page
    And I log in as "standard_user" with password "secret_sauce"
    And I am on the inventory page

  Scenario: Logout returns user to the login page
    When I open the burger navigation menu
    And I click the Logout link
    Then I should be redirected to the login page
    And the login form should be visible
