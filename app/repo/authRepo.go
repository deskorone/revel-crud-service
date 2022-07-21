package repo

import (
	"database/sql"
)

type AuthRepoImpl struct {
	db *sql.DB
}

var authRepoInstance AuthRepo

func NewAuthRepo(DB *sql.DB) AuthRepo {
	if authRepoInstance == nil {
		authRepoInstance = AuthRepoImpl{db: DB}
	}


	return authRepoInstance
}
