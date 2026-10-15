package models

// Order represents a store order for a pet.
type Order struct {
	ID       int64  `json:"id,omitempty"`
	PetID    int64  `json:"petId"`
	Quantity int    `json:"quantity"`
	ShipDate string `json:"shipDate,omitempty"`
	Status   string `json:"status,omitempty"` // placed | approved | delivered
	Complete bool   `json:"complete,omitempty"`
}
