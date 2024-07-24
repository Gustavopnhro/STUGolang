package models

type Product struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
	ID    int     `json:"id"`
}
