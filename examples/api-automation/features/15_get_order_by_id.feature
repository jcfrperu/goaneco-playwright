Feature: Get an order by ID
  As an API consumer
  I want to retrieve a specific order by its ID
  So that I can view its details

  Scenario: Successfully retrieve an existing order by ID
    Given an order exists in the store
    When I send GET /store/order/{id} with the order's ID
    Then the response status should be 200
    And the response should contain the matching petId and quantity
