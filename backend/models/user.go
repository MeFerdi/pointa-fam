package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username          string   `json:"username"`
	Email             string   `json:"email" gorm:"unique"`
	Password          string   `json:"password"`
	PhoneNumber       string   `json:"phone_number"`
	Location          string   `json:"location"`
	Role              string   `json:"role"`
	FarmerID          uint     `json:"farmer_id"`
	Farmer            Farmer   `json:"farmer" gorm:"foreignKey:FarmerID"`
	RetailerID        uint     `json:"retailer_id"`
	Retailer          Retailer `json:"retailer" gorm:"foreignKey:RetailerID"`
	ProfilePictureUrl string   `json:"profile_picture_url"`
}

// CreateUser inserts a new user into the database
func (u *User) CreateUser(db *gorm.DB) error {
	return db.Create(u).Error
}

// GetUserByEmail retrieves a user by email from the database
func GetUserByEmail(db *gorm.DB, email string) (User, error) {
	var user User
	err := db.Where("email = ?", email).First(&user).Error
	return user, err
}

// DeleteUser removes a user from the database by ID.
func DeleteUser(db *gorm.DB, id uint) error {
	return db.Delete(&User{}, id).Error
}
