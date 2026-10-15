Feature: Delete a user
  As an API consumer
  I want to remove a user account from the system
  So that their data is no longer accessible

  Scenario: Successfully delete a user and verify the account is gone
    Given a user exists in the system
    When I send DELETE /user/{username}
    Then the response status should be 200
    And a subsequent GET /user/{username} should return 404
