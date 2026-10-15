Feature: Update a user
  As an API consumer
  I want to update a user's profile information
  So that their account details remain current

  Scenario: Successfully update a user's email address
    Given a user exists in the system
    When I send PUT /user/{username} with an updated email address
    Then the response status should be 200
    And a subsequent GET /user/{username} should return the updated email
