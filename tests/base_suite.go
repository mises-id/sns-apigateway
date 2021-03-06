//go:build tests
// +build tests

package tests

import (
	"context"
	"fmt"
	"net"
	"sort"
	"time"

	"github.com/khaiql/dbcleaner"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/services/session"
	_ "github.com/mises-id/sns-socialsvc/config"
	"github.com/mises-id/sns-socialsvc/config/env"
	"github.com/mises-id/sns-socialsvc/handlers"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/mises-id/sns-socialsvc/svc/server"

	/* storagesvcenv "github.com/mises-id/sns-storagesvc/config/env"
	storagehandler "github.com/mises-id/sns-storagesvc/handlers"
	storagesvcdb "github.com/mises-id/sns-storagesvc/lib/db"
	storagesvc "github.com/mises-id/sns-storagesvc/svc/server" */
	"github.com/stretchr/testify/suite"
)

func init() {
	fmt.Println("this is test build...")
}

type BaseTestSuite struct {
	suite.Suite
	dbcleaner.DbCleaner
}

func (suite *BaseTestSuite) SetupSuite() {
	suite.DbCleaner = dbcleaner.New()
	// TODO the env should read through api
	duration, _ := time.ParseDuration("24h")
	env.Envs = &env.Env{
		DBName:           "mises_unit_test",
		DBUser:           "",
		DBPass:           "",
		MongoURI:         "mongodb://localhost:27017",
		DebugMisesPrefix: "1001",
		TokenDuration:    duration,
	}
	db.SetupMongo(context.Background())
	models.EnsureIndex()
	session.SetupMisesClient()

	port := 15040
	cfg := server.DefaultConfig
	for {

		cfg.GRPCAddr = fmt.Sprintf(":%d", port)
		cfg.HTTPAddr = fmt.Sprintf(":%d", port+1)
		cfg.DebugAddr = fmt.Sprintf(":%d", port+2)
		ln, err := net.Listen("tcp", cfg.GRPCAddr)

		if err != nil {
			port += 10
		} else {
			_ = ln.Close()
			break
		}

	}
	/* storagesvcenv.Envs = &storagesvcenv.Env{
		DBName:         "mises_unit_test",
		DBUser:         "",
		DBPass:         "",
		MongoURI:       "mongodb://localhost:27017",
		AWSAccessKeyId: "1",
		AWSSecretKey:   "2",
		AWSRegion:      "3",
	}
	storagesvcdb.SetupMongo(context.Background())

	sport := 16040
	scfg := storagesvc.DefaultConfig
	for {

		scfg.GRPCAddr = fmt.Sprintf(":%d", sport)
		scfg.HTTPAddr = fmt.Sprintf(":%d", sport+1)
		//cfg.DebugAddr = fmt.Sprintf(":%d", port+2)
		ln, err := net.Listen("tcp", scfg.GRPCAddr)

		if err != nil {
			sport += 10
		} else {
			_ = ln.Close()
			break
		}

	} */
	rest.ResetSvrPool(rest.PoolCfg{cfg.GRPCAddr, ":6050", 1, 60})
	/* go func() {

		scfg = storagehandler.SetConfig(scfg)

		storagesvc.Run(scfg)

	}() */

	go func() {

		cfg = handlers.SetConfig(cfg)

		server.Run(cfg)

	}()

}

func (suite *BaseTestSuite) TearDownSuite() {
}

func (suite *BaseTestSuite) Clean(collections ...string) {
	sort.Strings(collections)
	suite.DbCleaner.Acquire(collections...)
	for _, collection := range collections {
		_ = db.DB().Collection(collection).Drop(context.Background())
	}
	models.EnsureIndex()
}

func (suite *BaseTestSuite) Acquire(collections ...string) {
	sort.Strings(collections)
	suite.DbCleaner.Acquire(collections...)
}
