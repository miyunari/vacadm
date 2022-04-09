package model

import "time"

type Team struct {
	ID        string     `json:"id"`
	OwnerID   string     `json:"owner_id"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (t *Team) Copy() *Team {
	var createdAt, deletedAt, updatedAt *time.Time
	if t.CreatedAt != nil {
		ct := time.Unix(0, t.CreatedAt.UnixNano())
		createdAt = &ct
	}
	if t.DeletedAt != nil {
		dt := time.Unix(0, t.DeletedAt.UnixNano())
		deletedAt = &dt
	}
	if t.UpdatedAt != nil {
		ct := time.Unix(0, t.UpdatedAt.UnixNano())
		updatedAt = &ct
	}
	return &Team{
		ID:        t.ID,
		OwnerID:   t.OwnerID,
		Name:      t.Name,
		CreatedAt: createdAt,
		DeletedAt: deletedAt,
		UpdatedAt: updatedAt,
	}
}
