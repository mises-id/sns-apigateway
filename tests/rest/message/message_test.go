// +build tests

package message

import (
	"net/http"
	"testing"

	"github.com/mises-id/sns-apigateway/tests/factories"
	"github.com/mises-id/sns-apigateway/tests/rest"
	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/stretchr/testify/suite"
)

type MessageServerSuite struct {
	rest.RestBaseTestSuite
	collections []string
	statuses    []*models.Status
}

func (suite *MessageServerSuite) SetupSuite() {
	suite.RestBaseTestSuite.SetupSuite()
	suite.collections = []string{"counters", "users", "follows", "messages"}
}

func (suite *MessageServerSuite) TearDownSuite() {
	suite.RestBaseTestSuite.TearDownSuite()
}

func (suite *MessageServerSuite) SetupTest() {
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

func (suite *MessageServerSuite) TearDownTest() {
	suite.Clean(suite.collections...)
}

func TestFollowServer(t *testing.T) {
	suite.Run(t, &MessageServerSuite{})
}

func (suite *MessageServerSuite) TestNewCommentMessage() {
	token1001 := suite.MockLoginUser("1001:1001")
	token1002 := suite.MockLoginUser("1002:1002")

	suite.T().Run("summary message empty", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/user/message/summary").WithHeader("Authorization", "Bearer "+token1001).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Object().Value("total").Equal(0)
		resp.Value("data").Object().Value("notifications_count").Equal(0)
		resp.Value("data").Object().Value("users_count").Equal(0)

	})
	suite.T().Run("list message empty", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/user/message").WithQuery("limit", 3).WithHeader("Authorization", "Bearer "+token1001).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Array().Length().Equal(0)
		resp.Value("pagination").Object().Value("last_id").Equal("")
	})

	suite.T().Run("add comment", func(t *testing.T) {
		resp := suite.Expect.POST("/api/v1/comment").WithJSON(map[string]interface{}{
			"status_id": suite.statuses[0].ID.Hex(),
			"parent_id": "",
			"content":   "comment a  status",
		}).WithHeader("Authorization", "Bearer "+token1002).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
	})

	suite.T().Run("summary message new comment", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/user/message/summary").WithHeader("Authorization", "Bearer "+token1001).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Object().Value("total").Equal(1)
		resp.Value("data").Object().Value("notifications_count").Equal(1)
		resp.Value("data").Object().Value("users_count").Equal(0)

	})

	var pending_message_id string
	suite.T().Run("list message new comment", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/user/message").WithQuery("limit", 3).WithHeader("Authorization", "Bearer "+token1001).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Array().Length().Equal(1)
		resp.Value("data").Array().First().Object().Value("message_type").Equal("new_comment")
		resp.Value("data").Array().First().Object().Value("state").Equal("unread")
		pending_message_id = resp.Value("data").Array().First().Object().Value("id").String().Raw()
		resp.Value("pagination").Object().Value("last_id").Equal("")

	})

	suite.T().Run("read message", func(t *testing.T) {
		resp := suite.Expect.PUT("/api/v1/user/message/read").WithJSON(map[string]interface{}{
			"ids": []string{pending_message_id},
		}).WithHeader("Authorization", "Bearer "+token1001).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
	})

	suite.T().Run("summary message after read", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/user/message/summary").WithHeader("Authorization", "Bearer "+token1001).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Object().Value("total").Equal(0)
		resp.Value("data").Object().Value("notifications_count").Equal(0)
		resp.Value("data").Object().Value("users_count").Equal(0)

	})

	suite.T().Run("list message after read", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/user/message").WithQuery("limit", 3).WithHeader("Authorization", "Bearer "+token1001).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Array().Length().Equal(1)
		resp.Value("data").Array().First().Object().Value("message_type").Equal("new_comment")
		resp.Value("data").Array().First().Object().Value("state").Equal("read")
		resp.Value("data").Array().First().Object().Value("id").Equal(pending_message_id)
		resp.Value("pagination").Object().Value("last_id").Equal("")

	})

}
