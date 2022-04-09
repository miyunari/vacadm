package database

import (
	"context"

	"github.com/MninaTB/vacadm/pkg/database/inmemory"
	"github.com/MninaTB/vacadm/pkg/database/mariadb"
	"github.com/MninaTB/vacadm/pkg/model"
)

var _ Database = (*inmemory.InmemoryDB)(nil)
var _ Database = (*mariadb.MariaDB)(nil)

type Database interface {
	CreateUser(context.Context, *model.User) (*model.User, error)
	GetUserByID(context.Context, string) (*model.User, error)
	ListUsers(context.Context) ([]*model.User, error)
	UpdateUser(context.Context, *model.User) (*model.User, error)
	DeleteUser(context.Context, string) error

	CreateTeam(context.Context, *model.Team) (*model.Team, error)
	GetTeamByID(context.Context, string) (*model.Team, error)
	ListTeams(context.Context) ([]*model.Team, error)
	ListTeamUsers(context.Context, string) ([]*model.User, error)
	UpdateTeam(context.Context, *model.Team) (*model.Team, error)
	DeleteTeam(context.Context, string) error

	GetVacationByID(context.Context, string) (*model.Vacation, error)
	ListVacations(context.Context) ([]*model.Vacation, error)
	DeleteVacation(context.Context, string) error

	CreateVacationRequest(context.Context, *model.VacationRequest) (*model.VacationRequest, error)
	GetVacationRequestByID(context.Context, string) (*model.VacationRequest, error)
	ListVacationRequests(context.Context) ([]*model.VacationRequest, error)
	UpdateVacationRequest(context.Context, *model.VacationRequest) (*model.VacationRequest, error)
	DeleteVacationRequest(context.Context, string) error

	CreateVacationRessource(context.Context, *model.VacationRessource) (*model.VacationRessource, error)
	GetVacationRessourceByID(context.Context, string) (*model.VacationRessource, error)
	ListVacationRessource(context.Context) ([]*model.VacationRessource, error)
	UpdateVacationRessource(context.Context, *model.VacationRessource) (*model.VacationRessource, error)
	DeleteVacationRessource(context.Context, string) error
}
