package baseController

import (
	"github.com/astaxie/beego"
	"GoApp/services"
	"GoApp/Utilities/helper"
)

type (
	BaseController struct {
		beego.Controller
		services.Service
	}
)

func (baseController *BaseController) Prepare() {
	baseController.UserID = baseController.GetString("main")
	if baseController.UserID == "" {
		baseController.UserID = baseController.GetString(":main")
	}
	if baseController.UserID == "" {
		baseController.UserID = "Unknown"
	}

	if err := baseController.Service.Prepare(); err != nil {
		return
	}
}

func (baseController *BaseController) Finish() {
	defer func() {
		if baseController.MongoSession != nil {
			mongo.CloseSession("Main", baseController.MongoSession)
			baseController.MongoSession = nil
		}
	}()
}


