package database

import "context"

// RelationDB is implemented by any structure providing all RelationDB methods.
type RelationDB interface {
	// IsParentUser verifies if the given parentID is parent of userID in some form.
	IsParentUser(ctx context.Context, userID, parentID string) (bool, error)
	// IsTeamMember verifies if the given userID belongs to teamID.
	IsTeamMember(ctx context.Context, teamID, userID string) (bool, error)
	// IsTeamOwner verifies if the given userID refers to an owner of the teamID.
	IsTeamOwner(ctx context.Context, teamID, userID string) (bool, error)
}

// NewRelationDB returns initialized RelationDB that matches
// the RelationDB interface.
func NewRelationDB(db Database) RelationDB {
	return &relationDB{
		db: db,
	}
}

type relationDB struct {
	db Database
}

// IsParentUser verifies if the given parentID is parent of userID in some form.
// Parent is recursive in this case. This means that the parent of the parent is
// also valid.
func (r *relationDB) IsParentUser(ctx context.Context, userID, parentID string) (bool, error) {
	u, err := r.db.GetUserByID(ctx, userID)
	if err != nil {
		return false, nil
	}
	if u.ParentID == nil {
		return true, nil
	}
	next := u
	for next.ParentID != nil {
		if *next.ParentID == parentID {
			return true, nil
		}
		parent, err := r.db.GetUserByID(ctx, *next.ParentID)
		if err != nil {
			return false, nil
		}
		next = parent
	}
	return false, nil
}

// IsTeamMember verifies if the given userID belongs to teamID.
func (r *relationDB) IsTeamMember(ctx context.Context, teamID, userID string) (bool, error) {
	u, err := r.db.GetUserByID(ctx, userID)
	if err != nil {
		return false, nil
	}
	if u.TeamID == nil {
		return false, nil
	}
	return *u.TeamID == teamID, nil
}

// IsTeamOwner verifies if the given userID refers to an owner of the teamID.
func (r *relationDB) IsTeamOwner(ctx context.Context, teamID, userID string) (bool, error) {
	t, err := r.db.GetTeamByID(ctx, teamID)
	if err != nil {
		return false, nil
	}
	return t.OwnerID == userID, nil
}
