package vaccation

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MninaTB/vacadm/pkg/database"
	"github.com/MninaTB/vacadm/pkg/model"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func NewVaccation(store database.Database, logger logrus.FieldLogger) *vaccation {
	return &vaccation{
		store:  store,
		logger: logger.WithField("component", "vaccation-service"),
	}
}

type vaccation struct {
	store  database.Database
	logger logrus.FieldLogger
}

func (v *vaccation) Create(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("component", "create")
	var vac model.Vaccation
	err := json.NewDecoder(r.Body).Decode(&vac)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	vaccation, err := v.store.CreateVaccation(&vac)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(&vaccation)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	v.logger.Info("created vaccation with id: ", vaccation.ID)
	w.WriteHeader(http.StatusCreated)
}

func (v *vaccation) GetByID(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("component", "read")
	vacID, err := extractVaccationID(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	vaccation, err := v.store.GetVaccationByID(vacID)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(&vaccation)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	v.logger.Info("get vaccation with id: ", vacID)
}

func (v *vaccation) List(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("component", "list")
	list, err := v.store.ListVaccations()
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(&list)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	v.logger.Info("get list of vaccations")
}

func (v *vaccation) Update(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	v.logger.Warn("update is not implemented")
}

func (v *vaccation) Delete(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("component", "delete")
	vacID, err := extractVaccationID(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = v.store.DeleteVaccation(vacID)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	v.logger.Info("delete vaccation with id: ", vacID)
	w.WriteHeader(http.StatusAccepted)
}

func extractVaccationID(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	vaccationID, ok := vars["vaccationID"]
	if !ok {
		return "", errors.New("could not extract vaccationID")
	}
	return vaccationID, nil
}
