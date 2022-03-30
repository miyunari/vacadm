package vaccation

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MninaTB/vacadm/pkg/database"
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

func (v *vaccation) GetByID(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("component", "read")
	logger.Info("get vaccation by id")
	vacID, err := extractVaccationID(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	vaccation, err := v.store.GetVaccationByID(r.Context(), vacID)
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
	logger.Info("get vaccation list")
	list, err := v.store.ListVaccations(r.Context())
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

func (v *vaccation) Delete(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("component", "delete")
	logger.Info("delete vaccation")
	vacID, err := extractVaccationID(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = v.store.DeleteVaccation(r.Context(), vacID)
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
