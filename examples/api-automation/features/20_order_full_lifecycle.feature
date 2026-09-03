Feature: Order full lifecycle
  As an API consumer
  I want to place, retrieve, delete and confirm deletion of an order in sequence
  So that I can verify the complete order management workflow

  Scenario: Complete order lifecycle — place, retrieve, delete, confirm deletion
    Given I place a new order for pet ID 1 with quantity 2
    When I retrieve the order by its ID
    Then the order is found with quantity 2
    When I delete the order
    Then a subsequent retrieval returns 404
