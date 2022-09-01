package repo

import (
	"database/sql"
	"testAuth/app/models"
)

type UserRepoImpl struct {
	DB *sql.DB
}

const (
	findByNameQ = "select * from usr u where u.username = $1"
	saveUserQ   = "insert into usr (username, password, balance) values ($1, $2, $3) returning *"
	findUserQ   = "select * from usr u where u.user_id = $1"
)

// FindUserByName implements UserRepo поиск по имени
func (c *UserRepoImpl) FindUserByName(Name string) (*models.User, error) {
	u := models.User{}
	err := c.DB.QueryRow(findByNameQ, Name).Scan(&u.Id, &u.Name, &u.Password, &u.Balance)
	return &u, err
}

// Save implements UserRepo сохрание пользователя
func (c *UserRepoImpl) Save(u models.User) (*models.User, error) {
	err := c.DB.QueryRow(saveUserQ, u.Name, u.Password, u.Balance).Scan(&u.Id, &u.Name, &u.Password, &u.Balance)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// FindUserById implements UserRepo поиск потльзователя по имени
func (c *UserRepoImpl) FindUserById(id int) (*models.User, error) {
	u := models.User{}
	err := c.DB.QueryRow(findUserQ, id).Scan(&u.Id, &u.Name, &u.Password, &u.Balance)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

var userRepoInstance UserRepo

func NewUserRepo(DB *sql.DB) UserRepo {
	if userRepoInstance == nil {
		userRepoInstance = &UserRepoImpl{DB: DB}
	}
	return userRepoInstance
}
