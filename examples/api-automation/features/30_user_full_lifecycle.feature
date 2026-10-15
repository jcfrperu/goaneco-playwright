Feature: User full lifecycle
  As an API consumer
  I want to create, login, retrieve, update, delete and confirm deletion of a user
  So that I can verify the complete user management workflow

  Scenario: Complete user lifecycle — create, login, retrieve, update, delete, confirm deletion
    Given I create a new user with username "goaneco-lifecycle30"
    When I login with the user's credentials
    Then the login response status should be 200
    When I retrieve the user by username
    Then the user is found with the correct username
    When I update the user's email
    Then the update is successful
    When I delete the user
    Then a subsequent retrieval returns 404
