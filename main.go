package main

import (
	_ "GoApp/routers"
	"github.com/astaxie/beego"
	"GoApp/Utilities/helper"
	"github.com/goinggo/tracelog"
	"os"
)

const MainGoRoutine = "main"


func main() {
	beego.SetStaticPath("/static", "static")

	tracelog.Start(tracelog.LevelTrace)

	// Init mongo
	tracelog.Started("main", "Initializing Mongo")
	err := mongo.Startup(MainGoRoutine)
	if err != nil {
		tracelog.CompletedError(err, MainGoRoutine, "initApp")
		os.Exit(1)
	}

	beego.Run()

	tracelog.Completed(MainGoRoutine, "Website Shutdown")
	tracelog.Stop()
}

