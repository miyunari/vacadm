package database

import "context"

type RelationDB interface {
	IsParentUser(ctx context.Context, userID, parentID string) (bool, error)
	IsTeamMember(ctx context.Context, teamID, userID string) (bool, error)
	IsTeamOwner(ctx context.Context, teamID, userID string) (bool, error)
}

func NewRelationDB(db Database) RelationDB {
	return &relationDB{
		db: db,
	}
}

type relationDB struct {
	db Database
}

func (r *relationDB) IsParentUser(ctx context.Context, userID, parentID string) (bool, error) {
	u, err := r.db.GetUserByID(ctx, userID)
	if err != nil {
		return false, err
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

func (r *relationDB) IsTeamMember(ctx context.Context, teamID, userID string) (bool, error) {
	u, err := r.db.GetUserByID(ctx, userID)
	if err != nil {
		return false, err
	}
	if u.TeamID == nil {
		return false, nil
	}
	return *u.TeamID == teamID, nil
}

func (r *relationDB) IsTeamOwner(ctx context.Context, teamID, userID string) (bool, error) {
	t, err := r.db.GetTeamByID(ctx, teamID)
	if err != nil {
		return false, err
	}
	return t.OwnerID == userID, nil
}
