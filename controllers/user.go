package controllers

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/utils"
	"github.com/qiangjiyi/casbin_auth/cache"
	"github.com/qiangjiyi/casbin_auth/common"
	"github.com/qiangjiyi/casbin_auth/common/errors"
	"github.com/qiangjiyi/casbin_auth/common/util"
	"github.com/qiangjiyi/casbin_auth/services"

	"strconv"
	"time"
)

type UserController struct {
	BaseController
}

// Register
func (u *UserController) Register() {
	var param common.RegisterRequest

	if err := util.ParseRequestParams(u.Ctx.Input.RequestBody, &param); err != nil {
		u.Response(nil, err)
		return
	}

	// TODO validate input params

	if payload, err := services.Register(param); err.Code() != errors.SuccessCode {
		u.Response(nil, err)
	} else {
		u.Response(payload, err)
	}
}

// Login
func (u *UserController) Login() {
	var param common.LoginRequest

	if err := util.ParseRequestParams(u.Ctx.Input.RequestBody, &param); err != nil {
		u.Response(nil, err)
		return
	}

	// TODO validate input params
	/*if u.ParseAndValidate(&params) == false {
		return
	}*/

	payload, err := services.Login(param)

	if err.Code() != errors.SuccessCode { // if login failed
		u.Response(nil, err)
		return
	}

	if u.SessionOn {
		// set userId to session for current user. browser/client will receive by cookie
		u.SetSession(util.SessionKey, payload.UserId)
	} else {
		// set token to cookie
		// u.Ctx.SetCookie("Authorization", token)

		redis := cache.RedisCacheForSession()
		if redis == nil {
			u.Response(nil, errors.ErrInternalServer)
			return
		}

		// generate a token for current user
		token := string(utils.RandomCreateBytes(32))
		if err := redis.Put(token, payload.UserId, 30*time.Minute); err != nil {
			logs.Error("Redis put key[%s] failed. ---> %+v", token, err)
			u.Response(nil, errors.ErrInternalServer)
			return
		}
		payload.Token = token
	}

	u.Response(payload, err)
}

// Logout exit system
func (u *UserController) Logout() {
	if u.SessionOn {
		// delete session id from SessionStore.values
		u.DelSession(util.SessionKey)
		u.Response(nil, errors.Success)
	} else {
		// delete session id from redis cache
		redis := cache.RedisCacheForSession()
		if redis == nil {
			u.Response(nil, errors.ErrInternalServer)
			return
		}

		if _, ok := redis.Get(u.Ctx.Input.Header(util.TokenCacheKey)).([]byte); ok {
			if err := redis.Delete(u.Ctx.Input.Header(util.TokenCacheKey)); err == nil { // delete token from redis cache
				u.Response(nil, errors.Success)
				return
			}
		}
		u.Response(nil, errors.ErrUserLogoutFail)
	}
}

// UpdateUser update user info
func (u *UserController) UpdateUser() {
	var param common.RegisterRequest

	if err := util.ParseRequestParams(u.Ctx.Input.RequestBody, &param); err != nil {
		u.Response(nil, err)
		return
	}

	// TODO validate input params

	if payload, err := services.Update(param); err.Code() != errors.SuccessCode {
		u.Response(nil, err)
	} else {
		u.Response(payload, err)
	}
}

// DeleteUser
func (u *UserController) DeleteUser() {
	// obtain user id
	userId := u.Ctx.Input.Param(":userId")

	// verify user id is a valid value
	if id, err := strconv.Atoi(userId); err != nil {
		logs.Error("The user id to be deleted invalid ---> %+v", err)
		u.Response(nil, errors.ErrUserIdInvalid)
	} else {
		u.Response(nil, services.Delete(u.Enforcer, id))
	}
}
