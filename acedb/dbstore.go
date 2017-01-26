package dbhandle

import (
	"fmt"
	"sync"

	log "github.com/gkontos/gasket/acelog"
	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/mongo"
	"github.com/spf13/viper"
)

type Store interface {
	GetStore() (*cayley.Handle, error)
	SetConfig(conf *viper.Viper)
}

type graphStore struct {
	dbstore *cayley.Handle
	config  *viper.Viper
}

var (
	once sync.Once
)

func New() *graphStore {

	gs := &graphStore{}
	return gs
}

func (repo *graphStore) SetConfig(conf *viper.Viper) {
	repo.config = conf
}
func (repo *graphStore) GetStore() (*cayley.Handle, error) {
	if repo.config == nil {
		return nil, fmt.Errorf("Configuration for datastore not set")
	}
	var err error
	once.Do(func() {
		// repo = &graphStore{}
		err = repo.configure()
	})
	return repo.dbstore, err
}

// config will return a cayley.Handle for a mongo database
func (repo *graphStore) configure() error {

	if repo.dbstore == nil {

		server := repo.config.GetString("db.server")
		port := repo.config.GetString("db.port")
		opts := make(graph.Options)
		//		opts["username"] = ""
		//		opts["password"] = ""
		//		opts["database_name"] = "cayley"

		// Create a brand new graph
		addr := server + ":" + port
		graph.InitQuadStore("mongo", addr, opts)

		var err error
		repo.dbstore, err = cayley.NewGraph("mongo", addr, opts)
		if err != nil {
			log.Fatal(err)
			return err
		}

	}
	return nil
}
