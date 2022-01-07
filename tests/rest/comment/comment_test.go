// +build tests

package comment

import (
	"net/http"
	"testing"

	"github.com/mises-id/sns-apigateway/tests/factories"
	"github.com/mises-id/sns-apigateway/tests/rest"
	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/stretchr/testify/suite"
)

type CommentServerSuite struct {
	rest.RestBaseTestSuite
	collections []string
	statuses    []*models.Status
}

func (suite *CommentServerSuite) SetupSuite() {
	suite.RestBaseTestSuite.SetupSuite()
	suite.collections = []string{"counters", "users", "follows"}
}

func (suite *CommentServerSuite) TearDownSuite() {
	suite.RestBaseTestSuite.TearDownSuite()
}

func (suite *CommentServerSuite) SetupTest() {
	suite.Clean(suite.collections...)
	suite.Acquire(suite.collections...)
	factories.InitUsers(&models.User{
		UID:     uint64(1001),
		Misesid: "1001",
	}, &models.User{
		UID:     uint64(1002),
		Misesid: "1002",
	})
	suite.statuses = factories.InitDefaultStatuses()
}

func (suite *CommentServerSuite) TearDownTest() {
	suite.Clean(suite.collections...)
}

func TestFollowServer(t *testing.T) {
	suite.Run(t, &CommentServerSuite{})
}

func (suite *CommentServerSuite) TestListComment() {
	token := suite.MockLoginUser("1001:1001")

	suite.T().Run("list comments empty", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/comment").WithQuery("status_id", suite.statuses[0].ID.Hex()).
			Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Array().Length().Equal(0)
		resp.Value("pagination").Object().Value("last_id").Equal("")
	})

	suite.T().Run("add comment", func(t *testing.T) {
		resp := suite.Expect.POST("/api/v1/comment").WithJSON(map[string]interface{}{
			"status_id": suite.statuses[0].ID.Hex(),
			"parent_id": "",
			"content":   "comment a  status",
		}).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
	})

	suite.T().Run("list comments ", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/comment").WithQuery("status_id", suite.statuses[0].ID.Hex()).
			Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Array().Length().Equal(1)
		resp.Value("pagination").Object().Value("last_id").Equal("")
	})

}
