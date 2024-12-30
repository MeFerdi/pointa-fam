package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email" gorm:"unique"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
}

// DeleteUser removes a user from the database by ID.
func DeleteUser(db *gorm.DB, id uint) error {
	return db.Delete(&User{}, id).Error
}
