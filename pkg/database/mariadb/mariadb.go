package mariadb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/MninaTB/vacadm/pkg/model"
)

const (
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

	/*
		userUpdate = `
		  	UPDATE user
			SET (
				id ?, parent_id ?,
				firstname ?, lastname ?,
				email ?
			)
			WHERE uuid = ?
		`
	*/
	userDelete = `
		DELETE FROM user
		WHERE id = ?
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

	teamDelete = `
		DELETE FROM team
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
		DELETE FROM vaccation
		WHERE id = ?
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

	vaccationRequestDelete = `
		DELETE FROM vaccation_request
		WHERE id = ?
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

	vaccationRessourceDelete = `
  		DELETE FROM vaccation_ressource
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

func (m *MariaDB) CreateUser(u *model.User) (*model.User, error) {
	fmt.Println(*u.TeamID)
	row, err := m.db.QueryContext(context.Background(), userCreate, u.ParentID, u.TeamID, u.Email, u.FirstName, u.LastName)
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

func (m *MariaDB) GetUserByID(uuid string) (*model.User, error) {
	row := m.db.QueryRowContext(context.Background(), userSelectByID, uuid)
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

func (m *MariaDB) ListUsers() ([]*model.User, error) {
	allusr := make([]*model.User, 0)
	rows, err := m.db.QueryContext(context.Background(), basicUserSelect)
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

/*
func (m *MariaDB) UpdateUser(u *model.User) (*model.User, error) {
	row := m.db.QueryRowContext(context.Background(), userUpdate)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	user := &model.User{}
	var parentID sql.NullString
	var createdAt, updatedAt sql.NullTime
	err = row.Scan(&user.ID, &parentID, &createdAt, &updatedAt, &user.FirstName, &user.LastName, &user.Email)
	if err != nil {
		return nil, err
	}
	if parentID.Valid {
		user.ParentID = &parentID.String
	}
	if createdAt.Valid {
		user.CreatedAt = &createdAt.Time
	}
	if updatedAt.Valid {
		user.UpdatedAt = &updatedAt.Time
	}
	return user, nil
}*/

func (m *MariaDB) DeleteUser(uuid string) error {
	row := m.db.QueryRowContext(context.Background(), userDelete, uuid)
	err := row.Err()
	if err != nil {
		return err
	}
	return err
}

/*
func (m *MariaDB) CreateTeam(t *model.Team) (*model.Team, error) {

}*/

func (m *MariaDB) GetTeamByID(uuid string) (*model.Team, error) {
	row := m.db.QueryRowContext(context.Background(), teamSelectByID, uuid)
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

func (m *MariaDB) ListTeams() ([]*model.Team, error) {
	allTeams := make([]*model.Team, 0)
	rows, err := m.db.QueryContext(context.Background(), basicTeamSelect)
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

func (m *MariaDB) ListTeamUsers(uuid string) ([]*model.User, error) {
	teamUser := make([]*model.User, 0)
	rows, err := m.db.QueryContext(context.Background(), teamUserSelectByID, uuid)
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

/*
func (m *MariaDB) UpdateTeam(t *model.Team) (*model.Team, error) {

}*/

func (m *MariaDB) DeleteTeam(uuid string) error {
	row := m.db.QueryRowContext(context.Background(), teamDelete, uuid)
	err := row.Err()
	if err != nil {
		return err
	}
	return err
}

func (m *MariaDB) GetVaccationByID(uuid string) (*model.Vaccation, error) {
	row := m.db.QueryRowContext(context.Background(), vaccationSelectByID, uuid)
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

func (m *MariaDB) ListVaccations() ([]*model.Vaccation, error) {
	allVaccations := make([]*model.Vaccation, 0)
	rows, err := m.db.QueryContext(context.Background(), basicVaccationSelect)
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

func (m *MariaDB) DeleteVaccation(uuid string) error {
	row := m.db.QueryRowContext(context.Background(), vaccationDelete, uuid)
	err := row.Err()
	if err != nil {
		return err
	}
	return err
}

/*
func (m *MariaDB) CreateVaccationRequest(*model.VaccationRequest) (*model.VaccationRequest, error) {

}*/

func (m *MariaDB) GetVaccationRequestByID(uuid string) (*model.VaccationRequest, error) {
	row := m.db.QueryRowContext(context.Background(), vaccationRequestSelectByID, uuid)
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

func (m *MariaDB) ListVaccationRequests() ([]*model.VaccationRequest, error) {
	allVaccationRequests := make([]*model.VaccationRequest, 0)
	rows, err := m.db.QueryContext(context.Background(), basicVaccationRequestSelect)
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

/*
func (m *MariaDB) UpdateVaccationRequest(v *model.VaccationRequest) (*model.VaccationRequest, error) {

}*/

func (m *MariaDB) DeleteVaccationRequest(uuid string) error {
	row := m.db.QueryRowContext(context.Background(), vaccationRequestDelete, uuid)
	err := row.Err()
	if err != nil {
		return err
	}
	return err
}

/*
func (m *MariaDB) CreateVaccationRessource(v *model.VaccationRessource) (*model.VaccationRessource, error) {

}*/

func (m *MariaDB) GetVaccationRessourceByID(uuid string) (*model.VaccationRessource, error) {
	row := m.db.QueryRowContext(context.Background(), vaccationRessourceSelectByID, uuid)
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

func (m *MariaDB) ListVaccationRessources() ([]*model.VaccationRessource, error) {
	allVaccationRessources := make([]*model.VaccationRessource, 0)
	rows, err := m.db.QueryContext(context.Background(), basicVaccationRessourceSelect)
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

/*
func (m *MariaDB) UpdateVaccationRessource(v *model.VaccationRessource) (*model.VaccationRessource, error) {

}*/

func (m *MariaDB) DeleteVaccationRessource(uuid string) error {
	row := m.db.QueryRowContext(context.Background(), vaccationRessourceDelete, uuid)
	err := row.Err()
	if err != nil {
		return err
	}
	return err
}
