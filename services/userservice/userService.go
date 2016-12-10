package userService

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/goinggo/tracelog"
	"gopkg.in/mgo.v2"
	"GoApp/services"
	"GoApp/models"
)

const MainGoRoutine = "main"

type (
	serviceConfiguration struct {
		Database string
	}
)

var Config serviceConfiguration

func init() {
	Config.Database = "Database"

	if err := envconfig.Process("Database", &Config); err != nil {
		log.CompletedError(err, MainGoRoutine, "Init")
	}

}

func GetAllUsers(service services.Service) (*[]models.User, error) {

	var results []models.User

	f := func(collection *mgo.Collection) error {
		return collection.Find(nil).All(&results)
	}

	if err := service.DBAction(Config.Database, "users", f); err != nil {
		if err != mgo.ErrNotFound {
			return nil, err
		}
	}

	return &results, nil
}