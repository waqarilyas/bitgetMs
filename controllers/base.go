package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
	"github.com/kryptomind/bidboxapi/bitgetms/models"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

type AppData struct {
	Keys_list  []models.Key
	Conditions []models.Conditions
}

var app_data AppData

func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {

	var err error
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	server.DB, err = gorm.Open(Dbdriver, DBURL)
	if err != nil {
		log.Info("Cannot connect to the %s database", Dbdriver)
		log.Fatal("This is the error:", err)
	} else {
		log.Info("Connected to the %s database", Dbdriver)
	}

	// server.DB.Debug().AutoMigrate(&models.Key{}) //database migration
	server.Router = mux.NewRouter()
	server.initializeRoutes()
	server.InitKeys()

	//server.InitConditions()
}

func (server *Server) InitKeys() {
	key := models.Key{}
	keys, err := key.FindAllKeys(server.DB)
	if err != nil {
		log.Fatal("error getting keys")
		app_data.Keys_list = []models.Key{}
		return
	}
	log.Info("retreived keys")
	app_data.Keys_list = *keys
}

func (server *Server) GetExchangeSpecificKeys(service string) []models.Key {
	key := models.Key{}
	keys, err := key.FindKeysByService(server.DB, service)
	if err != nil {
		log.Fatal("error getting keys")
		app_data.Keys_list = []models.Key{}
		return []models.Key{}
	}

	log.Info("retrieved keys")
	return *keys
}

func (server *Server) InitConditions() {
	cond := models.Conditions{}
	conds, err := cond.FindAllConditions(server.DB)
	if err != nil {
		log.Fatal("error getting keys")
		app_data.Conditions = []models.Conditions{}
		return
	}
	log.Info("retreived conditions")
	app_data.Conditions = *conds
}

func (server *Server) Run(addr string) {
	log.Info("Listening on port 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))

}
