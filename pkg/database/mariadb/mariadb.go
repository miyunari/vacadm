package mariadb

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/MninaTB/vacadm/pkg/model"
)

const (
	basicUserSelect = `
	  	SELECT
			id,
		  	parent_id,
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
		  	uuid,
			parent_uuid,
			fistname, 
			lastname,
		  	email,
		  	created_at, updated_at
		)
	  VALUES (
		 	:uuid, :parent_uuid,
		  	:firstname, :lastname,
		  	email,
			DEFAULT, DEFAULT
	  )
	`

	userUpdate = `
	  	UPDATE user 
		SET (
			uuid :uuid, parent_uuid :parent_uuid, 
			firstname :firstname, lastname :lastname, 
			email :email
		)
		WHERE uuid = ?
	`

	userDelete = `
		DELETE FROM user
		WHERE uuid = ?
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
	u.ID = uuid.NewString()
	_, err := m.db.ExecContext(context.Background(), userCreate, u.ID, u.FirstName, u.LastName, u.Email, u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return m.GetUserByID(u.ID)
}

func (m *MariaDB) GetUserByID(uuid string) (*model.User, error) {
	row := m.db.QueryRowContext(context.Background(), userSelectByID, uuid)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	u := &model.User{}
	var parentID sql.NullString
	var createdAt, updatedAt sql.NullTime
	err = row.Scan(&u.ID, &parentID, &createdAt, &updatedAt, &u.FirstName, &u.LastName, &u.Email)
	if err != nil {
		return nil, err
	}
	if parentID.Valid {
		u.ParentID = &parentID.String
	}
	if createdAt.Valid {
		u.CreatedAt = &createdAt.Time
	}
	if updatedAt.Valid {
		u.UpdatedAt = &updatedAt.Time
	}
	return u, nil
}

func (m *MariaDB) ListUsers() (*model.User, error) {
	row := m.db.QueryRowContext(context.Background(), userSelectByID)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	u := &model.User{}
	row.Scan(u.ID, u.ParentID, u.CreatedAt, u.UpdatedAt, u.FirstName, u.LastName, u.Email)
	return u, nil
}

func (m *MariaDB) UpdateUser(u *model.User) (*model.User, error) {
	row := m.db.QueryRowContext(context.Background(), userUpdate)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	user := &model.User{}
	row.Scan(user.ID, user.ParentID, user.CreatedAt, user.UpdatedAt, user.FirstName, user.LastName, user.Email)
	return user, nil
}

func (m *MariaDB) DeleteUser(uuid string) error {
	row := m.db.QueryRowContext(context.Background(), userDelete, uuid)
	err := row.Err()
	if err != nil {
		return err
	}
	return err
}
