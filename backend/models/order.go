package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	RetailerID uint
	ProductID  uint
	Quantity   int
	Status     string // "pending", "completed"
}
