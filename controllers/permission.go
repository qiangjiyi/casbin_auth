package controllers

import (
	"github.com/astaxie/beego/logs"
	"github.com/qiangjiyi/casbin_auth/common"
	"github.com/qiangjiyi/casbin_auth/common/errors"
	"github.com/qiangjiyi/casbin_auth/common/util"
	"github.com/qiangjiyi/casbin_auth/services"
	"strconv"
)

type PermissionController struct {
	BaseController
}

// QueryPermission query user's all roles and permissions
func (p *PermissionController) QueryPermission() {
	// obtain user id
	userId := p.Ctx.Input.Param(":userId")

	// verify user id is a valid value
	if _, err := strconv.Atoi(userId); err != nil {
		logs.Error("The user id is invalid ---> %+v", err)
		p.Response(nil, errors.ErrUserIdInvalid)
	} else {
		if payload, err := services.QueryPermission(p.Enforcer, userId); err.Code() != errors.SuccessCode {
			p.Response(nil, err)
		} else {
			p.Response(payload, err)
		}
	}
}

// AddPermission add a new permission rule in casbin_rule
func (p *PermissionController) AddPermission() {

	// parse request param
	var param common.AddPermissionRequest
	if err := util.ParseRequestParams(p.Ctx.Input.RequestBody, &param); err != nil {
		p.Response(nil, err)
		return
	}

	if payload, err := services.AddPermission(p.Enforcer, param); err.Code() != errors.SuccessCode {
		p.Response(nil, err)
	} else {
		p.Response(payload, err)
	}
}
