// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/1046102779/grbac/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1/grbac",

		beego.NSNamespace("/entities",
			beego.NSInclude(
				&controllers.EntitiesController{},
			),
		),

		beego.NSNamespace("/functions",
			beego.NSInclude(
				&controllers.FunctionsController{},
			),
		),

		beego.NSNamespace("/regions",
			beego.NSInclude(
				&controllers.RegionsController{},
			),
		),

		beego.NSNamespace("/role_function_relationships",
			beego.NSInclude(
				&controllers.RoleFunctionRelationshipsController{},
			),
		),

		beego.NSNamespace("/roles",
			beego.NSInclude(
				&controllers.RolesController{},
			),
		),

		beego.NSNamespace("/user_roles",
			beego.NSInclude(
				&controllers.UserRolesController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
