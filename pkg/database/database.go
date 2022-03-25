package database

import (
	"github.com/MninaTB/vacadm/pkg/model"
)

type Database interface {
	CreateUser(*model.User) (*model.User, error)
	GetUserByID(string) (*model.User, error)
	ListUsers() ([]*model.User, error)
	UpdateUser(*model.User) (*model.User, error)
	DeleteUser(string) error

	CreateTeam(*model.Team) (*model.Team, error)
	GetTeamByID(string) (*model.Team, error)
	ListTeams() ([]*model.Team, error)
	ListTeamUsers(string) ([]*model.User, error)
	UpdateTeam(*model.Team) (*model.Team, error)
	DeleteTeam(string) error

	GetVaccationByID(string) (*model.Vaccation, error)
	ListVaccations() ([]*model.Vaccation, error)
	DeleteVaccation(string) error

	CreateVaccationRequest(*model.VaccationRequest) (*model.VaccationRequest, error)
	GetVaccationRequestByID(string) (*model.VaccationRequest, error)
	ListVaccationRequests() ([]*model.VaccationRequest, error)
	UpdateVaccationRequest(*model.VaccationRequest) (*model.VaccationRequest, error)
	DeleteVaccationRequest(string) error

	CreateVaccationRessource(*model.VaccationRessource) (*model.VaccationRessource, error)
	GetVaccationRessourceByID(string) (*model.VaccationRessource, error)
	ListVaccationRessource() ([]*model.VaccationRessource, error)
	UpdateVaccationRessource(*model.VaccationRessource) (*model.VaccationRessource, error)
	DeleteVaccationRessource(string) error
}
