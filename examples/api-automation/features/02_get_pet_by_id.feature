Feature: Get a pet by ID
  As an API consumer
  I want to retrieve a specific pet by its ID
  So that I can view its details

  Scenario: Successfully retrieve an existing pet by ID
    Given a pet exists in the store
    When I send GET /pet/{id} with the pet's ID
    Then the response status should be 200
    And the response should contain the pet's ID, name, and status
