package model

import "time"

// VacationResource represents the VacationResource model.
type VacationResource struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	YearlyDays int        `json:"yearly_days"`
	From       time.Time  `json:"from"`
	To         time.Time  `json:"to"`
	CreatedAt  *time.Time `json:"created_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

// Copy returns a deep copy.
func (v *VacationResource) Copy() *VacationResource {
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
	return &VacationResource{
		ID:         v.ID,
		UserID:     v.UserID,
		YearlyDays: v.YearlyDays,
		From:       v.From,
		To:         v.To,
		CreatedAt:  createdAt,
		DeletedAt:  deletedAt,
		UpdatedAt:  updatedAt,
	}
}
