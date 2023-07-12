package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	ListingID int     `json:"listingID"`
	Seller    string  `json:"seller"`
	Buyer     string  `json:"buyer"`
	Price     float64 `json:"price"`
	Status    string  `json:"string"` // "success" | "pending" | "error"
	Error     string  `json:"error"`
	// success or error timestamp
	FinishAt time.Time `json:"finishAt"`
	// misc blockchain details if necessary
}
