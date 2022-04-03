package util

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	ErrDoesNotExistUserID = errors.New("could not extract userID")
	ErrDoesNotExistTeamID = errors.New("could not extract teamID")
)

func TeamIDFromRequest(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	teamID, ok := vars["teamID"]
	if !ok {
		return "", errors.New("could not extract teamID")
	}
	return teamID, nil
}

func UserIDFromRequest(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	usrID, ok := vars["userID"]
	if !ok {
		return "", errors.New("could not extract userID")
	}
	return usrID, nil
}
