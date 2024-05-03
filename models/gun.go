package model

import (
	"time"
)

type Gun struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Manufacturer string    `json:"manufacturer"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
