package model

import "time"

// Vacation represents the Vacation model.
type Vacation struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	ApprovedBy *string    `json:"approved_by"`
	From       time.Time  `json:"from"`
	To         time.Time  `json:"to"`
	CreatedAt  *time.Time `json:"created_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

// Copy returns a deep copy.
func (v *Vacation) Copy() *Vacation {
	var approvedBy *string
	if v.ApprovedBy != nil {
		apBy := *v.ApprovedBy
		approvedBy = &apBy
	}
	var createdAt, deletedAt *time.Time
	if v.CreatedAt != nil {
		ct := time.Unix(0, v.CreatedAt.UnixNano())
		createdAt = &ct
	}
	if v.DeletedAt != nil {
		dt := time.Unix(0, v.DeletedAt.UnixNano())
		deletedAt = &dt
	}
	return &Vacation{
		ID:         v.ID,
		UserID:     v.UserID,
		ApprovedBy: approvedBy,
		From:       v.From,
		To:         v.To,
		CreatedAt:  createdAt,
		DeletedAt:  deletedAt,
	}
}
