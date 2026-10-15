Feature: Find pets by status
  As an API consumer
  I want to find all pets by status
  So that I can browse pets filtered by their availability

  # All three statuses are covered by a single table-driven Go test: TestScenario05FindByStatus

  Scenario Outline: All pets returned for status "<status>" have the correct status
    Given the store contains pets with status "<status>"
    When I send GET /pet/findByStatus?status=<status>
    Then the response status should be 200
    And the response should contain at least one pet
    And every pet in the response should have status "<status>"

    Examples:
      | status    |
      | available |
      | pending   |
      | sold      |
