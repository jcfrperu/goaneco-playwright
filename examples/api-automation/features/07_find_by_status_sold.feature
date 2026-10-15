Feature: Find pets by status "sold"
  As an API consumer
  I want to find all sold pets
  So that I can view purchase history

  Scenario: All pets returned for status "sold" have the correct status
    Given the store contains pets with status "sold"
    When I send GET /pet/findByStatus?status=sold
    Then the response status should be 200
    And the response should contain at least one pet
    And every pet in the response should have status "sold"
