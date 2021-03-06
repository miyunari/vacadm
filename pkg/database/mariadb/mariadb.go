package mariadb

import (
	"context"
	"database/sql"
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
			created_at
		)
		VALUES (
			UUID(), ?,
			?, ?,
			?, ?,
			NOW()
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
			id,
			owner_id, name,
			created_at
		)
		VALUES (
			UUID(),
			?, ?,
			NOW()
		)RETURNING id, created_at
	`

	basicTeamSelect = `
		SELECT
			id,
			owner_id, name,
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
			owner_id = ?,
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

	vacationCreate = `
	INSERT INTO vacation (
		id, user_id,
		from, to,
		created_at
	)
	VALUES (
		UUID(), ?
		?, ?
		NOW()
	) RETURNING id, created_at
`

	basicVacationSelect = `
		SELECT
			id,
			user_id,
			approved_id,
			from, to,
			created_at
		FROM vacation
	`

	teamVacationSelect = basicVacationSelect + `
		INNER JOIN user ON vacation.user_id=user.id
		WHERE user.team_id = ?
	`

	vacationSelectByID = basicVacationSelect + `
		WHERE id = ?
	`

	vacationDelete = `
		UPDATE vacation
		SET
			deleted_at = Now()
		WHERE id = ?
	`

	vacationRequestCreate = `
		INSERT INTO vaccation_request (
			id, user_id,
			from, to,
			created_at
		)
		VALUES (
			UUID(), ?
			?, ?
			NOW()
		) RETURNING id, created_at
	`

	basicVacationRequestSelect = `
		SELECT
			id,
			user_id,
			from, to,
			created_at, updated_at,
		FROM vacation_request
	`

	vacationRequestSelectByID = basicVacationRequestSelect + `
		WHERE id = ?
	`

	vacationRequestUpdate = `
		UPDATE vacation_request
		SET
			user_id = ?,
			from = ?, to = ?,
			updated_at = NOW()
		WHERE id = ?
	`

	vacationRequestDelete = `
		UPDATE vacation_request
		SET
			updated_at = NOW(),
			deleted_at = Now()
		WHERE id = ?
	`

	vacationResourceCreate = `
		INSERT INTO vacation_resource (
			id,
			user_id, yearly_days,
			created_at
		)
		VALUES (
			UUID(),
			?, ?
			NOW()
		) RETURNING id, created_at
	`

	basicVacationResourceSelect = `
		SELECT
			id,
			user_id,
			yearly_days,
			from, to,
			created_at, updated_at,
		FROM vacation_resource
	`

	vacationResourceSelectByID = basicVacationResourceSelect + `
		WHERE id = ?
	`

	vacationResourceUpdate = `
		UPDATE vacation_resource
		SET
			user_id = ?,
			yearly_days = ?,
			from = ?, to = ?,
			updated_at = NOW()
		WHERE id = ?
	`

	vacationResourceDelete = `
		UPDATE vacationResource
		SET
			updated_at = NOW(),
			deleted_at = Now()
		WHERE id = ?
	`
)

// NewMariaDB returns initialized MariaDB that fulfills
// the database interface.
func NewMariaDB(db *sql.DB) *MariaDB {
	return &MariaDB{
		db:     db,
		logger: logrus.New().WithField("component", "mariaDB"),
	}
}

// MariaDB implements the database interface.
type MariaDB struct {
	db     *sql.DB
	logger logrus.FieldLogger
}

// CreateUser stores an internal copy of the given user, if email address is
// not already in use, given parentID and/or teamID exists.
// Returns copy with assigned userID.
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

// GetUserByID returns the associated user by the given id.
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

// ListUsers returns a copy of the internal user list.
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

// UpdateUser updates user entry by the given user.
func (m *MariaDB) UpdateUser(ctx context.Context, u *model.User) (*model.User, error) {
	tx, err := m.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	_, err = tx.ExecContext(ctx, userUpdate, u.ParentID, u.TeamID, u.FirstName, u.LastName, u.Email, u.ID)
	if err != nil {
		if errTX := tx.Rollback(); err != nil {
			return nil, errTX
		}
		return nil, err
	}
	row := tx.QueryRowContext(ctx, "SELECT updated_at WHERE id = ?", u.ID)
	var updatedAt time.Time
	err = row.Scan(&updatedAt)
	if err != nil {
		if errTX := tx.Rollback(); err != nil {
			return nil, errTX
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

// DeleteUser removes user entry by the given id.
func (m *MariaDB) DeleteUser(ctx context.Context, uuid string) error {
	row := m.db.QueryRowContext(ctx, userDelete, uuid)
	err := row.Err()
	if err != nil {
		return err
	}
	return nil
}

// CreateTeam stores an internal copy of the given team.
// Returns copy with assigned teamID.
func (m *MariaDB) CreateTeam(ctx context.Context, t *model.Team) (*model.Team, error) {
	row, err := m.db.QueryContext(ctx, teamCreate, t.OwnerID, t.Name)
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

// GetTeamByID returns the associated team by the given id.
func (m *MariaDB) GetTeamByID(ctx context.Context, uuid string) (*model.Team, error) {
	row := m.db.QueryRowContext(ctx, teamSelectByID, uuid)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	t := &model.Team{}
	var createdAt, updatedAt sql.NullTime
	err = row.Scan(&t.ID, &t.OwnerID, &t.Name, &createdAt, &updatedAt)
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

// ListTeams returns a copy of the internal team list.
func (m *MariaDB) ListTeams(ctx context.Context) ([]*model.Team, error) {
	allTeams := make([]*model.Team, 0)
	rows, err := m.db.QueryContext(ctx, basicTeamSelect)
	if err != nil {
		return nil, err
	}
	var createdAt, updatedAt sql.NullTime
	for rows.Next() {
		t := model.Team{}
		err = rows.Scan(&t.ID, &t.OwnerID, t.Name, &createdAt, &updatedAt)
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

// ListTeamUsers returns a list of users associated by the given teamID
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

// UpdateTeam updates team entry by the given team.
func (m *MariaDB) UpdateTeam(ctx context.Context, t *model.Team) (*model.Team, error) {
	tx, err := m.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	_, err = tx.ExecContext(ctx, teamUpdate, t.OwnerID, t.Name, t.ID)
	if err != nil {
		if errTX := tx.Rollback(); err != nil {
			return nil, errTX
		}
		return nil, err
	}
	row := tx.QueryRowContext(ctx, "SELECT updated_at WHERE id = ?", t.ID)
	var updatedAt time.Time
	err = row.Scan(&updatedAt)
	if err != nil {
		if errTX := tx.Rollback(); err != nil {
			return nil, errTX
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

// DeleteTeam removes team entry by the given id.
func (m *MariaDB) DeleteTeam(ctx context.Context, uuid string) error {
	row := m.db.QueryRowContext(ctx, teamDelete, uuid)
	err := row.Err()
	if err != nil {
		return err
	}
	return err
}

// CreateVacation stores an internal copy of the given vacation resource.
// Returns copy with assigned vacationID.
func (m *MariaDB) CreateVacation(ctx context.Context, v *model.Vacation) (*model.Vacation, error) {
	row, err := m.db.QueryContext(ctx, vacationCreate, v.UserID, v.ApprovedBy, v.From, v.To)
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

// GetVacationByID returns the associated vacation by the given id.
func (m *MariaDB) GetVacationByID(ctx context.Context, uuid string) (*model.Vacation, error) {
	row := m.db.QueryRowContext(ctx, vacationSelectByID, uuid)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	v := &model.Vacation{}
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
		v.ApprovedBy = &approvedID.String
	}
	return v, nil
}

// GetVacationByTeamID returns the list of vacations of one team by given teamID.
func (m *MariaDB) GetVacationsByTeamID(ctx context.Context, tID string) ([]*model.Vacation, error) {
	teamVacations := make([]*model.Vacation, 0)
	rows, err := m.db.QueryContext(ctx, teamVacationSelect, tID)
	if err != nil {
		return nil, err
	}
	var createdAt, from, to sql.NullTime
	var userID sql.NullString
	var approvedID sql.NullString
	for rows.Next() {
		v := model.Vacation{}
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
			v.ApprovedBy = &approvedID.String
		}
		teamVacations = append(teamVacations, &v)
	}
	return teamVacations, nil
}

// ListVacations returns a copy of the internal vacation list.
func (m *MariaDB) ListVacations(ctx context.Context) ([]*model.Vacation, error) {
	allVacations := make([]*model.Vacation, 0)
	rows, err := m.db.QueryContext(ctx, basicVacationSelect)
	if err != nil {
		return nil, err
	}
	var createdAt, from, to sql.NullTime
	var userID sql.NullString
	var approvedID sql.NullString
	for rows.Next() {
		v := model.Vacation{}
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
			v.ApprovedBy = &approvedID.String
		}
		allVacations = append(allVacations, &v)
	}
	return allVacations, nil
}

// DeleteVacation removes vacation entry by the given id.
func (m *MariaDB) DeleteVacation(ctx context.Context, uuid string) error {
	row := m.db.QueryRowContext(ctx, vacationDelete, uuid)
	err := row.Err()
	if err != nil {
		return err
	}
	return err
}

// CreateVacationRequest stores an internal copy of the given vacationRequest.
// Returns copy with assigned vacationRequestID.
func (m *MariaDB) CreateVacationRequest(ctx context.Context, v *model.VacationRequest) (*model.VacationRequest, error) {
	row, err := m.db.QueryContext(ctx, vacationRequestCreate, v.UserID, v.From, v.To)
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

// GetVacationRequestByID returns the associated vacationRequest by the given id.
func (m *MariaDB) GetVacationRequestByID(ctx context.Context, uuid string) (*model.VacationRequest, error) {
	row := m.db.QueryRowContext(ctx, vacationRequestSelectByID, uuid)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	v := &model.VacationRequest{}
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

// ListVacationRequests returns a copy of the internal vacationRequest list.
func (m *MariaDB) ListVacationRequests(ctx context.Context) ([]*model.VacationRequest, error) {
	allVacationRequests := make([]*model.VacationRequest, 0)
	rows, err := m.db.QueryContext(ctx, basicVacationRequestSelect)
	if err != nil {
		return nil, err
	}
	var userID sql.NullString
	var createdAt, updatedAt, from, to sql.NullTime
	for rows.Next() {
		v := model.VacationRequest{}
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
		allVacationRequests = append(allVacationRequests, &v)
	}
	return allVacationRequests, nil
}

// UpdateVacationRequest updates vacationRequest entry by the given vacationRequest.
func (m *MariaDB) UpdateVacationRequest(ctx context.Context, v *model.VacationRequest) (*model.VacationRequest, error) {
	tx, err := m.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	_, err = tx.ExecContext(ctx, vacationRequestUpdate, v.UserID, v.From, v.To, v.ID)
	if err != nil {
		if errTX := tx.Rollback(); err != nil {
			return nil, errTX
		}
		return nil, err
	}
	row := tx.QueryRowContext(ctx, "SELECT updated_at WHERE id = ?", v.ID)
	var updatedAt time.Time
	err = row.Scan(&updatedAt)
	if err != nil {
		if errTX := tx.Rollback(); err != nil {
			return nil, errTX
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

// DeleteVacationRequest removes vacationRequest entry by the given id.
func (m *MariaDB) DeleteVacationRequest(ctx context.Context, uuid string) error {
	row := m.db.QueryRowContext(ctx, vacationRequestDelete, uuid)
	err := row.Err()
	if err != nil {
		return err
	}
	return err
}

// CreateVacationResource stores an internal copy of the given vacationResource.
// Returns copy with assigned vacationResourceID.
func (m *MariaDB) CreateVacationResource(ctx context.Context, v *model.VacationResource) (*model.VacationResource, error) {
	row, err := m.db.QueryContext(ctx, vacationResourceCreate, v.UserID, v.YearlyDays)
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

// GetVacationResourceByID returns the associated vacationResource by the given id.
func (m *MariaDB) GetVacationResourceByID(ctx context.Context, uuid string) (*model.VacationResource, error) {
	row := m.db.QueryRowContext(ctx, vacationResourceSelectByID, uuid)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	v := &model.VacationResource{}
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

// ListVacationResource returns a copy of the internal vacationResource list.
func (m *MariaDB) ListVacationResource(ctx context.Context) ([]*model.VacationResource, error) {
	allVacationResources := make([]*model.VacationResource, 0)
	rows, err := m.db.QueryContext(ctx, basicVacationResourceSelect)
	if err != nil {
		return nil, err
	}
	var userID sql.NullString
	var createdAt, updatedAt, from, to sql.NullTime
	for rows.Next() {
		v := model.VacationResource{}
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
		allVacationResources = append(allVacationResources, &v)
	}
	return allVacationResources, nil
}

// UpdateVacationResource updates vacationResource entry by the given vacationResource.
func (m *MariaDB) UpdateVacationResource(ctx context.Context, v *model.VacationResource) (*model.VacationResource, error) {
	tx, err := m.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	_, err = tx.ExecContext(ctx, vacationResourceUpdate, v.UserID, v.YearlyDays, v.From, v.To, v.ID)
	if err != nil {
		if errTX := tx.Rollback(); err != nil {
			return nil, errTX
		}
		return nil, err
	}
	row := tx.QueryRowContext(ctx, "SELECT updated_at WHERE id = ?", v.ID)
	var updatedAt time.Time
	err = row.Scan(&updatedAt)
	if err != nil {
		if errTX := tx.Rollback(); err != nil {
			return nil, errTX
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

// DeleteVacationResource removes vacationResource entry by the given id.
func (m *MariaDB) DeleteVacationResource(ctx context.Context, uuid string) error {
	row := m.db.QueryRowContext(ctx, vacationResourceDelete, uuid)
	err := row.Err()
	if err != nil {
		return err
	}
	return err
}
