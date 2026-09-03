Feature: Get store inventory
  As an API consumer
  I want to retrieve the current inventory
  So that I can see how many pets are in each status category

  Scenario: Successfully retrieve the store inventory
    When I send GET /store/inventory
    Then the response status should be 200
    And the response body should be a non-empty map of status counts
