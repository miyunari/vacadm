package model

import "time"

type VaccationRequest struct {
	ID        string    `json:"id"`
	To        time.Time `json:"to"`
	From      time.Time `json:"from"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
