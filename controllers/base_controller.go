package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/validation"
	"github.com/casbin/casbin"
	"github.com/qiangjiyi/casbin_auth/common"
	"github.com/qiangjiyi/casbin_auth/common/errors"
	"github.com/qiangjiyi/casbin_auth/common/util"
	"reflect"
	"runtime"
)

// BaseController contains basic property and function
type BaseController struct {
	beego.Controller
	Enforcer   *casbin.Enforcer
	CurrUserId int
	SessionOn  bool
}

// Prepare is called prior to the base controller method.
func (baseController *BaseController) Prepare() {
	baseController.SessionOn = beego.BConfig.WebConfig.Session.SessionOn

	// get casbin.Enforcer from context
	e, ok := baseController.Ctx.Input.GetData("enforcer").(*casbin.Enforcer)
	if !ok {
		logs.Error("Get casbin enforcer failed.")
		baseController.Response(nil, errors.ErrPermissionEnforcerLoadFailed)
		return
	}
	baseController.Enforcer = e

	// get current user id from context
	baseController.CurrUserId, _ = baseController.Ctx.Input.GetData(util.SessionKey).(int)

	logs.Info("BaseController.Prepare -> Path[%s]", baseController.Ctx.Request.URL.Path)
}

// Finish is called once the base controller method completes.
func (baseController *BaseController) Finish() {
	logs.Info("BaseController.Finish -> Path[%s]", baseController.Ctx.Request.URL.Path)
}

// ** VALIDATION

// ParseAndValidate will run the params through the validation framework and then
// response with the specified localized or provided message.
func (baseController *BaseController) ParseAndValidate(params interface{}) bool {
	// This is not working anymore :(
	if err := baseController.ParseForm(params); err != nil {
		baseController.ServeError(err)
		return false
	}

	var valid validation.Validation
	ok, err := valid.Valid(params)
	if err != nil {
		baseController.ServeError(err)
		return false
	}

	if ok == false {
		// Build a map of the Error messages for each field
		messages2 := make(map[string]string)

		val := reflect.ValueOf(params).Elem()
		for i := 0; i < val.NumField(); i++ {
			// Look for an Error tag in the field
			typeField := val.Type().Field(i)
			tag := typeField.Tag
			tagValue := tag.Get("Error")

			// Was there an Error tag
			if tagValue != "" {
				messages2[typeField.Name] = tagValue
			}
		}

		// Build the Error response
		var errors []string
		/*for _, err := range valid.Errors {
			// Match an Error from the validation framework Errors
			// to a field name we have a mapping for
			message, ok := messages2[err.Field]
			if ok == true {
				// Use a localized message if one exists
				errors = append(errors, localize.T(message))
				continue
			}

			// No match, so use the message as is
			errors = append(errors, err.Message)
		}*/

		baseController.ServeValidationErrors(errors)
		return false
	}

	return true
}

// ** EXCEPTIONS

// ServeError prepares and serves an Error exception.
func (baseController *BaseController) ServeError(err error) {
	baseController.Data["json"] = struct {
		Error string `json:"Error"`
	}{err.Error()}
	baseController.Ctx.Output.SetStatus(500)
	baseController.ServeJSON()
}

// ServeValidationErrors prepares and serves a validation exception.
func (baseController *BaseController) ServeValidationErrors(Errors []string) {
	baseController.Data["json"] = struct {
		Errors []string `json:"Errors"`
	}{Errors}
	baseController.Ctx.Output.SetStatus(409)
	baseController.ServeJSON()
}

// ** CATCHING PANICS

// CatchPanic is used to catch any Panic and flogging exceptions. Returns a 500 as the response.
func (baseController *BaseController) CatchPanic(functionName string) {
	if r := recover(); r != nil {
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)

		baseController.ServeError(fmt.Errorf("%v", r))
	}
}

// Response returns a standard restful response.
func (baseController *BaseController) Response(payload interface{}, err errors.AuthError) {
	baseController.Data["json"] = common.Response{
		ErrCode:    err.Code(),
		ErrMessage: err.Error(),
		Payload:    payload,
	}
	baseController.ServeJSON()
}
