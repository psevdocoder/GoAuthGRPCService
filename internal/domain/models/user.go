package models

type User struct {
	ID           uint
	Username     string
	PasswordHash []byte
	Role         int
}
