Feature: Get a non-existent pet
  As an API consumer
  I want to receive a clear error when requesting a pet that does not exist
  So that I can handle missing resources gracefully

  Scenario: Requesting a non-existent pet returns 404
    When I send GET /pet/999999999
    Then the response status should be 404
    And the response should not be OK
