package main

import (
	"net/http"

	//"github.com/gorilla/handlers"
	"os"

	"github.com/gkontos/gasket/acedb"
	log "github.com/gkontos/gasket/acelog"
	"github.com/gkontos/gasket/aceservice"
	"github.com/gkontos/gasket/aceweb"
	"github.com/spf13/viper"
)

type controller interface {
	registerRoute()
}

func main() {

	v := viper.New()
	v.SetConfigName("config") // no need to include file extension

	if _, err := os.Stat(os.Getenv("ACES_CFG")); err == nil {
		v.AddConfigPath(os.Getenv("ACES_CFG"))
	} else {
		v.AddConfigPath("./config") // set the path of your config file
	}
	err := v.ReadInConfig()

	if err != nil {
		log.Fatal("Unable to get configuration - ", err)
	}

	aceweb.SetVersion(v.GetString("app.version"))
	log.Info("Starting Server ", v.GetString("app.version"))

	ds := dbhandle.New()
	ds.SetConfig(v)
	graphStore, dberr := ds.GetStore()
	if dberr != nil {
		log.Fatal("Unable to get datastore connection - ", dberr)
	}

	aceservice.SetStore(graphStore)

	router := aceweb.SysViewRouter()

	log.Fatal(http.ListenAndServe(":8080", log.RequestLogHandler(router)))

}
