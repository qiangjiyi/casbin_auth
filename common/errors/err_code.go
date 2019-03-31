package errors

// define error type
type ErrorCodeType uint32

// success code
const SuccessCode ErrorCodeType = 0

// common error code and message defined
const (
	InvalidParamsErrCode    ErrorCodeType = 100000 // 参数不合法
	MissingParamsErrCode    ErrorCodeType = 100001 // 缺少参数
	DatabaseOperationFailed ErrorCodeType = 100002 // 数据库操作失败
	ParseRequestParamsError ErrorCodeType = 100003 // 解析请求体失败
	SerializeDataFail       ErrorCodeType = 100004 // 序列化数据失败
	DeserializeDataFail     ErrorCodeType = 100005 // 反序列化(解析)数据失败
	GetServerContextFail    ErrorCodeType = 100006 // 获取服务的上下文失败
	PermissionDenied        ErrorCodeType = 100007 // 没有权限
	InternalServerError     ErrorCodeType = 100008 // 服务内部错误
	BashShellRunFailed      ErrorCodeType = 100009 // shell脚本运行失败
)

// auth module error code
const (
	UserIdInvalid      ErrorCodeType = 100100 // 用户ID不合法
	UserNotExistById   ErrorCodeType = 100101 // 用户ID对应的用户不存在
	UsernameHasExisted ErrorCodeType = 100102 // 用户名已经被占用
	UsernameNotCorrect ErrorCodeType = 100103 // 用户名不正确
	PasswordNotCorrect ErrorCodeType = 100104 // 密码不正确
	NotLoginStatus     ErrorCodeType = 100105 // 未登录状态
	UserLogoutFail     ErrorCodeType = 100106 // 注销登录失败

	RoleIdInvalid    ErrorCodeType = 100200 // 角色ID不合法
	RoleNotExistById ErrorCodeType = 100201 // 角色ID对应的角色不存在

	PermissionEnforcerLoadFailed ErrorCodeType = 100300 // 权限执行模块加载失败
	PermissionTypeInvalid        ErrorCodeType = 100301 // 权限类型不合法 ['p' 或 'g']
	PermissionSubTypeInvalid     ErrorCodeType = 100302 // 权限主体类型不合法 ['u' 或 'r']
	PermissionSubInvalid         ErrorCodeType = 100303 // 权限主体值不合法(不存在对应的主体)
	PermissionObjInvalid         ErrorCodeType = 100304 // 权限资源值不合法(不存在对应的角色)
	PermissionHasExisted         ErrorCodeType = 100305 // 权限已存在
)

// ... other module

var (
	Success = New(SuccessCode, "OK")

	ErrInvalidParams        = New(InvalidParamsErrCode, "Request params invalid.")
	ErrMissingParams        = New(MissingParamsErrCode, "Missing required params.")
	ErrDatabaseOperation    = New(DatabaseOperationFailed, "Database operation failed.")
	ErrParseRequestParams   = New(ParseRequestParamsError, "Parse request params error.")
	ErrSerializeDataFail    = New(SerializeDataFail, "Serialize data failed.")
	ErrDeserializeDataFail  = New(DeserializeDataFail, "Deserialize data failed.")
	ErrGetServerContextFail = New(GetServerContextFail, "Get server context failed.")
	ErrPermissionDenied     = New(PermissionDenied, "Permission not allowed.")
	ErrInternalServer       = New(InternalServerError, "Internal server error.")
	ErrBashShellRunFailed   = New(BashShellRunFailed, "Bash shell run failed.")

	ErrUserIdInvalid      = New(UserIdInvalid, "User id invalid.")
	ErrUserNotExistById   = New(UserNotExistById, "User not exist.")
	ErrUsernameHasExisted = New(UsernameHasExisted, "Username has existed.")
	ErrUsernameNotCorrect = New(UsernameNotCorrect, "Username not correct.")
	ErrPasswordNotCorrect = New(PasswordNotCorrect, "Password not correct.")
	ErrNotLoginStatus     = New(NotLoginStatus, "Not login yet.")
	ErrUserLogoutFail     = New(UserLogoutFail, "User logout fail.")

	ErrRoleIdInvalid    = New(RoleIdInvalid, "Role id invalid.")
	ErrRoleNotExistById = New(RoleNotExistById, "Role not exist.")

	ErrPermissionEnforcerLoadFailed = New(PermissionEnforcerLoadFailed, "Permission enforcer load failed.")
	ErrPermissionTypeInvalid        = New(PermissionTypeInvalid, "Permission type invalid.")
	ErrPermissionSubTypeInvalid     = New(PermissionSubTypeInvalid, "Permission sub type invalid.")
	ErrPermissionSubInvalid         = New(PermissionSubInvalid, "Permission sub invalid.")
	ErrPermissionObjInvalid         = New(PermissionObjInvalid, "Permission obj invalid.")
	ErrPermissionHasExisted         = New(PermissionHasExisted, "Permission has existed.")
)
