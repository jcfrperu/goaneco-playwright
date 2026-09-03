Feature: Get a user by username
  As an API consumer
  I want to retrieve a user's profile by username
  So that I can view their account details

  Scenario: Successfully retrieve an existing user by username
    Given a user exists with a known username
    When I send GET /user/{username}
    Then the response status should be 200
    And the response body should contain the matching username
