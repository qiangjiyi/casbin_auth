package filter

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/casbin/casbin"
	"github.com/qiangjiyi/casbin_auth/cache"
	"github.com/qiangjiyi/casbin_auth/common"
	"github.com/qiangjiyi/casbin_auth/common/errors"
	"github.com/qiangjiyi/casbin_auth/common/util"
	"net/http"
	"strconv"
	"time"
)

// NewAuthorizer returns the authorizer.
// Use a casbin enforcer as input
func NewAuthorizer(e *casbin.Enforcer, sessionOn bool) beego.FilterFunc {
	return func(ctx *context.Context) {
		if ctx.Request.RequestURI == "/v1/api/auth/user/register" || ctx.Request.RequestURI == "/v1/api/auth/user/login" || ctx.Request.RequestURI == "/v1/api/auth/user/logout" { // if request uri is register/login/logout then no need to enforce
			ctx.Input.SetData("enforcer", e)
			return
		}

		a := &BasicAuthorizer{enforcer: e}

		if sessionOn { // if beego's session opened
			if ctx.Input.Session(util.SessionKey) == nil { // not login
				a.RequireLogin(ctx)
				return
			}

			a.uid, _ = ctx.Input.Session(util.SessionKey).(int) // get user id from session
		} else { // if beego's session closed then use token mode
			redis := cache.RedisCacheForSession()
			if redis == nil { // get redis cache failed
				a.InternalServerError(ctx)
				return
			}

			if value, ok := redis.Get(ctx.Input.Header(util.TokenCacheKey)).([]byte); ok {
				redis.Put(ctx.Input.Header(util.TokenCacheKey), value, 30*time.Minute) // refresh token expire time
				a.uid, _ = strconv.Atoi(string(value))
			} else {
				a.RequireLogin(ctx)
				return
			}
		}

		// query user's all roles by user id from casbin_rule
		// a.roles = e.GetRolesForUser(strconv.Itoa(a.uid))

		if !a.CheckPermission(ctx.Request) {
			a.RequirePermission(ctx)
		} else { // save casbin.Enforcer to context
			ctx.Input.SetData("enforcer", e)
			ctx.Input.SetData(util.SessionKey, a.uid)
		}
	}
}

// BasicAuthorizer stores the casbin handler
type BasicAuthorizer struct {
	uid int
	// roles    []string
	enforcer *casbin.Enforcer
}

// CheckPermission checks the user/method/path combination from the request.
// Returns true (permission granted) or false (permission forbidden)
func (a *BasicAuthorizer) CheckPermission(r *http.Request) bool {
	method := r.Method
	path := r.URL.Path

	ok := a.enforcer.Enforce(strconv.Itoa(a.uid), path, method) // use casbin enforcer to enforce if user has permission
	return ok
}

// RequirePermission returns the 403 Forbidden to the client
func (a *BasicAuthorizer) RequirePermission(ctx *context.Context) {
	ctx.ResponseWriter.WriteHeader(403)
	ctx.Output.JSON(common.Response{
		ErrCode:    errors.ErrPermissionDenied.Code(),
		ErrMessage: errors.ErrPermissionDenied.Error(),
	}, true, false)
}

// RequireLogin returns the 302 Redirect to the client
func (a *BasicAuthorizer) RequireLogin(ctx *context.Context) {
	ctx.ResponseWriter.WriteHeader(302)
	ctx.Output.JSON(common.Response{
		ErrCode:    errors.ErrNotLoginStatus.Code(),
		ErrMessage: errors.ErrNotLoginStatus.Error(),
	}, true, false)
}

// InternalServerError returns the 500 to the client
func (a *BasicAuthorizer) InternalServerError(ctx *context.Context) {
	ctx.ResponseWriter.WriteHeader(500)
	ctx.Output.JSON(common.Response{
		ErrCode:    errors.ErrInternalServer.Code(),
		ErrMessage: errors.ErrInternalServer.Error(),
	}, true, false)
}
