package services

import (
	"github.com/astaxie/beego/logs"
	"github.com/casbin/casbin"
	"github.com/qiangjiyi/casbin_auth/common"
	"github.com/qiangjiyi/casbin_auth/common/errors"
	"github.com/qiangjiyi/casbin_auth/common/util"
	"github.com/qiangjiyi/casbin_auth/models"
	"strconv"
	"strings"
)

// QueryPermission
func QueryPermission(e *casbin.Enforcer, userId string) (*common.QueryPermissionResponse, errors.AuthError) {
	// query user info by userId
	id, _ := strconv.Atoi(userId)
	user, err := models.GetUserById(id)
	if err != nil {
		if err.Error() == "<QuerySeter> no row found" {
			logs.Error("The query user does not exist.")
			return nil, errors.ErrUserNotExistById
		}
	}

	response := &common.QueryPermissionResponse{
		UserId:   user.Id,
		Username: user.Username,
		RealName: user.RealName,
		Mobile:   user.Mobile,
		Email:    user.Email,
	}

	// query user's all roles by casbin.Enforcer
	var roleList []common.UpdateRoleResponse
	roles := getRoles(e, userId)
	for _, r := range roles {
		roleId, _ := strconv.Atoi(strings.Replace(r, "role_", "", -1))
		role, err := models.GetRoleById(roleId)
		if err != nil {
			if err.Error() == "<QuerySeter> no row found" {
				logs.Error("The query role does not exist.")
				return nil, errors.ErrRoleNotExistById
			}
		}
		roleList = append(roleList, common.UpdateRoleResponse{
			Id:          role.Id,
			RoleName:    role.RoleName,
			Description: role.Description,
			UpdateTime:  role.UpdateTime.Format(util.DateTimeFormat),
		})
	}
	response.RoleList = roleList

	// query user's all permissions by casbin.Enforcer
	var permissionList []string
	for _, p := range e.GetPermissionsForUser(userId) {
		permissionList = append(permissionList, strings.Join(p, ","))
	}
	for _, r := range roles { // query every role's permission
		for _, p := range e.GetPermissionsForUser(r) {
			permissionList = append(permissionList, strings.Join(p, ","))
		}
	}
	response.PermissionList = permissionList

	return response, errors.Success
}

// getRoles Recursive find roles by user name and add result to string slice
func getRoles(e *casbin.Enforcer, name string) (roles []string) {
	for _, r := range e.GetRolesForUser(name) {
		// remove repeat role
		var flag = false
		for _, v := range roles {
			if v == r {
				flag = true
				break
			}
		}
		if !flag {
			roles = append(roles, r)
		}

		roles = append(roles, getRoles(e, r)...)
	}
	return
}

// AddPermission
func AddPermission(e *casbin.Enforcer, param common.AddPermissionRequest) (*common.AddPermissionResponse, errors.AuthError) {
	switch param.Type {
	case "p": // permission
		var sub, obj, action string

		// set sub value by sub type
		sub, err := setSub(param.SubType, param.Sub)
		if err != nil {
			return nil, err
		}

		// set obj value
		obj = param.Obj

		// set action value
		action = param.Action

		if e.AddPermissionForUser(sub, obj, action) {
			return &common.AddPermissionResponse{
				Policy: "p, " + sub + ", " + obj + ", " + action,
			}, errors.Success
		} else { // the permission add has existed.
			logs.Error("Add permission failed because of has existed.")
			return nil, errors.ErrPermissionHasExisted
		}
	case "g": // group
		var sub, obj string

		// set sub value by sub type
		sub, err := setSub(param.SubType, param.Sub)
		if err != nil {
			return nil, err
		}

		// set obj value by obj type
		// query role by obj verify if obj is a valid role
		roleId, _ := strconv.Atoi(param.Obj)
		if _, err := models.GetRoleById(roleId); err != nil {
			logs.Error("Query role by id error. ---> %+v", err)
			return nil, errors.ErrPermissionObjInvalid
		}
		obj = "role_" + param.Obj

		if e.AddRoleForUser(sub, obj) {
			return &common.AddPermissionResponse{
				Policy: "g, " + sub + ", " + obj,
			}, errors.Success
		} else {
			logs.Error("Add permission group failed because of has existed.")
			return nil, errors.ErrPermissionHasExisted
		}
	default:
		return nil, errors.ErrPermissionTypeInvalid
	}

}

// setSub set sub value by sub type
func setSub(subType string, subValue string) (string, errors.AuthError) {

	switch subType {
	case "u": // user
		// query user by sub verify if sub is a valid user
		if user, exist := models.GetUserByUsername(subValue); !exist {
			return "", errors.ErrPermissionSubInvalid
		} else {
			return strconv.Itoa(user.Id), nil
		}
	case "r": // role
		// query role by sub verify if sub is a valid role
		roleId, _ := strconv.Atoi(subValue)
		if _, err := models.GetRoleById(roleId); err != nil {
			logs.Error("Query role by id error. ---> %+v", err)
			return "", errors.ErrPermissionSubInvalid
		}
		return "role_" + subValue, nil
	default:
		return "", errors.ErrPermissionSubTypeInvalid
	}
}
