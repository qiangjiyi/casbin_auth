// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/qiangjiyi/casbin_auth/controllers"
)

func init() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))

	ns := beego.NewNamespace("/v1/api",

		beego.NSNamespace("/auth",
			beego.NSNamespace("/user",
				beego.NSRouter("/register", &controllers.UserController{}, "post:Register"),
				beego.NSRouter("/login", &controllers.UserController{}, "post:Login"),
				beego.NSRouter("/logout", &controllers.UserController{}, "get:Logout"),
				beego.NSRouter("/update", &controllers.UserController{}, "post:UpdateUser"),
				beego.NSRouter("/delete/:userId", &controllers.UserController{}, "get:DeleteUser"),
			),

			beego.NSNamespace("/role",
				beego.NSRouter("/create", &controllers.RoleController{}, "post:CreateRole"),
				beego.NSRouter("/update", &controllers.RoleController{}, "post:UpdateRole"),
				beego.NSRouter("/delete/:roleId", &controllers.RoleController{}, "get:DeleteRole"),
			),

			beego.NSNamespace("/permission",
				beego.NSRouter("/query/:userId", &controllers.PermissionController{}, "get:QueryPermission"),
				beego.NSRouter("/add", &controllers.PermissionController{}, "post:AddPermission"),
			),

		),

	)
	beego.AddNamespace(ns)
}
