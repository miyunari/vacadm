package vacation

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MninaTB/vacadm/pkg/database"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func NewVacation(store database.Database, logger logrus.FieldLogger) *vacation {
	return &vacation{
		store:  store,
		logger: logger.WithField("component", "vacation-service"),
	}
}

type vacation struct {
	store  database.Database
	logger logrus.FieldLogger
}

func (v *vacation) GetByID(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("component", "read")
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

func (v *vacation) List(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("component", "list")
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
				approvedBy = l.ApprovedBy.ID
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

func (v *vacation) Delete(w http.ResponseWriter, r *http.Request) {
	logger := v.logger.WithField("component", "delete")
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
