package service

import (
	"context"
	"database/sql"
	"errors"
	"roomino/ctl"
	"roomino/dao"
	"roomino/model"
	"roomino/types"
	"roomino/util"
	"sync"
)

type UserSrv struct {
}

var UserSrvIns *UserSrv
var UserSrvOnce sync.Once

func GetUserSrv() *UserSrv {
	UserSrvOnce.Do(func() {
		UserSrvIns = &UserSrv{}
	})
	return UserSrvIns
}

func (s *UserSrv) Register(ctx context.Context, req *types.UserServiceReq) (interface{}, error) {
	userDao := dao.NewUserDao(ctx)
	u, err := userDao.FindUserByUserName(req.Username)
	if err == sql.ErrNoRows {
		u = &model.Users{
			Username:  req.Username,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			DOB:       req.DOB,
			Gender:    req.Gender,
			Email:     req.Email,
			Phone:     req.Phone,
		}
		if err := u.SetPassword(req.Passwd); err != nil {
			return nil, err
		}

		if err := userDao.CreateUser(u); err != nil {
			return nil, err
		}
		return ctl.RespSuccess(), nil
	}
	if err == nil {
		return nil, errors.New("Userexists")
	}
	return nil, err
}

func (s *UserSrv) Login(ctx context.Context, req *types.UserServiceReq) (resp interface{}, err error) {
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.FindUserByUserName(req.Username)
	if err == sql.ErrNoRows {
		err = errors.New("UsernotExist")
		return
	}
	if !user.CheckPassword(req.Passwd) {
		err = errors.New("WrongPassword")
		return
	}
	token, err := util.GenerateToken(req.Username, 0)
	if err != nil {
		return
	}
	u := &types.UserResp{
		UserName: user.Username,
	}
	uResp := &types.TokenData{
		User:  u,
		Token: token,
	}
	return ctl.RespSuccessWithData(uResp), nil
}

func (s *UserSrv) GetUserProfile(ctx context.Context) (interface{}, error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return nil, errors.New("failed to get user info")
	}
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.FindUserByUserName(u.UserName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to retrieve user info")
	}
	userProfile := types.UserProfile{
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		DOB:       user.DOB,
		Gender:    user.Gender,
		Email:     user.Email,
		Phone:     user.Phone,
	}
	return ctl.RespSuccessWithData(userProfile), nil
}

func (s *UserSrv) InterestUserProfile(ctx context.Context, req types.UserResp) (interface{}, error) {
	username := req.UserName
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.FindUserByUserName(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to retrieve user info")
	}
	userProfile := types.UserProfile{
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		DOB:       user.DOB,
		Gender:    user.Gender,
		Email:     user.Email,
		Phone:     user.Phone,
	}
	return ctl.RespSuccessWithData(userProfile), nil
}
