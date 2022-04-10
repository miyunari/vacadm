package util

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	ErrDoesNotExistUserID   = errors.New("could not extract userID")
	ErrDoesNotExistTeamID   = errors.New("could not extract teamID")
	ErrDoesNotExistParentID = errors.New("could not extract parentID")
)

func TeamIDFromRequest(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	teamID, ok := vars["teamID"]
	if !ok {
		return "", ErrDoesNotExistTeamID
	}
	return teamID, nil
}

func UserIDFromRequest(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	usrID, ok := vars["userID"]
	if !ok {
		return "", ErrDoesNotExistUserID
	}
	return usrID, nil
}

func ParentIDFromRequest(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	usrID, ok := vars["parentID"]
	if !ok {
		return "", ErrDoesNotExistParentID
	}
	return usrID, nil
}
