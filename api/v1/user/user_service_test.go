package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MninaTB/vacadm/pkg/database/inmemory"
	"github.com/MninaTB/vacadm/pkg/model"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func TestUserService_Create(t *testing.T) {

	usr := &model.User{}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(usr); err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPut, "/v1/user", &buf)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	svc := NewUserService(inmemory.NewInmemoryDB(), logrus.New())

	handler := http.HandlerFunc(svc.Create)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	usr.Email = "test@test.com"
	usr.FirstName = "abc"
	usr.LastName = "def"
	buf.Reset()
	if err := json.NewEncoder(&buf).Encode(usr); err != nil {
		t.Fatal(err)
	}

	req, err = http.NewRequest(http.MethodPut, "/v1/user", &buf)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	got := &model.User{}
	if err := json.NewDecoder(rr.Body).Decode(got); err != nil {
		t.Fatal(err)
	}

	// NOTE: check if uuid is set
	_, err = uuid.Parse(got.ID)
	if err != nil {
		t.Error(err)
	}

	// NOTE: check if createdAt timestamp is set
	if got.CreatedAt == nil {
		t.Error("missing timestamp created_at")
	}

	ignoreFields := cmp.FilterPath(func(p cmp.Path) bool {
		return strings.Contains(p.String(), "ID") ||
			strings.Contains(p.String(), "CreatedAt")
	}, cmp.Ignore())

	// NOTE: compare original struct, ignore ID and CreatedAt (should be different)
	if !cmp.Equal(usr, got, ignoreFields) {
		t.Fatal(cmp.Diff(usr, got, ignoreFields))
	}
}
