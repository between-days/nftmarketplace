package models

import "gorm.io/gorm"

type Listing struct {
	gorm.Model
	Seller   string
	Price    float64
	ImageURL string
}
