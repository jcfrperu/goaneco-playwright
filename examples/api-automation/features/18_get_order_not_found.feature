Feature: Get a non-existent order
  As an API consumer
  I want to receive a clear error when requesting an order that does not exist
  So that I can handle missing resources gracefully

  Scenario: Requesting a non-existent order returns 404
    When I send GET /store/order/999999
    Then the response status should be 404
    And the response should not be OK
