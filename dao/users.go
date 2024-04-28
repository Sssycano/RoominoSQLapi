package dao

import (
	"context"
	"database/sql"

	"roomino/model"
)

type UserDao struct {
	DB *sql.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	if ctx == nil {
		ctx = context.Background()
	}
	return &UserDao{NewDBClient(ctx)}
}

func (dao *UserDao) FindUserByUserName(userName string) (*model.Users, error) {
	query := "SELECT username, first_name, last_name, DOB, gender, email, phone, passwd FROM users WHERE username = ?"
	row := dao.DB.QueryRow(query, userName)
	user := &model.Users{}
	err := row.Scan(&user.Username, &user.FirstName, &user.LastName, &user.DOB, &user.Gender, &user.Email, &user.Phone, &user.Passwd)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (dao *UserDao) CreateUser(user *model.Users) error {
	query := "INSERT INTO users (username, first_name, last_name, DOB, gender, email, phone, passwd) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := dao.DB.Exec(query, user.Username, user.FirstName, user.LastName, user.DOB, user.Gender, user.Email, user.Phone, user.Passwd)
	if err != nil {
		return err
	}
	return nil
}
