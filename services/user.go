package services

import (
	"github.com/astaxie/beego/logs"
	"github.com/casbin/casbin"
	"github.com/qiangjiyi/casbin_auth/common"
	"github.com/qiangjiyi/casbin_auth/common/errors"
	"github.com/qiangjiyi/casbin_auth/common/util"
	"github.com/qiangjiyi/casbin_auth/models"
	"strconv"
	"time"
)

// Register
func Register(param common.RegisterRequest) (*common.RegisterResponse, errors.AuthError) {
	// query user by username and return username error if has existed.
	_, exist := models.GetUserByUsername(param.Username)
	if exist {
		logs.Error("Request username %q has existed in database", param.Username)
		return nil, errors.ErrUsernameHasExisted
	}

	user := models.User{
		Username:   param.Username,
		Password:   util.Md5Encrypt(param.Password), // use md5 hash algorithm for request password encrypt
		RealName:   param.RealName,
		Mobile:     param.Mobile,
		Email:      param.Email,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	if userId, err := models.AddUser(&user); err != nil {
		// print error log info and return
		logs.Error("Register user failed ---> %+v", err)
		return nil, errors.ErrDatabaseOperation
	} else {
		logs.Info("Register user success. New userId ---> %d", userId)
		return &common.RegisterResponse{
			UserId:     int(userId),
			Username:   user.Username,
			CreateTime: user.CreateTime.Format(util.DateTimeFormat),
		}, errors.Success
	}
}

// Login
func Login(param common.LoginRequest) (*common.LoginResponse, errors.AuthError) {
	// query user by username and return username error if not exist.
	user, exist := models.GetUserByUsername(param.Username)
	if !exist {
		logs.Error("Not found user by username: %s", param.Username)
		return nil, errors.ErrUsernameNotCorrect
	}

	// using md5 hash algorithm for request password encrypt to compare with database
	// if not equal return password not correct error.
	if user.Password != util.Md5Encrypt(param.Password) {
		logs.Error("Request password is not matched with database")
		return nil, errors.ErrPasswordNotCorrect
	}

	logs.Info("Login success...")
	return &common.LoginResponse{UserId: user.Id}, errors.Success
}

// Update
func Update(param common.RegisterRequest) (*common.UpdateUserResponse, errors.AuthError) {
	user := models.User{
		Id:         param.Id,
		Username:   param.Username,
		Password:   util.Md5Encrypt(param.Password), // use md5 hash algorithm for request password encrypt
		RealName:   param.RealName,
		Mobile:     param.Mobile,
		Email:      param.Email,
		UpdateTime: time.Now(),
	}

	if err := models.UpdateUserById(&user); err != nil {
		// print error log info and return
		if err.Error() == "<QuerySeter> no row found" { // to be updated user not exist.
			logs.Error("The user to be updated does not exist.")
			return nil, errors.ErrUserNotExistById
		}

		logs.Error("Update user failed ---> %+v", err)
		return nil, errors.ErrDatabaseOperation
	} else {
		logs.Info("Update user [userId = %d] success.", param.Id)
		return &common.UpdateUserResponse{
			UserId:     user.Id,
			Username:   user.Username,
			UpdateTime: user.UpdateTime.Format(util.DateTimeFormat),
		}, errors.Success
	}
}

// Delete
func Delete(e *casbin.Enforcer, id int) errors.AuthError {
	if err := models.DeleteUser(id); err != nil {
		// print error log info and return
		if err.Error() == "<QuerySeter> no row found" { // to be deleted user not exist.
			logs.Error("The user to be deleted does not exist.")
			return errors.ErrUserNotExistById
		}

		logs.Error("Delete user failed ---> %+v", err)
		return errors.ErrDatabaseOperation
	} else {
		// sync delete role and permission related to user by casbin.Enforcer
		e.DeleteUser(strconv.Itoa(id))
		e.DeletePermissionsForUser(strconv.Itoa(id))
		logs.Info("Delete user [userId = %d] success.", id)
		return errors.Success
	}
}
