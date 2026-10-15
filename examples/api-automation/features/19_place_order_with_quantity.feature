Feature: Place an order with a specific quantity
  As an API consumer
  I want to order multiple units of a pet
  So that bulk purchases are supported

  Scenario: Successfully place an order with quantity=5
    Given I have order data for pet ID 1 with quantity 5
    When I send POST /store/order with the order payload
    Then the response status should be 200
    And the response body should show quantity as 5
