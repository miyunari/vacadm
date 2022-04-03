package mariadb

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/MninaTB/vacadm/pkg/model"
)

const (
	userCreate = `
		INSERT INTO user (
			id, parent_id,
			team_id, email,
			firstname, lastname,
			created_at, updated_at
		)
  		VALUES (
			UUID(), ?,
			?, ?,
			?, ?, 
			NOW(), NOW()
  		) RETURNING id, created_at
	`

	basicUserSelect = `
	  	SELECT
			id,
			parent_id,
		  	team_id,
		  	created_at, updated_at,
		  	firstname, lastname,
		  	email
	  	FROM user
	`

	userSelectByID = basicUserSelect + `
	  	WHERE id = ?
	`

	userUpdate = `
	  	UPDATE user
		SET 
			parent_id = ?, team_id = ?,
			firstname = ?, lastname = ?,
			email = ?, updated_at = NOW()
		WHERE id = ?
	`

	userDelete = `
		UPDATE user
		SET 
			updated_at = NOW(), 
			deleted_at = Now()
		WHERE id = ?
	`

	teamCreate = `
		INSERT INTO team (
			id, name,
			created_at, updated_at
		)
	  	VALUES (
			UUID(), ?,
			NOW(), NOW()
	  	)RETURNING id, created_at
	`

	basicTeamSelect = `
	  	SELECT
			id,
			name,
		  	created_at, updated_at,
	  	FROM team
	`

	teamSelectByID = basicTeamSelect + `
	  	WHERE id = ?
	`

	teamUserSelectByID = basicUserSelect + `
		WHERE team_id = ?
	`

	teamUpdate = `
		UPDATE team
  		SET 
	  		name = ?,
	  		updated_at = NOW()
  		WHERE id = ?
	`

	teamDelete = `
		UPDATE team
		SET 
			updated_at = NOW(), 
			deleted_at = Now()
		WHERE id = ?
	`

	basicVaccationSelect = `
		SELECT
  			id,
  			user_id,
			approved_id,
			from, to,
			created_at
		FROM vaccation
	`

	vaccationSelectByID = basicVaccationSelect + `
		WHERE id = ?
	`

	vaccationDelete = `
		UPDATE vaccation
		SET 
			deleted_at = Now()
		WHERE id = ?
	`

	vaccationRequestCreate = `
		INSERT INTO user (
			id, user_id,
			from, to,
		  	created_at, updated_at
		)
  		VALUES (
			UUID(), ?
			NOW(), NOW()
			NOW(), NOW()
  		) RETURNING id, created_at
	`

	basicVaccationRequestSelect = `
	  	SELECT
			id,
			user_id,
			from, to,
		  	created_at, updated_at,
	  	FROM vaccation_request
	`

	vaccationRequestSelectByID = basicVaccationRequestSelect + `
	  	WHERE id = ?
	`

	vaccationRequestUpdate = `
		UPDATE vaccation_request
  		SET 
			user_id = ?,
	  		from = ?, to = ?,
	  		updated_at = NOW()
  		WHERE id = ?
	`

	vaccationRequestDelete = `
		UPDATE vaccation_request
		SET 
			updated_at = NOW(), 
			deleted_at = Now()
		WHERE id = ?
	`

	vaccationRessourceCreate = `
		INSERT INTO user (
			id, 
			user_id, yearly_days,
			created_at, updated_at
		)
	  	VALUES (
			UUID(),
			?, ?
			NOW(), NOW()
  		) RETURNING id, created_at
	`

	basicVaccationRessourceSelect = `
		SELECT
	  		id,
	  		user_id,
			yearly_days,
			from, to,
			created_at, updated_at,
		FROM vaccation_ressource
	`

	vaccationRessourceSelectByID = basicVaccationRessourceSelect + `
		WHERE id = ?
	`

	vaccationRessourceUpdate = `
		UPDATE vaccation_ressource
	  	SET 
			user_id = ?,
		  	yearly_days = ?,
			from = ?, to = ?,
		  	updated_at = NOW()
	  	WHERE id = ?
	`

	vaccationRessourceDelete = `
		UPDATE vaccationRessource
		SET 
			updated_at = NOW(), 
			deleted_at = Now()
		WHERE id = ?
	`
)

func NewMariaDB(db *sql.DB) *MariaDB {
	return &MariaDB{
		db:     db,
		logger: logrus.New().WithField("component", "mariaDB"),
	}
}

type MariaDB struct {
	db     *sql.DB
	logger logrus.FieldLogger
}

func (m *MariaDB) CreateUser(ctx context.Context, u *model.User) (*model.User, error) {
	row, err := m.db.QueryContext(ctx, userCreate, u.ParentID, u.TeamID, u.Email, u.FirstName, u.LastName)
	if err != nil {
		return nil, err
	}
	err = row.Err()
	if err != nil {
		return nil, err
	}
	var id string
	var createdAt time.Time
	row.Next()
	err = row.Scan(&id, &createdAt)
	if err != nil {
		return nil, err
	}
	u.ID = id
	u.CreatedAt = &createdAt
	return u, nil
}

func (m *MariaDB) GetUserByID(ctx context.Context, uuid string) (*model.User, error) {
	row := m.db.QueryRowContext(ctx, userSelectByID, uuid)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	u := &model.User{}
	var parentID, teamID sql.NullString
	var createdAt, updatedAt sql.NullTime
	err = row.Scan(&u.ID, &parentID, &teamID, &createdAt, &updatedAt, &u.FirstName, &u.LastName, &u.Email)
	if err != nil {
		return nil, err
	}
	if parentID.Valid {
		u.ParentID = &parentID.String
	}
	if teamID.Valid {
		u.TeamID = &teamID.String
	}
	if createdAt.Valid {
		u.CreatedAt = &createdAt.Time
	}
	if updatedAt.Valid {
		u.UpdatedAt = &updatedAt.Time
	}
	return u, nil
}

func (m *MariaDB) ListUsers(ctx context.Context) ([]*model.User, error) {
	allusr := make([]*model.User, 0)
	rows, err := m.db.QueryContext(ctx, basicUserSelect)
	if err != nil {
		return nil, err
	}
	var parentID, teamID sql.NullString
	var createdAt, updatedAt sql.NullTime
	for rows.Next() {
		u := model.User{}
		err = rows.Scan(&u.ID, &parentID, &teamID, &createdAt, &updatedAt, &u.FirstName, &u.LastName, &u.Email)
		if err != nil {
			return nil, err
		}
		if parentID.Valid {
			u.ParentID = &parentID.String
		}
		if teamID.Valid {
			u.TeamID = &teamID.String
		}
		if createdAt.Valid {
			u.CreatedAt = &createdAt.Time
		}
		if updatedAt.Valid {
			u.UpdatedAt = &updatedAt.Time
		}
		allusr = append(allusr, &u)
	}
	return allusr, nil
}

func (m *MariaDB) UpdateUser(ctx context.Context, u *model.User) (*model.User, error) {
	tx, err := m.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	_, err = tx.ExecContext(ctx, userUpdate, u.ParentID, u.TeamID, u.FirstName, u.LastName, u.Email, u.ID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}
	row := tx.QueryRowContext(ctx, "SELECT updated_at WHERE id = ?", u.ID)
	var updatedAt time.Time
	err = row.Scan(&updatedAt)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	u.UpdatedAt = &updatedAt
	return u, nil
}

func (m *MariaDB) DeleteUser(ctx context.Context, uuid string) error {
	row := m.db.QueryRowContext(ctx, userDelete, uuid)
	err := row.Err()
	if err != nil {
		return err
	}
	return nil
}

func (m *MariaDB) CreateTeam(ctx context.Context, t *model.Team) (*model.Team, error) {
	row, err := m.db.QueryContext(ctx, teamCreate, t.Name)
	if err != nil {
		return nil, err
	}
	err = row.Err()
	if err != nil {
		return nil, err
	}
	var id string
	var createdAt time.Time
	row.Next()
	err = row.Scan(&id, &createdAt)
	if err != nil {
		return nil, err
	}
	t.ID = id
	t.CreatedAt = &createdAt
	return t, nil
}

func (m *MariaDB) GetTeamByID(ctx context.Context, uuid string) (*model.Team, error) {
	row := m.db.QueryRowContext(ctx, teamSelectByID, uuid)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	t := &model.Team{}
	var createdAt, updatedAt sql.NullTime
	err = row.Scan(&t.ID, &createdAt, &updatedAt, &t.Name)
	if err != nil {
		return nil, err
	}
	if createdAt.Valid {
		t.CreatedAt = &createdAt.Time
	}
	if updatedAt.Valid {
		t.UpdatedAt = &updatedAt.Time
	}
	return t, nil
}

func (m *MariaDB) ListTeams(ctx context.Context) ([]*model.Team, error) {
	allTeams := make([]*model.Team, 0)
	rows, err := m.db.QueryContext(ctx, basicTeamSelect)
	if err != nil {
		return nil, err
	}
	var createdAt, updatedAt sql.NullTime
	for rows.Next() {
		t := model.Team{}
		err = rows.Scan(&t.ID, t.Name, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}
		if createdAt.Valid {
			t.CreatedAt = &createdAt.Time
		}
		if updatedAt.Valid {
			t.UpdatedAt = &updatedAt.Time
		}
		allTeams = append(allTeams, &t)
	}
	return allTeams, nil
}

func (m *MariaDB) ListTeamUsers(ctx context.Context, uuid string) ([]*model.User, error) {
	teamUser := make([]*model.User, 0)
	rows, err := m.db.QueryContext(ctx, teamUserSelectByID, uuid)
	if err != nil {
		return nil, err
	}
	var parentID, teamID sql.NullString
	var createdAt, updatedAt sql.NullTime
	for rows.Next() {
		u := model.User{}
		err = rows.Scan(&u.ID, &parentID, &teamID, &createdAt, &updatedAt, &u.FirstName, &u.LastName, &u.Email)
		if err != nil {
			return nil, err
		}
		if parentID.Valid {
			u.ParentID = &parentID.String
		}
		if teamID.Valid {
			u.TeamID = &teamID.String
		}
		if createdAt.Valid {
			u.CreatedAt = &createdAt.Time
		}
		if updatedAt.Valid {
			u.UpdatedAt = &updatedAt.Time
		}
		teamUser = append(teamUser, &u)
	}
	return teamUser, nil
}

func (m *MariaDB) UpdateTeam(ctx context.Context, t *model.Team) (*model.Team, error) {
	tx, err := m.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	_, err = tx.ExecContext(ctx, teamUpdate, t.Name, t.ID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}
	row := tx.QueryRowContext(ctx, "SELECT updated_at WHERE id = ?", t.ID)
	var updatedAt time.Time
	err = row.Scan(&updatedAt)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	t.UpdatedAt = &updatedAt
	return t, nil
}

func (m *MariaDB) DeleteTeam(ctx context.Context, uuid string) error {
	row := m.db.QueryRowContext(ctx, teamDelete, uuid)
	err := row.Err()
	if err != nil {
		return err
	}
	return err
}

func (m *MariaDB) GetVaccationByID(ctx context.Context, uuid string) (*model.Vaccation, error) {
	row := m.db.QueryRowContext(ctx, vaccationSelectByID, uuid)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	v := &model.Vaccation{}
	var createdAt, from, to sql.NullTime
	var userID sql.NullString
	var approvedID sql.NullString
	err = row.Scan(&v.ID, &createdAt, &from, &to, &userID, &approvedID)
	if err != nil {
		return nil, err
	}
	if createdAt.Valid {
		v.CreatedAt = &createdAt.Time
	}
	if from.Valid {
		v.From = from.Time
	}
	if to.Valid {
		v.To = to.Time
	}
	if userID.Valid {
		v.UserID = userID.String
	}
	if approvedID.Valid {
		v.ApprovedBy.ID = approvedID.String
	}
	return v, nil
}

func (m *MariaDB) ListVaccations(ctx context.Context) ([]*model.Vaccation, error) {
	allVaccations := make([]*model.Vaccation, 0)
	rows, err := m.db.QueryContext(ctx, basicVaccationSelect)
	if err != nil {
		return nil, err
	}
	var createdAt, from, to sql.NullTime
	var userID sql.NullString
	var approvedID sql.NullString
	for rows.Next() {
		v := model.Vaccation{}
		err = rows.Scan(&v.ID, &createdAt, &from, &to, &userID, &approvedID)
		if err != nil {
			return nil, err
		}
		if createdAt.Valid {
			v.CreatedAt = &createdAt.Time
		}
		if from.Valid {
			v.From = from.Time
		}
		if to.Valid {
			v.To = from.Time
		}
		if userID.Valid {
			v.UserID = userID.String
		}
		if approvedID.Valid {
			v.ApprovedBy.ID = approvedID.String
		}
		allVaccations = append(allVaccations, &v)
	}
	return allVaccations, nil
}

func (m *MariaDB) DeleteVaccation(ctx context.Context, uuid string) error {
	row := m.db.QueryRowContext(ctx, vaccationDelete, uuid)
	err := row.Err()
	if err != nil {
		return err
	}
	return err
}

func (m *MariaDB) CreateVaccationRequest(ctx context.Context, v *model.VaccationRequest) (*model.VaccationRequest, error) {
	row, err := m.db.QueryContext(ctx, vaccationRequestCreate, v.UserID, v.From, v.To)
	if err != nil {
		return nil, err
	}
	err = row.Err()
	if err != nil {
		return nil, err
	}
	var id string
	var createdAt, from, to time.Time
	row.Next()
	err = row.Scan(&id, &createdAt, &from, &to)
	if err != nil {
		return nil, err
	}
	v.ID = id
	v.CreatedAt = &createdAt
	v.From = from
	v.To = to
	return v, nil
}

func (m *MariaDB) GetVaccationRequestByID(ctx context.Context, uuid string) (*model.VaccationRequest, error) {
	row := m.db.QueryRowContext(ctx, vaccationRequestSelectByID, uuid)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	v := &model.VaccationRequest{}
	var userID sql.NullString
	var createdAt, updatedAt, from, to sql.NullTime
	err = row.Scan(&v.ID, &userID, &createdAt, &updatedAt, &from, &to)
	if err != nil {
		return nil, err
	}
	if userID.Valid {
		v.UserID = userID.String
	}
	if createdAt.Valid {
		v.CreatedAt = &createdAt.Time
	}
	if updatedAt.Valid {
		v.UpdatedAt = &updatedAt.Time
	}
	if from.Valid {
		v.From = from.Time
	}
	if to.Valid {
		v.To = to.Time
	}
	return v, nil
}

func (m *MariaDB) ListVaccationRequests(ctx context.Context) ([]*model.VaccationRequest, error) {
	allVaccationRequests := make([]*model.VaccationRequest, 0)
	rows, err := m.db.QueryContext(ctx, basicVaccationRequestSelect)
	if err != nil {
		return nil, err
	}
	var userID sql.NullString
	var createdAt, updatedAt, from, to sql.NullTime
	for rows.Next() {
		v := model.VaccationRequest{}
		err = rows.Scan(&v.ID, &createdAt, &updatedAt, &from, &to, &userID)
		if err != nil {
			return nil, err
		}
		if createdAt.Valid {
			v.CreatedAt = &createdAt.Time
		}
		if from.Valid {
			v.From = from.Time
		}
		if to.Valid {
			v.To = from.Time
		}
		if userID.Valid {
			v.UserID = userID.String
		}
		allVaccationRequests = append(allVaccationRequests, &v)
	}
	return allVaccationRequests, nil
}

func (m *MariaDB) UpdateVaccationRequest(ctx context.Context, v *model.VaccationRequest) (*model.VaccationRequest, error) {
	tx, err := m.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	_, err = tx.ExecContext(ctx, vaccationRequestUpdate, v.UserID, v.From, v.To, v.ID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}
	row := tx.QueryRowContext(ctx, "SELECT updated_at WHERE id = ?", v.ID)
	var updatedAt time.Time
	err = row.Scan(&updatedAt)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	v.UpdatedAt = &updatedAt
	return v, nil
}

func (m *MariaDB) DeleteVaccationRequest(ctx context.Context, uuid string) error {
	row := m.db.QueryRowContext(ctx, vaccationRequestDelete, uuid)
	err := row.Err()
	if err != nil {
		return err
	}
	return err
}

func (m *MariaDB) CreateVaccationRessource(ctx context.Context, v *model.VaccationRessource) (*model.VaccationRessource, error) {
	row, err := m.db.QueryContext(ctx, vaccationRessourceCreate, v.UserID, v.YearlyDays)
	if err != nil {
		return nil, err
	}
	err = row.Err()
	if err != nil {
		return nil, err
	}
	var id string
	var createdAt time.Time
	row.Next()
	err = row.Scan(&id, &createdAt)
	if err != nil {
		return nil, err
	}
	v.ID = id
	v.CreatedAt = &createdAt
	return v, nil
}

func (m *MariaDB) GetVaccationRessourceByID(ctx context.Context, uuid string) (*model.VaccationRessource, error) {
	row := m.db.QueryRowContext(ctx, vaccationRessourceSelectByID, uuid)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	v := &model.VaccationRessource{}
	var userID sql.NullString
	var createdAt, updatedAt, from, to sql.NullTime
	err = row.Scan(&v.ID, &userID, &createdAt, &updatedAt, &from, &to, v.YearlyDays)
	if err != nil {
		return nil, err
	}
	if userID.Valid {
		v.UserID = userID.String
	}
	if createdAt.Valid {
		v.CreatedAt = &createdAt.Time
	}
	if updatedAt.Valid {
		v.UpdatedAt = &updatedAt.Time
	}
	if from.Valid {
		v.From = from.Time
	}
	if to.Valid {
		v.To = to.Time
	}
	return v, nil
}

func (m *MariaDB) ListVaccationRessource(ctx context.Context) ([]*model.VaccationRessource, error) {
	allVaccationRessources := make([]*model.VaccationRessource, 0)
	rows, err := m.db.QueryContext(ctx, basicVaccationRessourceSelect)
	if err != nil {
		return nil, err
	}
	var userID sql.NullString
	var createdAt, updatedAt, from, to sql.NullTime
	for rows.Next() {
		v := model.VaccationRessource{}
		err = rows.Scan(&v.ID, &createdAt, &updatedAt, &from, &to, &userID)
		if err != nil {
			return nil, err
		}
		if createdAt.Valid {
			v.CreatedAt = &createdAt.Time
		}
		if from.Valid {
			v.From = from.Time
		}
		if to.Valid {
			v.To = from.Time
		}
		if userID.Valid {
			v.UserID = userID.String
		}
		allVaccationRessources = append(allVaccationRessources, &v)
	}
	return allVaccationRessources, nil
}

func (m *MariaDB) UpdateVaccationRessource(ctx context.Context, v *model.VaccationRessource) (*model.VaccationRessource, error) {
	tx, err := m.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	_, err = tx.ExecContext(ctx, vaccationRessourceUpdate, v.UserID, v.YearlyDays, v.From, v.To, v.ID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}
	row := tx.QueryRowContext(ctx, "SELECT updated_at WHERE id = ?", v.ID)
	var updatedAt time.Time
	err = row.Scan(&updatedAt)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	v.UpdatedAt = &updatedAt
	return v, nil
}

func (m *MariaDB) DeleteVaccationRessource(ctx context.Context, uuid string) error {
	row := m.db.QueryRowContext(ctx, vaccationRessourceDelete, uuid)
	err := row.Err()
	if err != nil {
		return err
	}
	return err
}

func (m *MariaDB) IsParentUser(ctx context.Context, userID, parentID string) (bool, error) {
	return false, errors.New("not implemented")
}

func (m *MariaDB) IsTeamMember(ctx context.Context, teamID, userID string) (bool, error) {
	return false, errors.New("not implemented")
}

func (m *MariaDB) IsTeamOwner(ctx context.Context, teamID, userID string) (bool, error) {
	return false, errors.New("not implemented")
}
