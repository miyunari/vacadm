package util

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	// ErrDoesNotExistUserID is an error returned when a userID does not exist
	// in a URL.
	ErrDoesNotExistUserID = errors.New("could not extract userID")
	// ErrDoesNotExistTeamID is an error returned when a teamID does not exist
	// in a URL.
	ErrDoesNotExistTeamID = errors.New("could not extract teamID")
	// ErrDoesNotExistParentID is an error returned when a parentID does not exist
	// in a URL.
	ErrDoesNotExistParentID = errors.New("could not extract parentID")
)

// TeamIDFromRequest reads a teamID from the given request.
// if no teamID can be found, an error is returned.
func TeamIDFromRequest(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	teamID, ok := vars["teamID"]
	if !ok {
		return "", ErrDoesNotExistTeamID
	}
	return teamID, nil
}

// UserIDFromRequest reads a userID from the given request.
// if no userID can be found, an error is returned.
func UserIDFromRequest(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	usrID, ok := vars["userID"]
	if !ok {
		return "", ErrDoesNotExistUserID
	}
	return usrID, nil
}

// ParentIDFromRequest reads a parentID from the given request.
// if no parentID can be found, an error is returned.
func ParentIDFromRequest(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	usrID, ok := vars["parentID"]
	if !ok {
		return "", ErrDoesNotExistParentID
	}
	return usrID, nil
}
