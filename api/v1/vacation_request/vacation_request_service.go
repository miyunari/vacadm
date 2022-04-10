package vacationrequest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/MninaTB/vacadm/api/v1/util"
	"github.com/MninaTB/vacadm/pkg/database"
	"github.com/MninaTB/vacadm/pkg/model"
	"github.com/MninaTB/vacadm/pkg/notify"
)

// NewVacationRequestService returns a VacationRequestService.
func NewVacationRequestService(
	store database.Database,
	notifier notify.Notifier,
	logger logrus.FieldLogger,
) *VacationRequestService {
	return &VacationRequestService{
		store:         store,
		relationStore: database.NewRelationDB(store),
		notifier:      notifier,
		logger:        logger.WithField("component", "vacation-request-service"),
	}
}

// VacationRequestService implements http.HandlerFunc's to operate on VacationRequest
// resources.
type VacationRequestService struct {
	store         database.Database
	relationStore database.RelationDB
	notifier      notify.Notifier
	logger        logrus.FieldLogger
}

// Create reads the given payload and creates a store representation accordingly.
func (v *VacationRequestService) Create(w http.ResponseWriter, r *http.Request) {
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

// GetByID extracts a VacationRequestID from URL and writes all user information
// into the given response writer.
func (v *VacationRequestService) GetByID(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("method", "read")
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

// Approve checks if a user has the necessary permissions to approve a request.
// If this is the case, a confirmed Vacation entry is created in the store.
func (v *VacationRequestService) Approve(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("method", "approve")
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

	userID, err := util.UserIDFromRequest(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	parentID, err := util.ParentIDFromRequest(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger = logger.WithFields(logrus.Fields{
		"vac-request": vrID,
		"parentID":    parentID,
		"userID":      userID,
	})

	ok, err := v.relationStore.IsParentUser(r.Context(), userID, parentID)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		logger.Error("missing permission - can not approve")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	logger.Info("approve vacation-request")
	vacation := &model.Vacation{
		UserID:     vR.UserID,
		ApprovedBy: &parentID,
		From:       vR.From,
		To:         vR.To,
	}

	parent, err := v.store.GetUserByID(r.Context(), parentID)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	vac, err := v.store.CreateVacation(r.Context(), vacation)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	msg := fmt.Sprintf(
		"your vacation request '%s', from: %s, to: %s got approved by: %s %s",
		vrID, vR.From.String(), vR.To.String(), parent.FirstName, parent.LastName,
	)
	err = v.notifier.NotifyUser(r.Context(), userID, msg)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(&vac)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// List retuns a list of all VacationRequests available on the internal store.
func (v *VacationRequestService) List(w http.ResponseWriter, r *http.Request) {
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

// Update reads new VacationRequest information from the request body and
// updates the store representation accordingly.
func (v *VacationRequestService) Update(w http.ResponseWriter, r *http.Request) {
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

// Delete a VacationRequest associated to the given VacationRequestID in the URL.
func (v *VacationRequestService) Delete(w http.ResponseWriter, r *http.Request) {
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
