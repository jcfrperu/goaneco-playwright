package models

// Category represents a pet category.
type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Tag represents a pet tag.
type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Pet represents a pet in the store.
type Pet struct {
	ID        int64     `json:"id,omitempty"`
	Category  *Category `json:"category,omitempty"`
	Name      string    `json:"name"`
	PhotoURLs []string  `json:"photoUrls"`
	Tags      []Tag     `json:"tags,omitempty"`
	Status    string    `json:"status,omitempty"` // available | pending | sold
}
