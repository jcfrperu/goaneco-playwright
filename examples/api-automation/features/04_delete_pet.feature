Feature: Delete a pet
  As an API consumer
  I want to remove a pet from the store
  So that it is no longer available

  Scenario: Successfully delete a pet and verify it is gone
    Given a pet exists in the store
    When I send DELETE /pet/{id} with the pet's ID
    Then the response status should be 200
    And a subsequent GET /pet/{id} should return 404
