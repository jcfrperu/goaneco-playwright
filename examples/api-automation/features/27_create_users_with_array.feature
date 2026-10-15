Feature: Create users with array
  As an API consumer
  I want to create multiple users in a single request using an array
  So that batch user creation is efficient

  Scenario: Successfully create multiple users via POST /user/createWithArray
    Given I have an array of 2 user objects
    When I send POST /user/createWithArray with the user array
    Then the response status should be 200
    And each user should be retrievable by GET /user/{username}
