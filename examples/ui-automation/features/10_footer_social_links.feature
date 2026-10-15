Feature: Footer social media links

  As a visitor
  I want to see social media links in the footer
  So that I can follow SauceDemo on social platforms

  Background:
    Given I navigate to the SauceDemo login page
    And I log in as "standard_user" with password "secret_sauce"
    And I am on the inventory page

  Scenario: Footer contains Twitter, Facebook, and LinkedIn links
    Then the footer should contain a visible link to Twitter
    And the Twitter link should have a non-empty href
    And the footer should contain a visible link to Facebook
    And the Facebook link should have a non-empty href
    And the footer should contain a visible link to LinkedIn
    And the LinkedIn link should have a non-empty href
