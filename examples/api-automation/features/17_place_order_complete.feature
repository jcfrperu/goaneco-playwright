Feature: Place a completed order
  As an API consumer
  I want to place an order marked as complete
  So that it is recorded as fully processed

  Scenario: Successfully place an order with complete=true
    Given I have order data with complete set to true
    When I send POST /store/order with the order payload
    Then the response status should be 200
    And the response body should show complete as true
