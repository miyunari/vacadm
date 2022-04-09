package vacationrequest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/MninaTB/vacadm/api/v1/util"
	"github.com/MninaTB/vacadm/pkg/database"
	"github.com/MninaTB/vacadm/pkg/model"
	"github.com/MninaTB/vacadm/pkg/notify"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func NewVacationRequest(store database.Database, notifier notify.Notifier, logger logrus.FieldLogger) *vacationRequest {
	return &vacationRequest{
		store:    store,
		notifier: notifier,
		logger:   logger.WithField("component", "vacation-request-service"),
	}
}

type vacationRequest struct {
	store    database.Database
	notifier notify.Notifier
	logger   logrus.FieldLogger
}

func (v *vacationRequest) Create(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("method", "create")
	logger.Info("create new vacation-request")
	var vr model.VacationRequest
	err := json.NewDecoder(r.Body).Decode(&vr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error(err)
		return
	}
	newVR, err := v.store.CreateVacationRequest(r.Context(), &vr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err)
		return
	}
	userID, err := util.UserIDFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error(err)
		return
	}
	user, err := v.store.GetUserByID(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err)
		return
	}
	if user.ParentID != nil {
		action := fmt.Sprintf("new vacation request from %s %s, id: %s", user.FirstName, user.LastName, user.ID)
		err = v.notifier.NotifyUser(r.Context(), *user.ParentID, action)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(err)
			return
		}
	}
	err = json.NewEncoder(w).Encode(newVR)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err)
		return
	}
	v.logger.Info("create vacation-request with ID: ", newVR.ID)
	w.WriteHeader(http.StatusCreated)
}

func (v *vacationRequest) GetByID(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("component", "read")
	logger.Info("get vacation-request by id")
	vrID, err := extractVacationRequestID(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	vR, err := v.store.GetVacationRequestByID(r.Context(), vrID)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(&vR)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	v.logger.Info("get vacation-request with id: ", vR)
}

func (v *vacationRequest) List(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("method", "list")
	logger.Info("retrieve vacation-request list")
	list, err := v.store.ListVacationRequests(r.Context())
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
	v.logger.Info("get list of vacation-requests")
}

func (v *vacationRequest) Update(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("method", "update")
	logger.Info("update vacation-request")
	var vr model.VacationRequest
	err := json.NewDecoder(r.Body).Decode(&vr)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newVR, err := v.store.UpdateVacationRequest(r.Context(), &vr)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID, err := util.UserIDFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error(err)
		return
	}
	user, err := v.store.GetUserByID(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err)
		return
	}
	if user.ParentID != nil {
		action := fmt.Sprintf("updated vacation request from %s %s, id: %s", user.FirstName, user.LastName, user.ID)
		err = v.notifier.NotifyUser(r.Context(), *user.ParentID, action)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(err)
			return
		}
	}
	err = json.NewEncoder(w).Encode(&newVR)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	v.logger.Info("update vacation-request with id: ", newVR.ID)

}

func (v *vacationRequest) Delete(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("method", "delete")
	logger.Info("delete vacation-request")
	vrID, err := extractVacationRequestID(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = v.store.DeleteVacationRequest(r.Context(), vrID)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	v.logger.Info("delete vacation-request with id: ", vrID)
	w.WriteHeader(http.StatusAccepted)
}

func extractVacationRequestID(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	vacationRequestID, ok := vars["vacationRequestID"]
	if !ok {
		return "", errors.New("could not extract vacationRequestID")
	}
	return vacationRequestID, nil
}
