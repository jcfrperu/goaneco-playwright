Feature: Create a user
  As an API consumer
  I want to create a new user account
  So that they can interact with the Petstore

  Scenario: Successfully create a user and retrieve the account
    Given I have user data with username "goaneco-create21" and email "goaneco-create21@example.com"
    When I send POST /user with the user payload
    Then the response status should be 200
    And a subsequent GET /user/goaneco-create21 should return 200 with matching fields
