package tests

import (
	"context"
	"sort"

	"github.com/khaiql/dbcleaner"
	"github.com/mises-id/socialsvc/app/models"
	"github.com/mises-id/socialsvc/config/env"
	"github.com/mises-id/socialsvc/lib/db"
	"github.com/stretchr/testify/suite"
)

type BaseTestSuite struct {
	suite.Suite
	dbcleaner.DbCleaner
}

func (suite *BaseTestSuite) SetupSuite() {
	suite.DbCleaner = dbcleaner.New()
	// TODO the env should read through api
	env.Envs = &env.Env{
		DBName:   "mises_dev",
		DBUser:   "root",
		DBPass:   "example",
		MongoURI: "mongodb://localhost:27017",
	}
	db.SetupMongo(context.Background())
	models.EnsureIndex()
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
