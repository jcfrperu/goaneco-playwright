Feature: Find pets by multiple statuses
  As an API consumer
  I want to filter pets by more than one status at a time
  So that I can query a broader set of pets efficiently

  Scenario: Filtering by "available" and "pending" returns only pets with those statuses
    When I send GET /pet/findByStatus?status=available&status=pending
    Then the response status should be 200
    And the response should contain at least one pet
    And every pet in the response should have status "available" or "pending"
