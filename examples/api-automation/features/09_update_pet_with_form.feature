Feature: Update a pet using form data
  As an API consumer
  I want to update a pet's name and status via form-encoded data
  So that I can use standard HTML form submission semantics

  Scenario: Successfully update a pet with form-encoded data
    Given a pet exists in the store
    When I send POST /pet/{id} with Content-Type application/x-www-form-urlencoded and name "UpdatedName" and status "pending"
    Then the response status should be 200
