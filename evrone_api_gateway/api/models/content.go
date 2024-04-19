package models

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name" example:"Laptop"`
	Description string  `json:"description" example:"MacBook Pro 9999"`
	CategoryID  string  `json:"category_id" example:"9ffd6a49-3c85-4561-bbc0-0d1445741576"`
	UnitPrice   float64 `json:"unit_price" example:"99.9"`
}
