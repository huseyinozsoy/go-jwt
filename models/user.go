package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// User represents a User schema
type User struct {
	Base
	Email    string `json:"email" gorm:"unique;size:100;not null"`
	Username string `json:"username" gorm:"unique;size:100;not null"`
	Password string `json:"password" gorm:"unique;size:100;not null"`
}

// UserErrors represent the error format for user routes
type UserErrors struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}

// Claims represent the structure of the JWT token
type Claims struct {
	jwt.StandardClaims
	ID uuid.UUID `gorm:"primaryKey`
}
