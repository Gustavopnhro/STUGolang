package structs

import "github.com/google/uuid"

type Product struct {
	Name  string    `json:"name"`
	Price float64   `json:"price"`
	Stock int       `json:"stock"`
	ID    uuid.UUID `json:"id"`
}
