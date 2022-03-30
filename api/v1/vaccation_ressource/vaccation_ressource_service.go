package vaccationressources

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MninaTB/vacadm/pkg/database"
	"github.com/MninaTB/vacadm/pkg/model"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func NewVaccationRessource(store database.Database, logger logrus.FieldLogger) *vaccationRessource {
	return &vaccationRessource{
		store:  store,
		logger: logger.WithField("component", "vaccation-ressource"),
	}
}

type vaccationRessource struct {
	store  database.Database
	logger logrus.FieldLogger
}

func (v *vaccationRessource) Create(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("method", "create")
	logger.Info("create new vaccation-ressource")
	var vr model.VaccationRessource
	err := json.NewDecoder(r.Body).Decode(&vr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error(err)
		return
	}
	newVR, err := v.store.CreateVaccationRessource(r.Context(), &vr)
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
	v.logger.Info("create vaccation-ressource with ID: ", newVR.ID)
	w.WriteHeader(http.StatusCreated)
}

func (v *vaccationRessource) GetByID(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("component", "read")
	logger.Info("get vaccation-ressource by id")
	vrID, err := extractVaccationRessourceID(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	vr, err := v.store.GetVaccationRessourceByID(r.Context(), vrID)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(vr)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	v.logger.Info("get vaccation-ressource with id: ", vr)
}

func (v *vaccationRessource) List(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("method", "list")
	logger.Info("retrieve vaccation-ressource list")
	list, err := v.store.ListVaccationRessource(r.Context())
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
	v.logger.Info("get list of vaccation-ressource")
}

func (v *vaccationRessource) Update(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("method", "update")
	logger.Info("update vaccation-ressource")
	var vr model.VaccationRessource
	err := json.NewDecoder(r.Body).Decode(&vr)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newVR, err := v.store.UpdateVaccationRessource(r.Context(), &vr)
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
	v.logger.Info("update vaccation-ressource with id: ", newVR.ID)
}

func (v *vaccationRessource) Delete(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("method", "delete")
	logger.Info("delete vaccation-resscource")
	vrID, err := extractVaccationRessourceID(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = v.store.DeleteVaccationRessource(r.Context(), vrID)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	v.logger.Info("delete vaccation-ressource with id: ", vrID)
	w.WriteHeader(http.StatusAccepted)
}

func extractVaccationRessourceID(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	vaccationRessourceID, ok := vars["vaccationRessourceID"]
	if !ok {
		return "", errors.New("could not extract vaccationRessourceID")
	}
	return vaccationRessourceID, nil
}
