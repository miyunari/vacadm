package model

import "time"

type User struct {
	ID        string     `json:"id"`
	ParentID  *string    `json:"parent_id"`
	TeamID    *string    `json:"team_id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (u *User) Copy() *User {
	var parentID, teamID *string
	if u.ParentID != nil {
		pID := *u.ParentID
		parentID = &pID
	}
	if u.TeamID != nil {
		tID := *u.TeamID
		teamID = &tID
	}
	var createdAt, deletedAt, updatedAt *time.Time
	if u.CreatedAt != nil {
		ct := time.Unix(0, u.CreatedAt.UnixNano())
		createdAt = &ct
	}
	if u.DeletedAt != nil {
		dt := time.Unix(0, u.DeletedAt.UnixNano())
		deletedAt = &dt
	}
	if u.UpdatedAt != nil {
		ct := time.Unix(0, u.UpdatedAt.UnixNano())
		updatedAt = &ct
	}
	return &User{
		ID:        u.ID,
		ParentID:  parentID,
		TeamID:    teamID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		CreatedAt: createdAt,
		DeletedAt: deletedAt,
		UpdatedAt: updatedAt,
	}
}
