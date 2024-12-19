package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email       string `json:"email" gorm:"unique"`
	Password    string `json:"password"`
	Role        string `json:"role"`
	Name        string `json:"name"`
	Location    string `json:"location"`
	ContactInfo string `json:"contact_info"`
}
