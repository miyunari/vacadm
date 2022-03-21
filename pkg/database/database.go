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
	UpdateTeam(*model.Team) (*model.Team, error)
	DeleteTeam(string) error

	CreateVaccation(*model.Vaccation) (*model.Vaccation, error)
	GetVaccationByID(string) (*model.Vaccation, error)
	ListVaccations() ([]*model.Vaccation, error)
	UpdateVaccation(*model.Vaccation) (*model.Vaccation, error)
	DeleteVaccation(string) error
}
