package model

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model           // Adds some metadata fields to the table
	ID         uuid.UUID `gorm:"type:uuid"` // Explicitly specify the type to be uuid
	Title      string
	SubTitle   string
	Text       string
}
type CredentialsJsonLess struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	JwtToken string `json:"jwt-token"`
}
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
