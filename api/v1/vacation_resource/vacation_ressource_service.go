package vacationresources

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/MninaTB/vacadm/pkg/database"
	"github.com/MninaTB/vacadm/pkg/model"
)

// NewVacationResourceService returns a VacationResourceService.
func NewVacationResourceService(
	store database.Database,
	logger logrus.FieldLogger,
) *VacationResourceService {
	return &VacationResourceService{
		store:  store,
		logger: logger.WithField("component", "vacation-resource"),
	}
}

// VacationResourceService http.HandlerFunc's to operate on VacationResource resources.
type VacationResourceService struct {
	store  database.Database
	logger logrus.FieldLogger
}

// Create reads the given payload and creates a store representation accordingly.
func (v *VacationResourceService) Create(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("method", "create")
	logger.Info("create new vacation-resource")
	var vr model.VacationResource
	err := json.NewDecoder(r.Body).Decode(&vr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error(err)
		return
	}
	newVR, err := v.store.CreateVacationResource(r.Context(), &vr)
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
	v.logger.Info("create vacation-resource with ID: ", newVR.ID)
	w.WriteHeader(http.StatusCreated)
}

// GetByID extracts a VacationResourceID from URL and writes all user information
// into the given response writer.
func (v *VacationResourceService) GetByID(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("component", "read")
	logger.Info("get vacation-resource by id")
	vrID, err := extractVacationResourceID(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	vr, err := v.store.GetVacationResourceByID(r.Context(), vrID)
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
	v.logger.Info("get vacation-resource with id: ", vr)
}

// List retuns a list of all VacationResources available on the internal store.
func (v *VacationResourceService) List(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("method", "list")
	logger.Info("retrieve vacation-resource list")
	list, err := v.store.ListVacationResource(r.Context())
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
	v.logger.Info("get list of vacation-resource")
}

// Update reads new VacationResource information from the request body and updates the store
// representation accordingly.
func (v *VacationResourceService) Update(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("method", "update")
	logger.Info("update vacation-resource")
	var vr model.VacationResource
	err := json.NewDecoder(r.Body).Decode(&vr)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newVR, err := v.store.UpdateVacationResource(r.Context(), &vr)
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
	v.logger.Info("update vacation-resource with id: ", newVR.ID)
}

// Delete a VacationResource associated to the given VacationResourceID in the URL.
func (v *VacationResourceService) Delete(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("method", "delete")
	logger.Info("delete vacation-resscource")
	vrID, err := extractVacationResourceID(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = v.store.DeleteVacationResource(r.Context(), vrID)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	v.logger.Info("delete vacation-resource with id: ", vrID)
	w.WriteHeader(http.StatusAccepted)
}

func extractVacationResourceID(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	vacationResourceID, ok := vars["vacationResourceID"]
	if !ok {
		return "", errors.New("could not extract vacationResourceID")
	}
	return vacationResourceID, nil
}
