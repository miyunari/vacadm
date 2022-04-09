package vacationressources

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MninaTB/vacadm/pkg/database"
	"github.com/MninaTB/vacadm/pkg/model"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func NewVacationRessource(store database.Database, logger logrus.FieldLogger) *vacationRessource {
	return &vacationRessource{
		store:  store,
		logger: logger.WithField("component", "vacation-ressource"),
	}
}

type vacationRessource struct {
	store  database.Database
	logger logrus.FieldLogger
}

func (v *vacationRessource) Create(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("method", "create")
	logger.Info("create new vacation-ressource")
	var vr model.VacationRessource
	err := json.NewDecoder(r.Body).Decode(&vr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error(err)
		return
	}
	newVR, err := v.store.CreateVacationRessource(r.Context(), &vr)
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
	v.logger.Info("create vacation-ressource with ID: ", newVR.ID)
	w.WriteHeader(http.StatusCreated)
}

func (v *vacationRessource) GetByID(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("component", "read")
	logger.Info("get vacation-ressource by id")
	vrID, err := extractVacationRessourceID(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	vr, err := v.store.GetVacationRessourceByID(r.Context(), vrID)
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
	v.logger.Info("get vacation-ressource with id: ", vr)
}

func (v *vacationRessource) List(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("method", "list")
	logger.Info("retrieve vacation-ressource list")
	list, err := v.store.ListVacationRessource(r.Context())
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
	v.logger.Info("get list of vacation-ressource")
}

func (v *vacationRessource) Update(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("method", "update")
	logger.Info("update vacation-ressource")
	var vr model.VacationRessource
	err := json.NewDecoder(r.Body).Decode(&vr)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newVR, err := v.store.UpdateVacationRessource(r.Context(), &vr)
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
	v.logger.Info("update vacation-ressource with id: ", newVR.ID)
}

func (v *vacationRessource) Delete(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("method", "delete")
	logger.Info("delete vacation-resscource")
	vrID, err := extractVacationRessourceID(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = v.store.DeleteVacationRessource(r.Context(), vrID)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	v.logger.Info("delete vacation-ressource with id: ", vrID)
	w.WriteHeader(http.StatusAccepted)
}

func extractVacationRessourceID(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	vacationRessourceID, ok := vars["vacationRessourceID"]
	if !ok {
		return "", errors.New("could not extract vacationRessourceID")
	}
	return vacationRessourceID, nil
}
