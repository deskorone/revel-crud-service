package repo

import (
	"database/sql"
	"testAuth/app/models"
)

type AuthRepoImpl struct {
	db *sql.DB
}

// Registration implements AuthRepo
func (AuthRepoImpl) Registration(u models.User) (*models.User, error) {
	panic("unimplemented")
}

var authRepoInstance AuthRepo

func NewAuthRepo(DB *sql.DB) AuthRepo {
	if authRepoInstance == nil {
		authRepoInstance = AuthRepoImpl{db: DB}
	}


	return authRepoInstance
}
