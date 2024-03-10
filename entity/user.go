package entity

import "time"

type User struct {
	ID             string
	FullName       string
	PhoneNumber    string
	HashedPassword string
	CreatedAt      time.Time
	UpdatedAt      *time.Time
}
