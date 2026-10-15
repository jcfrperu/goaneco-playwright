Feature: Update a pet
  As an API consumer
  I want to update an existing pet's information
  So that its details remain accurate

  Scenario: Successfully update a pet's status to "sold"
    Given a pet exists in the store with status "available"
    When I send PUT /pet with the updated pet data setting status to "sold"
    Then the response status should be 200
    And the response should show the pet status as "sold"
