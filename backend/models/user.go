package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Location        string `json:"location"`
	Email           string `json:"email" gorm:"unique"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password" gorm:"-"`
	Role            string `json:"role"`
	ContactInfo     string `json:"contact_info"`
}
