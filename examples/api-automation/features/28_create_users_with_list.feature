Feature: Create users with list
  As an API consumer
  I want to create multiple users in a single request using a list
  So that batch user creation is supported via the list endpoint

  Scenario: Successfully create multiple users via POST /user/createWithList
    Given I have a list of 2 user objects
    When I send POST /user/createWithList with the user list
    Then the response status should be 200
    And each user should be retrievable by GET /user/{username}
