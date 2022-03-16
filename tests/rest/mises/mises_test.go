//go:build tests
// +build tests

package mises

import (
	"net/http"
	"testing"

	"github.com/mises-id/sns-apigateway/tests/rest"
	"github.com/stretchr/testify/suite"
)

type MisesServerSuite struct {
	rest.RestBaseTestSuite
}

func (suite *MisesServerSuite) SetupSuite() {
	suite.RestBaseTestSuite.SetupSuite()
}

func (suite *MisesServerSuite) TearDownSuite() {
	suite.RestBaseTestSuite.TearDownSuite()
}

func (suite *MisesServerSuite) SetupTest() {

}

func (suite *MisesServerSuite) TearDownTest() {
}

func TestMisesServer(t *testing.T) {
	suite.Run(t, &MisesServerSuite{})
}

func (suite *MisesServerSuite) TestGasPrice() {

	suite.T().Run("get gas price", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/mises/gasprices").
			Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
	})

}

func (suite *MisesServerSuite) TestChainInfo() {

	suite.T().Run("get  chaininfo", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/mises/chaininfo").
			Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		hash := resp.Value("data").Object().Value("block_hash").String().Raw()

		respCached := suite.Expect.GET("/api/v1/mises/chaininfo").
			Expect().Status(http.StatusOK).JSON().Object()
		respCached.Value("code").Equal(0)
		respCached.Value("data").Object().Value("block_hash").Equal(hash)

	})

}
