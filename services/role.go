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

// CreateRole
func CreateRole(param common.RoleRequest) (*common.CreateRoleResponse, errors.AuthError) {
	role := models.Role{
		RoleName:    param.RoleName,
		Description: param.Description,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}

	if roleId, err := models.AddRole(&role); err != nil {
		// print error log info and return
		logs.Error("CreateRole failed ---> %+v", err)
		return nil, errors.ErrDatabaseOperation
	} else {
		logs.Info("CreateRole success. New roleId ---> %d", roleId)
		param.Id = int(roleId)
		return &common.CreateRoleResponse{
			Id:          int(roleId),
			RoleName:    role.RoleName,
			Description: role.Description,
			CreateTime:  role.CreateTime.Format(util.DateTimeFormat),
		}, errors.Success
	}
}

// UpdateRole
func UpdateRole(param common.RoleRequest) (*common.UpdateRoleResponse, errors.AuthError) {
	role := models.Role{
		Id:          param.Id,
		RoleName:    param.RoleName,
		Description: param.Description,
		UpdateTime:  time.Now(),
	}

	if err := models.UpdateRoleById(&role); err != nil {
		if err.Error() == "<QuerySeter> no row found" { // to be updated role not exist.
			logs.Error("The role to be updated does not exist.")
			return nil, errors.ErrRoleNotExistById
		}

		logs.Error("UpdateRole failed ---> %+v", err)
		return nil, errors.ErrDatabaseOperation
	} else {
		logs.Info("UpdateRole [roleId = %d] success.", param.Id)
		return &common.UpdateRoleResponse{
			Id:          role.Id,
			RoleName:    role.RoleName,
			Description: role.Description,
			UpdateTime:  role.UpdateTime.Format(util.DateTimeFormat),
		}, errors.Success
	}
}

// DeleteRole
func DeleteRole(e *casbin.Enforcer, id int) errors.AuthError {
	if err := models.DeleteRole(id); err != nil {
		if err.Error() == "<QuerySeter> no row found" { // to be deleted role not exist.
			logs.Error("The role to be deleted does not exist.")
			return errors.ErrRoleNotExistById
		}

		logs.Error("DeleteRole failed ---> %+v", err)
		return errors.ErrDatabaseOperation
	} else {
		// sync delete user, role and permission related to role by casbin.Enforcer
		e.DeleteRole("role_" + strconv.Itoa(id))
		e.DeleteRolesForUser("role_" + strconv.Itoa(id))
		logs.Info("DeleteRole [roleId = %d] success.", id)
		return errors.Success
	}
}
