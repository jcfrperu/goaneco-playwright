Feature: Pet full lifecycle
  As an API consumer
  I want to create, read, update, and delete a pet in sequence
  So that I can verify the complete pet management workflow

  Scenario: Complete pet lifecycle — create, retrieve, update, delete, confirm deletion
    Given I create a new pet named "LifecyclePet" with status "available"
    When I retrieve the pet by its ID
    Then the pet is found with status "available"
    When I update the pet's status to "sold"
    Then the update is successful
    When I delete the pet
    Then a subsequent retrieval returns 404
