Feature: Place a store order
  As an API consumer
  I want to place an order for a pet
  So that the pet is reserved for purchase

  Scenario: Successfully place an order and receive a valid order ID
    Given I have order data for pet ID 1 with quantity 1 and status "placed"
    When I send POST /store/order with the order payload
    Then the response status should be 200
    And the response body should contain a valid order ID
    And the order status should be "placed"
