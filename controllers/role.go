package controllers

import (
	"github.com/astaxie/beego/logs"
	"github.com/qiangjiyi/casbin_auth/common"
	"github.com/qiangjiyi/casbin_auth/common/errors"
	"github.com/qiangjiyi/casbin_auth/common/util"
	"github.com/qiangjiyi/casbin_auth/services"
	"strconv"
)

type RoleController struct {
	BaseController
}

// CreateRole create a new role according request param
func (r *RoleController) CreateRole() {
	var param common.RoleRequest

	if err := util.ParseRequestParams(r.Ctx.Input.RequestBody, &param); err != nil {
		r.Response(nil, err)
		return
	}

	// TODO validate input params

	if payload, err := services.CreateRole(param); err.Code() != errors.SuccessCode {
		r.Response(nil, err)
	} else {
		r.Response(payload, err)
	}
}

// UpdateRole update role info
func (r *RoleController) UpdateRole() {
	var param common.RoleRequest

	if err := util.ParseRequestParams(r.Ctx.Input.RequestBody, &param); err != nil {
		r.Response(nil, err)
		return
	}

	// TODO validate input params

	if payload, err := services.UpdateRole(param); err.Code() != errors.SuccessCode {
		r.Response(nil, err)
	} else {
		r.Response(payload, err)
	}
}

// DeleteRole delete role info by role id
func (r *RoleController) DeleteRole() {
	// obtain role id
	roleId := r.Ctx.Input.Param(":roleId")

	// verify role id is a valid value
	if id, err := strconv.Atoi(roleId); err != nil {
		logs.Error("The role id to be deleted invalid ---> %+v", err)
		r.Response(nil, errors.ErrRoleIdInvalid)
	} else {
		r.Response(nil, services.DeleteRole(r.Enforcer, id))
	}

}
