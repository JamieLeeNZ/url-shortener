package models

import "time"

type User struct {
	ID        string    `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name,omitempty" db:"name"`
	Picture   string    `json:"picture,omitempty" db:"picture"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
}
