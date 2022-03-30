package database

import (
	"context"

	"github.com/MninaTB/vacadm/pkg/model"
)

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

	GetVaccationByID(context.Context, string) (*model.Vaccation, error)
	ListVaccations(context.Context) ([]*model.Vaccation, error)
	DeleteVaccation(context.Context, string) error

	CreateVaccationRequest(context.Context, *model.VaccationRequest) (*model.VaccationRequest, error)
	GetVaccationRequestByID(context.Context, string) (*model.VaccationRequest, error)
	ListVaccationRequests(context.Context) ([]*model.VaccationRequest, error)
	UpdateVaccationRequest(context.Context, *model.VaccationRequest) (*model.VaccationRequest, error)
	DeleteVaccationRequest(context.Context, string) error

	CreateVaccationRessource(context.Context, *model.VaccationRessource) (*model.VaccationRessource, error)
	GetVaccationRessourceByID(context.Context, string) (*model.VaccationRessource, error)
	ListVaccationRessource(context.Context) ([]*model.VaccationRessource, error)
	UpdateVaccationRessource(context.Context, *model.VaccationRessource) (*model.VaccationRessource, error)
	DeleteVaccationRessource(context.Context, string) error
}
