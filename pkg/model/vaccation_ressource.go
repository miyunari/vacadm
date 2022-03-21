package model

import "time"

type VaccationRessources struct {
	ID         string    `json:"id"`
	YearlyDays int       `json:"yearly_days"`
	From       time.Time `json:"from"`
	To         time.Time `json:"to"`
	CreatedAt  time.Time `json:"created_at"`
	DeletedAt  time.Time `json:"deleted_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
