Feature: Login with valid credentials
  As an API consumer
  I want to authenticate with valid credentials
  So that I receive a session token

  Scenario: Successful login returns 200 with a session token
    Given I have valid credentials username "test" and password "abc123"
    When I send GET /user/login with the credentials
    Then the response status should be 200
    And the response body should not be empty
