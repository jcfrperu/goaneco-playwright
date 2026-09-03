Feature: Login with invalid credentials
  As an API consumer
  I want to be rejected when using wrong credentials
  So that unauthorized access is prevented

  Scenario: Login with invalid credentials returns a non-2xx status
    Given I have invalid credentials
    When I send GET /user/login with the invalid credentials
    Then the response status should not be in the 2xx range
