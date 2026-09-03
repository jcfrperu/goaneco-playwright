Feature: Find pets by tags
  As an API consumer
  I want to search for pets by tag
  So that I can find pets with specific characteristics

  Scenario: Filtering by tag "goaneco" returns a valid response
    When I send GET /pet/findByTags?tags=goaneco
    Then the response status should be 200
    And the response body should be a valid JSON array
