Feature: Find pets by status "pending"
  As an API consumer
  I want to find all pending pets
  So that I can track pets awaiting processing

  Scenario: Filtering by "pending" returns a valid array of pets
    When I send GET /pet/findByStatus?status=pending
    Then the response status should be 200
    And the response body should be a valid JSON array
    And every pet in the response should have status "pending"
