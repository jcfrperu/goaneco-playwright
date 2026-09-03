Feature: Delete a store order
  As an API consumer
  I want to cancel and remove an order
  So that it is no longer tracked in the system

  Scenario: Successfully delete an order and verify it is gone
    Given an order exists in the store
    When I send DELETE /store/order/{id} with the order's ID
    Then the response status should be 200
    And a subsequent GET /store/order/{id} should return 404
