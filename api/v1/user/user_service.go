package user

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/MninaTB/vacadm/api/v1/util"
	"github.com/MninaTB/vacadm/pkg/database"
	"github.com/MninaTB/vacadm/pkg/model"
)

// NewUserService returns a UserService.
func NewUserService(store database.Database, logger logrus.FieldLogger) *UserService {
	return &UserService{
		store:  store,
		logger: logger.WithField("component", "user-service"),
	}
}

// UserService implements http.HandlerFunc's to operate on user resources.
type UserService struct {
	store  database.Database
	logger logrus.FieldLogger
}

// Create reads the given payload and creates a store representation accordingly.
func (u *UserService) Create(w http.ResponseWriter, r *http.Request) {
	logger := u.logger.WithField("method", "create")
	logger.Info("create new user")
	var usr model.User
	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error(err)
		return
	}
	user, err := u.store.CreateUser(r.Context(), &usr)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	u.logger.Info("create user with id: ", user.ID)
	w.WriteHeader(http.StatusCreated)
}

// GetByID extracts a userID from URL and writes all user information into the
// given response writer.
func (u *UserService) GetByID(w http.ResponseWriter, r *http.Request) {
	logger := u.logger.WithField("method", "read")
	logger.Info("get user by id")
	userID, err := util.UserIDFromRequest(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	usr, err := u.store.GetUserByID(r.Context(), userID)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(usr)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	u.logger.Info("get user with id: ", userID)
}

// List retuns a list of all users available on the internal store.
func (u *UserService) List(w http.ResponseWriter, r *http.Request) {
	logger := u.logger.WithField("method", "list")
	logger.Info("retrieve user list")
	list, err := u.store.ListUsers(r.Context())
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(&list)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	u.logger.Info("get list of users")
}

// Update reads new user information from the request body and updates the store
// representation accordingly.
func (u *UserService) Update(w http.ResponseWriter, r *http.Request) {
	logger := u.logger.WithField("method", "update")
	logger.Info("update user")
	var usr model.User
	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := u.store.UpdateUser(r.Context(), &usr)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(&user)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	u.logger.Info("update user with id: ", usr.ID)
}

// Delete a user associated to the given userID in the URL.
func (u *UserService) Delete(w http.ResponseWriter, r *http.Request) {
	logger := u.logger.WithField("method", "delete")
	logger.Info("delete user")
	userID, err := util.UserIDFromRequest(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = u.store.DeleteUser(r.Context(), userID)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	u.logger.Info("delete user with id: ", userID)
	w.WriteHeader(http.StatusAccepted)
}
