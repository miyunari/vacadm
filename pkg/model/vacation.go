package model

import "time"

type Vacation struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	ApprovedBy *User      `json:"approved_by"`
	From       time.Time  `json:"from"`
	To         time.Time  `json:"to"`
	CreatedAt  *time.Time `json:"created_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

func (v Vacation) Copy() *Vacation {
	return &v
}
