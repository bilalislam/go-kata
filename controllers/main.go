package controllers

import (
	bc "GoApp/controllers/baseController"
	"GoApp/services/userservice"
	"fmt"
)

type IndexController struct {
	bc.BaseController
}

func (c *IndexController) Index() {
	c.TplName = "main.html"
}

func (c *IndexController) GetAllUsers() {
	c.TplName = "main.html"
	results, err := userService.GetAllUsers(c.Service)

	if err != nil {
		fmt.Print("sorun var ")
	}

	c.Data["json"] = &results
	c.ServeJSON()
}
