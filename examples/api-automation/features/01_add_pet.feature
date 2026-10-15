Feature: Add a new pet
  As an API consumer
  I want to add a new pet to the store
  So that it becomes available for purchase

  Scenario: Successfully add a pet with valid data
    Given I have pet data with name "Buddy" and status "available"
    When I send POST /pet with the pet payload
    Then the response status should be 200
    And the response body should contain a valid pet ID
    And the pet name should be "Buddy"
    And the pet status should be "available"
