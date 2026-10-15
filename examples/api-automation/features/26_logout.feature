Feature: Logout
  As an API consumer
  I want to end my session
  So that I am securely logged out

  Scenario: Logout returns 200
    When I send GET /user/logout
    Then the response status should be 200
