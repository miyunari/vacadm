package vacation

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/MninaTB/vacadm/pkg/database"
)

// NewVacation returns a VacationService.
func NewVacationService(store database.Database, logger logrus.FieldLogger) *VacationService {
	return &VacationService{
		store:  store,
		logger: logger.WithField("component", "vacation-service"),
	}
}

// VacationService implements http.HandlerFunc's to operate on user resources.
type VacationService struct {
	store  database.Database
	logger logrus.FieldLogger
}

// GetByID extracts a vacationID from URL and writes all user information into the
// given response writer.
func (v *VacationService) GetByID(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("method", "read")
	logger.Info("get vacation by id")
	vacID, err := extractVacationID(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	vacation, err := v.store.GetVacationByID(r.Context(), vacID)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(&vacation)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	v.logger.Info("get vacation with id: ", vacID)
}

// List retuns a list of all vacations available on the internal store.
func (v *VacationService) List(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("method", "list")
	logger.Info("get vacation list")
	list, err := v.store.ListVacations(r.Context())
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if r.Header.Get("Content-Type") == "application/csv" {
		logger.Info("csv requested")
		csvWriter := csv.NewWriter(w)
		for _, l := range list {
			approvedBy := ""
			if l.ApprovedBy != nil {
				approvedBy = *l.ApprovedBy
			}
			var createdAt string
			if l.CreatedAt != nil {
				createdAt = l.CreatedAt.String()
			}
			var deletedAt string
			if l.DeletedAt != nil {
				deletedAt = l.DeletedAt.String()
			}
			err = csvWriter.Write([]string{l.ID, l.UserID, approvedBy, l.From.String(), l.To.String(), createdAt, deletedAt})
			if err != nil {
				logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		csvWriter.Flush()
		return
	}
	err = json.NewEncoder(w).Encode(&list)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Delete a vacation associated to the given vacationID in the URL.
func (v *VacationService) Delete(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("method", "delete")
	logger.Info("delete vacation")
	vacID, err := extractVacationID(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = v.store.DeleteVacation(r.Context(), vacID)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	v.logger.Info("delete vacation with id: ", vacID)
	w.WriteHeader(http.StatusAccepted)
}

func extractVacationID(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	vacationID, ok := vars["vacationID"]
	if !ok {
		return "", errors.New("could not extract vacationID")
	}
	return vacationID, nil
}
