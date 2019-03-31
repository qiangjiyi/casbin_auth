package common

/*
This file contains all own struct define related to auth module.
 */

type RegisterRequest struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	RealName string `json:"real_name"`
	Mobile   string `json:"mobile"`
	Email    string `json:"email"`
}

type RegisterResponse struct {
	UserId     int    `json:"user_id"`
	Username   string `json:"username"`
	CreateTime string `json:"create_time"`
}

type UpdateUserResponse struct {
	UserId     int    `json:"user_id"`
	Username   string `json:"username"`
	UpdateTime string `json:"update_time"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserId int    `json:"user_id"`
	Token  string `json:"token,omitempty"`
}

type RoleRequest struct {
	Id          int    `json:"id"`
	RoleName    string `json:"role_name"`
	Description string `json:"description"`
}

type CreateRoleResponse struct {
	Id          int    `json:"id"`
	RoleName    string `json:"role_name"`
	Description string `json:"description"`
	CreateTime  string `json:"create_time"`
}

type UpdateRoleResponse struct {
	Id          int    `json:"id"`
	RoleName    string `json:"role_name"`
	Description string `json:"description"`
	UpdateTime  string `json:"update_time"`
}

type AddPermissionRequest struct {
	Type    string `json:"type"`
	Sub     string `json:"sub"`
	SubType string `json:"sub_type"`
	Obj     string `json:"obj"`
	Action  string `json:"action"`
}

type AddPermissionResponse struct {
	Policy string `json:"policy"`
}

type QueryPermissionResponse struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	RealName string `json:"real_name"`
	Mobile   string `json:"mobile"`
	Email    string `json:"email"`

	RoleList       []UpdateRoleResponse `json:"role_list"`
	PermissionList []string             `json:"permission_list"`
}
