package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/1046102779/grbac/controllers:FunctionsController"] = append(beego.GlobalControllerRouter["github.com/1046102779/grbac/controllers:FunctionsController"],
		beego.ControllerComments{
			Method: "GetFuncId",
			Router: `/`,
			AllowHTTPMethods: []string{"POST"},
			Params: nil})

}
