package services

import (
	"gopkg.in/mgo.v2"
	"GoApp/Utilities/helper"
	log "github.com/goinggo/tracelog"
)

const MainGoRoutine = "main"


type (
	Service struct {
		MongoSession *mgo.Session
		UserID       string
	}
)

func (service *Service) Prepare() (err error) {

	service.MongoSession, err = mongo.CopyMonotonicSession(MainGoRoutine)
	if err != nil {
		log.Error(err, MainGoRoutine, "Service.Prepare")
		return err
	}

	return err
}

func (service *Service) Finish() (err error) {

	if service.MongoSession != nil {
		mongo.CloseSession(MainGoRoutine, service.MongoSession)
		service.MongoSession = nil
	}

	return err
}

func (service *Service) DBAction(databaseName string, collectionName string, dbCall mongo.DBCall) (err error) {
	return mongo.Execute(MainGoRoutine, service.MongoSession, databaseName, collectionName, dbCall)
}