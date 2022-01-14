// +build tests

package status

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/mises-id/sns-apigateway/lib/codes"
	"github.com/mises-id/sns-apigateway/tests/factories"
	"github.com/mises-id/sns-apigateway/tests/rest"
	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StatusServerSuite struct {
	rest.RestBaseTestSuite
	collections []string
	statuses    []*models.Status
}

func (suite *StatusServerSuite) SetupSuite() {
	suite.RestBaseTestSuite.SetupSuite()
	suite.collections = []string{"counters", "users", "follows", "statuses", "likes"}
}

func (suite *StatusServerSuite) TearDownSuite() {
	suite.RestBaseTestSuite.TearDownSuite()
}

func (suite *StatusServerSuite) SetupTest() {
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

func (suite *StatusServerSuite) TearDownTest() {
	suite.Clean(suite.collections...)
}

func TestStatusServer(t *testing.T) {
	suite.Run(t, &StatusServerSuite{})
}

func (suite *StatusServerSuite) TestListStatus() {
	token := suite.MockLoginUser("1001:1001")

	suite.T().Run("list user status", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/user/1001/status").Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("data").Array()
	})

	suite.T().Run("user timeline", func(t *testing.T) {
		suite.Expect.POST("/api/v1/user/follow").WithJSON(map[string]interface{}{"to_user_id": 1002}).
			WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK)

		resp := suite.Expect.GET("/api/v1/timeline/me").
			WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("data").Array()
	})
}

func (suite *StatusServerSuite) TestCreateStatus() {
	token := suite.MockLoginUser("1001:1001")
	linkMeta := &map[string]interface{}{
		"title":           "Test link title",
		"host":            "www.test.com",
		"attachment_path": "1",
		"link":            "http://www.test.com/articles/test/1",
	}
	suite.T().Run("create a text status", func(t *testing.T) {
		resp := suite.Expect.POST("/api/v1/status").WithJSON(map[string]interface{}{
			"status_type": "text",
			"content":     "post a text status",
		}).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(codes.SuccessCode)
		status := &models.Status{}
		err := db.ODM(context.Background()).Last(status).Error
		suite.Nil(err)
		suite.Equal("post a text status", status.Content)
		suite.Equal(enum.TextStatus, status.StatusType)
		suite.Equal(uint64(1001), status.UID)
	})
	suite.T().Run("create a link status", func(t *testing.T) {
		resp := suite.Expect.POST("/api/v1/status").WithJSON(map[string]interface{}{
			"status_type": "link",
			"content":     "post a link status",
			"link_meta":   linkMeta,
		}).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(codes.SuccessCode)
		status := &models.Status{}
		err := db.ODM(context.Background()).Last(status).Error
		suite.Nil(err)
		suite.Equal("post a link status", status.Content)
		suite.Equal(enum.LinkStatus, status.StatusType)
		suite.Equal("1", status.LinkMeta.ImagePath)
		suite.Equal(uint64(1001), status.UID)
	})
	suite.T().Run("forward a text status", func(t *testing.T) {
		resp := suite.Expect.POST("/api/v1/status").WithJSON(map[string]interface{}{
			"status_type": "text",
			"parent_id":   suite.statuses[0].ID.Hex(),
			"content":     "forward a text status",
		}).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(codes.SuccessCode)
		status := &models.Status{}
		err := db.ODM(context.Background()).Last(status).Error
		suite.Nil(err)
		suite.Equal("forward a text status", status.Content)
		suite.Equal(enum.TextStatus, status.StatusType)
		suite.Equal(suite.statuses[0].ID.Hex(), status.ParentID.Hex())
		suite.Equal(suite.statuses[0].ID.Hex(), status.OriginID.Hex())
		suite.Equal(uint64(1001), status.UID)

		parentStatus := &models.Status{}
		err = db.ODM(context.Background()).First(parentStatus, bson.M{"_id": suite.statuses[0].ID}).Error
		suite.Nil(err)
		suite.Equal(uint64(1), parentStatus.ForwardsCount)
	})
	suite.T().Run("forward a link status", func(t *testing.T) {
		resp := suite.Expect.POST("/api/v1/status").WithJSON(map[string]interface{}{
			"status_type": "text",
			"parent_id":   suite.statuses[1].ID.Hex(),
			"content":     "forward a link status",
		}).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(codes.SuccessCode)
		status := &models.Status{}
		err := db.ODM(context.Background()).Last(status).Error
		suite.Nil(err)
		suite.Equal("forward a link status", status.Content)
		suite.Equal(enum.TextStatus, status.StatusType)
		suite.Equal(suite.statuses[1].ID.Hex(), status.ParentID.Hex())
		suite.Equal(suite.statuses[1].ID.Hex(), status.OriginID.Hex())
		suite.Equal(uint64(1001), status.UID)

		parentStatus := &models.Status{}
		err = db.ODM(context.Background()).First(parentStatus, bson.M{"_id": suite.statuses[1].ID}).Error
		suite.Nil(err)
		suite.Equal(uint64(1), parentStatus.ForwardsCount)
	})
}

func (suite *StatusServerSuite) TestDeleteStatus() {
	token := suite.MockLoginUser("1001:1001")
	suite.T().Run("delete status not found", func(t *testing.T) {
		resp := suite.Expect.DELETE("/api/v1/status/xxxxxxx").
			WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusInternalServerError).JSON().Object()
		resp.Value("code").Equal(codes.InternalCode)

		resp = suite.Expect.DELETE("/api/v1/status/"+primitive.NewObjectID().Hex()).
			WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusNotFound).JSON().Object()
		resp.Value("code").Equal(codes.NotFoundCode)
	})

	suite.T().Run("delete status forbidden", func(t *testing.T) {
		resp := suite.Expect.DELETE("/api/v1/status/"+suite.statuses[1].ID.Hex()).
			WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusForbidden).JSON().Object()
		resp.Value("code").Equal(codes.ForbiddenCode)
	})

	suite.T().Run("delete status success", func(t *testing.T) {
		resp := suite.Expect.DELETE("/api/v1/status/"+suite.statuses[0].ID.Hex()).
			WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(codes.SuccessCode)
	})
}

func (suite *StatusServerSuite) TestLikeStatus() {
	token := suite.MockLoginUser("1001:1001")
	suite.T().Run("like a status", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/user/1001/like").
			WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Array().Length().Equal(0)

		resp = suite.Expect.POST(fmt.Sprintf("/api/v1/status/%s/like", suite.statuses[1].ID.Hex())).
			WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(codes.SuccessCode)

		resp = suite.Expect.POST(fmt.Sprintf("/api/v1/status/%s/like", suite.statuses[1].ID.Hex())).
			WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(codes.SuccessCode)

		likes := make([]*models.Like, 0)
		err := db.ODM(context.TODO()).Find(&likes).Error
		suite.Nil(err)
		suite.Equal(1, len(likes))

		status := &models.Status{}
		err = db.ODM(context.TODO()).First(status, bson.M{"_id": suite.statuses[1].ID}).Error
		suite.Nil(err)
		suite.Equal(uint64(1), status.LikesCount)

		resp = suite.Expect.GET("/api/v1/user/1001/like").
			WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Array().Length().Equal(1)
	})
}

func (suite *StatusServerSuite) TestUnlikeStatus() {
	token := suite.MockLoginUser("1001:1001")
	suite.T().Run("unlike a status", func(t *testing.T) {
		resp := suite.Expect.DELETE(fmt.Sprintf("/api/v1/status/%s/like", suite.statuses[0].ID.Hex())).
			WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusNotFound).JSON().Object()
		resp.Value("code").Equal(codes.NotFoundCode)

		resp = suite.Expect.POST(fmt.Sprintf("/api/v1/status/%s/like", suite.statuses[0].ID.Hex())).
			WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(codes.SuccessCode)

		resp = suite.Expect.DELETE(fmt.Sprintf("/api/v1/status/%s/like", suite.statuses[0].ID.Hex())).
			WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(codes.SuccessCode)

		likes := make([]*models.Like, 0)
		err := db.ODM(context.TODO()).Find(&likes).Error
		suite.Nil(err)
		suite.Equal(1, len(likes))
		suite.NotNil(likes[0].DeletedAt)

		resp = suite.Expect.GET("/api/v1/user/1001/like").
			WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Array().Length().Equal(0)
	})
}

func (suite *StatusServerSuite) TestGetUpdateStatus() {
	token := suite.MockLoginUser("1001:1001")
	var statusID string
	suite.T().Run("create a text status", func(t *testing.T) {
		resp := suite.Expect.POST("/api/v1/status").WithJSON(map[string]interface{}{
			"status_type": "text",
			"content":     "post a text status",
		}).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(codes.SuccessCode)
		statusID = resp.Value("data").Object().Value("id").String().Raw()
	})

	suite.T().Run("get a public status", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/status/"+statusID).
			WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Object().Value("id").Equal(statusID)
		resp.Value("data").Object().Value("is_public").Equal(true)
	})
	suite.T().Run("update a status", func(t *testing.T) {
		resp := suite.Expect.PATCH("/api/v1/status/"+statusID).WithJSON(map[string]interface{}{
			"is_private": true,
		}).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp = suite.Expect.GET("/api/v1/status/"+statusID).
			WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Object().Value("id").Equal(statusID)
		resp.Value("data").Object().Value("is_public").Equal(false)
	})

}

func (suite *StatusServerSuite) TestRecommendStatus() {
	token := suite.MockLoginUser("1001:1001")
	suite.T().Run("recommend status for guest", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/status/recommend").Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("data").Array()
	})

	suite.T().Run("recommend status pagination", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/status/recommend").WithQuery("limit", 3).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("data").Array()
		nextid := resp.Value("pagination").Object().Value("last_id").String().Raw()

		resp = suite.Expect.GET("/api/v1/status/recommend").
			WithQuery("limit", 1).WithQuery("last_id", nextid).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("data").Array()
	})

	suite.T().Run("recommend status for user", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/status/recommend").
			WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("data").Array()
	})
}
