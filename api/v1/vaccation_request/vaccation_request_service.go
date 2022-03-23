package vaccationrequest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MninaTB/vacadm/pkg/database"
	"github.com/MninaTB/vacadm/pkg/model"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func NewVaccationRequest(store database.Database, logger logrus.FieldLogger) *vaccationRequest {
	return &vaccationRequest{
		store:  store,
		logger: logger.WithField("component", "vaccation-request-service"),
	}
}

type vaccationRequest struct {
	store  database.Database
	logger logrus.FieldLogger
}

func (v *vaccationRequest) Create(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("method", "create")
	logger.Info("create new vaccation-request")
	var vr model.VaccationRequest
	err := json.NewDecoder(r.Body).Decode(&vr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error(err)
		return
	}
	newVR, err := v.store.CreateVaccationRequest(&vr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err)
		return
	}
	err = json.NewEncoder(w).Encode(newVR)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err)
		return
	}
	v.logger.Info("create vaccation-request with ID: ", newVR.ID)
	w.WriteHeader(http.StatusCreated)
}

func (v *vaccationRequest) GetByID(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("component", "read")
	logger.Info("get vaccation-request by id")
	vrID, err := extractVaccationRequestID(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	vR, err := v.store.GetVaccationRequestByID(vrID)
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
	v.logger.Info("get vaccation-request with id: ", vR)
}

func (v *vaccationRequest) List(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("method", "list")
	logger.Info("retrieve vaccation-request list")
	list, err := v.store.ListVaccationRequests()
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
	v.logger.Info("get list of vaccation-requests")
}

func (v *vaccationRequest) Update(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("method", "update")
	logger.Info("update vaccation-request")
	var vr model.VaccationRequest
	err := json.NewDecoder(r.Body).Decode(&vr)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newVR, err := v.store.UpdateVaccationRequest(&vr)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(&newVR)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	v.logger.Info("update vaccation-request with id: ", newVR.ID)

}

func (v *vaccationRequest) Delete(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("method", "delete")
	logger.Info("delete vaccation-request")
	vrID, err := extractVaccationRequestID(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = v.store.DeleteVaccationRequest(vrID)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	v.logger.Info("delete vaccation-request with id: ", vrID)
	w.WriteHeader(http.StatusAccepted)
}

func extractVaccationRequestID(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	vaccationRequestID, ok := vars["vaccationRequestID"]
	if !ok {
		return "", errors.New("could not extract vaccationRequestID")
	}
	return vaccationRequestID, nil
}
