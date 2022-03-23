package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/MninaTB/vacadm/pkg/database"
	"github.com/MninaTB/vacadm/pkg/model"
)

func NewUserService(store database.Database, logger logrus.FieldLogger) *userService {
	return &userService{
		store:  store,
		logger: logger.WithField("component", "user-service"),
	}
}

type userService struct {
	store  database.Database
	logger logrus.FieldLogger
}

func (u *userService) Create(w http.ResponseWriter, r *http.Request) {
	logger := u.logger.WithField("method", "create")
	logger.Info("create new user")
	var usr model.User
	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error(err)
		return
	}
	user, err := u.store.CreateUser(&usr)
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

func (u *userService) GetByID(w http.ResponseWriter, r *http.Request) {
	logger := u.logger.WithField("method", "read")
	logger.Info("get user by id")
	userID, err := extractUserID(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	usr, err := u.store.GetUserByID(userID)
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

func (u *userService) List(w http.ResponseWriter, r *http.Request) {
	logger := u.logger.WithField("method", "list")
	logger.Info("retrieve user list")
	list, err := u.store.ListUsers()
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(&list)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	u.logger.Info("get list of users")
}

func (u *userService) Update(w http.ResponseWriter, r *http.Request) {
	logger := u.logger.WithField("method", "update")
	logger.Info("update user")
	var usr model.User
	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := u.store.UpdateUser(&usr)
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

func (u *userService) Delete(w http.ResponseWriter, r *http.Request) {
	logger := u.logger.WithField("method", "delete")
	logger.Info("delete user")
	userID, err := extractUserID(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = u.store.DeleteUser(userID)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	u.logger.Info("delete user with id: ", userID)
	w.WriteHeader(http.StatusAccepted)
}

func extractUserID(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	usrID, ok := vars["userID"]
	if !ok {
		return "", errors.New("could not extract userID")
	}
	return usrID, nil
}
