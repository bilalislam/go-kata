package mongo

import (
	"gopkg.in/mgo.v2"
	"github.com/kelseyhightower/envconfig"
	log "github.com/goinggo/tracelog"
	"strings"
	"time"
	"fmt"
)

const (
	MasterSession = "master"
	MonotonicSession = "monotonic"
)

var (
	singleton mongoManager
)

type(
	mongoConfiguration struct {
		Hosts    string
		Database string
		UserName string
		Password string
	}

	mongoSession struct {
		mongoDBDialInfo *mgo.DialInfo
		mongoSession    *mgo.Session
	}

	mongoManager struct {
		sessions map[string]mongoSession
	}

	DBCall func(*mgo.Collection) error
)

func Startup(sessionID string) error {

	if singleton.sessions != nil {
		return nil
	}

	log.Started(sessionID, "Startup")

	var config mongoConfiguration
	config.Hosts = "localhost"
	config.Database = "MyAppDatabase"

	if err := envconfig.Process("mgo", &config); err != nil {
		log.CompletedError(err, sessionID, "Startup")
	}

	//Create the Mongo Manager
	singleton = mongoManager{
		sessions:make(map[string]mongoSession),
	}

	// Log the mongodb connection straps.
	log.Trace(sessionID, "Startup", "MongoDB : Hosts[%s]", config.Hosts)
	log.Trace(sessionID, "Startup", "MongoDB : Database[%s]", config.Database)
	log.Trace(sessionID, "Startup", "MongoDB : Username[%s]", config.UserName)

	hosts := strings.Split(config.Hosts, ",")

	// Create the strong session.
	if err := CreateSession(sessionID, "strong", MasterSession, hosts, config.Database, config.UserName, config.Password); err != nil {
		log.CompletedError(err, sessionID, "Startup")
		return err
	}

	// Create the monotonic session.
	if err := CreateSession(sessionID, "monotonic", MonotonicSession, hosts, config.Database, config.UserName, config.Password); err != nil {
		log.CompletedError(err, sessionID, "Startup")
		return err
	}

	log.Completed(sessionID, "Startup")

	return nil
}


// CreateSession creates a connection pool for use.
func CreateSession(sessionID string, mode string, sessionName string, hosts []string, databaseName string, username string, password string) error {

	log.Startedf(sessionID, "CreateSession", "Mode[%s] SessionName[%s] Hosts[%s] DatabaseName[%s] Username[%s]", mode, sessionName, hosts, databaseName, username)

	//Create the database object
	mongoSession := mongoSession{
		mongoDBDialInfo:&mgo.DialInfo{
			Addrs:hosts,
			Timeout:60 * time.Second,
			Database: databaseName,
			Username: username,
			Password: password,
		},
	}

	//Established the master session
	var err error
	mongoSession.mongoSession, err = mgo.DialWithInfo(mongoSession.mongoDBDialInfo)
	if err != nil {
		log.CompletedError(err, sessionID, "CreateSession")
		return err
	}

	switch mode {
	case "strong":
		mongoSession.mongoSession.SetMode(mgo.Strong, true)
		break
	case "monotonic":
		mongoSession.mongoSession.SetMode(mgo.Monotonic, true)
		break
	}

	mongoSession.mongoSession.SetSafe(&mgo.Safe{})
	// Add the database to the map.
	singleton.sessions[sessionName] = mongoSession
	log.Completed(sessionID, "CreateSession")

	return nil
}

// CopyMonotonicSession makes a copy of the monotonic session for client use.
func CopyMonotonicSession(sessionID string) (*mgo.Session, error) {
	return CopySession(sessionID, MonotonicSession)
}

// CopySession makes a copy of the specified session for client use.
func CopySession(sessionID string, useSession string) (*mgo.Session, error) {
	log.Startedf(sessionID, "CopySession", "UseSession[%s]", useSession)

	// Find the session object.
	session := singleton.sessions[useSession]

	if session.mongoSession == nil {
		err := fmt.Errorf("Unable To Locate Session %s", useSession)
		log.CompletedError(err, sessionID, "CopySession")
		return nil, err
	}

	// Copy the master session.
	mongoSession := session.mongoSession.Copy()

	log.Completed(sessionID, "CopySession")
	return mongoSession, nil
}

// Shutdown systematically brings the manager down gracefully.
func Shutdown(sessionID string) error {
	log.Started(sessionID, "Shutdown")

	// Close the databases
	for _, session := range singleton.sessions {
		CloseSession(sessionID, session.mongoSession)
	}

	log.Completed(sessionID, "Shutdown")
	return nil
}

// CloseSession puts the connection back into the pool.
func CloseSession(sessionID string, mongoSession *mgo.Session) {
	log.Started(sessionID, "CloseSession")
	mongoSession.Close()
	log.Completed(sessionID, "CloseSession")
}

// GetDatabase returns a reference to the specified database.
func GetDatabase(mongoSession *mgo.Session, useDatabase string) *mgo.Database {
	return mongoSession.DB(useDatabase)
}

// GetCollection returns a reference to a collection for the specified database and collection name.
func GetCollection(mongoSession *mgo.Session, useDatabase string, useCollection string) *mgo.Collection {
	return mongoSession.DB(useDatabase).C(useCollection)
}

// Execute the MongoDB literal function.
func Execute(sessionID string, mongoSession *mgo.Session, databaseName string, collectionName string, dbCall DBCall) error {
	log.Startedf(sessionID, "Execute", "Database[%s] Collection[%s]", databaseName, collectionName)

	// Capture the specified collection.
	collection := GetCollection(mongoSession, databaseName, collectionName)
	if collection == nil {
		err := fmt.Errorf("Collection %s does not exist", collectionName)
		log.CompletedError(err, sessionID, "Execute")
		return err
	}

	// Execute the MongoDB call.
	err := dbCall(collection)
	if err != nil {
		log.CompletedError(err, sessionID, "Execute")
		return err
	}

	log.Completed(sessionID, "Execute")
	return nil
}