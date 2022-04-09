package model

import "time"

type VacationRequest struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	To        time.Time  `json:"to"`
	From      time.Time  `json:"from"`
	CreatedAt *time.Time `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (v *VacationRequest) Copy() *VacationRequest {
	var createdAt, deletedAt, updatedAt *time.Time
	if v.CreatedAt != nil {
		ct := time.Unix(0, v.CreatedAt.UnixNano())
		createdAt = &ct
	}
	if v.DeletedAt != nil {
		dt := time.Unix(0, v.DeletedAt.UnixNano())
		deletedAt = &dt
	}
	if v.UpdatedAt != nil {
		ut := time.Unix(0, v.UpdatedAt.UnixNano())
		updatedAt = &ut
	}
	return &VacationRequest{
		ID:        v.ID,
		UserID:    v.UserID,
		From:      v.From,
		To:        v.To,
		CreatedAt: createdAt,
		DeletedAt: deletedAt,
		UpdatedAt: updatedAt,
	}
}
