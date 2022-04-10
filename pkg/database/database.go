package database

import (
	"context"

	"github.com/MninaTB/vacadm/pkg/database/inmemory"
	"github.com/MninaTB/vacadm/pkg/database/mariadb"
	"github.com/MninaTB/vacadm/pkg/model"
)

var _ Database = (*inmemory.InmemoryDB)(nil)
var _ Database = (*mariadb.MariaDB)(nil)

// Database is implemented by any structure providing all Database methods,
// defines how models are handled.
type Database interface {
	// CreateUser stores an internal copy of the given user, if email address is
	// not already in use, given parentID and/or teamID exists.
	// Returns copy with assigned userID.
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	// GetUserByID returns the associated user by the given id.
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
	// ListUsers returns a copy of the internal user list.
	ListUsers(ctx context.Context) ([]*model.User, error)
	// UpdateUser updates user entry by the given user.
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
	// DeleteUser removes user entry by the given id.
	DeleteUser(ctx context.Context, userID string) error

	// CreateTeam stores an internal copy of the given team.
	// Returns copy with assigned teamID.
	CreateTeam(ctx context.Context, team *model.Team) (*model.Team, error)
	// GetTeamByID returns the associated team by the given id.
	GetTeamByID(ctx context.Context, teamID string) (*model.Team, error)
	// ListTeams returns a copy of the internal team list.
	ListTeams(ctx context.Context) ([]*model.Team, error)
	// ListTeamUsers returns a list of users associated by the given teamID
	ListTeamUsers(ctx context.Context, teamID string) ([]*model.User, error)
	// UpdateTeam updates team entry by the given team.
	UpdateTeam(ctx context.Context, team *model.Team) (*model.Team, error)
	// DeleteTeam removes team entry by the given id.
	DeleteTeam(ctx context.Context, teamID string) error

	// CreateVacation stores an internal copy of the given vacation resource.
	// Returns copy with assigned vacationID.
	CreateVacation(ctx context.Context, vacation *model.Vacation) (*model.Vacation, error)
	// GetVacationsByTeamID returns the associated vacation by the given teamID.
	GetVacationsByTeamID(ctx context.Context, teamID string) ([]*model.Vacation, error)
	// GetVacationByID returns the associated vacation by the given id.
	GetVacationByID(ctx context.Context, vacationID string) (*model.Vacation, error)
	// ListVacations returns a copy of the internal vacation list.
	ListVacations(ctx context.Context) ([]*model.Vacation, error)
	// DeleteVacation removes vacation entry by the given id.
	DeleteVacation(ctx context.Context, vacationID string) error

	// CreateVacationRequest stores an internal copy of the given vacationRequest.
	// Returns copy with assigned vacationRequestID.
	CreateVacationRequest(ctx context.Context, vacationRequest *model.VacationRequest) (*model.VacationRequest, error)
	// GetVacationRequestByID returns the associated vacationRequest by the given id.
	GetVacationRequestByID(ctx context.Context, vacationRequestID string) (*model.VacationRequest, error)
	// ListVacationRequests returns a copy of the internal vacationRequest list.
	ListVacationRequests(ctx context.Context) ([]*model.VacationRequest, error)
	// UpdateVacationRequest updates vacationRequest entry by the given vacationRequest.
	UpdateVacationRequest(ctx context.Context, vacationRequest *model.VacationRequest) (*model.VacationRequest, error)
	// DeleteVacationRequest removes vacationRequest entry by the given id.
	DeleteVacationRequest(ctx context.Context, vacationRequestID string) error

	// CreateVacationRessource stores an internal copy of the given vacationRessource.
	// Returns copy with assigned vacationRessourceID.
	CreateVacationRessource(ctx context.Context, vacationRessource *model.VacationRessource) (*model.VacationRessource, error)
	// GetVacationRessourceByID returns the associated vacationRessource by the given id.
	GetVacationRessourceByID(ctx context.Context, vacationRessourceID string) (*model.VacationRessource, error)
	// ListVacationRessource returns a copy of the internal vacationRessource list.
	ListVacationRessource(ctx context.Context) ([]*model.VacationRessource, error)
	// UpdateVacationRessource updates vacationRessource entry by the given vacationRessource.
	UpdateVacationRessource(ctx context.Context, vacationRessource *model.VacationRessource) (*model.VacationRessource, error)
	// DeleteVacationRessource removes vacationRessource entry by the given id.
	DeleteVacationRessource(ctx context.Context, vacationRessourceID string) error
}
