package routers

import (
	"GoApp/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("Index", new(controllers.IndexController), "get:Index")
	beego.Router("GetAllUsers", new(controllers.IndexController), "get:GetAllUsers")
}
