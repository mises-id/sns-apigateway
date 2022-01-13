// +build tests

package follow

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/mises-id/sns-apigateway/tests/factories"
	"github.com/mises-id/sns-apigateway/tests/rest"
	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FollowServerSuite struct {
	rest.RestBaseTestSuite
	collections []string
}

func (suite *FollowServerSuite) SetupSuite() {
	suite.RestBaseTestSuite.SetupSuite()
	suite.collections = []string{"counters", "users", "follows"}
}

func (suite *FollowServerSuite) TearDownSuite() {
	suite.RestBaseTestSuite.TearDownSuite()
}

func (suite *FollowServerSuite) SetupTest() {
	suite.Clean(suite.collections...)
	suite.Acquire(suite.collections...)

}

func (suite *FollowServerSuite) TearDownTest() {
	suite.Clean(suite.collections...)
}

func TestFollowServer(t *testing.T) {
	suite.Run(t, &FollowServerSuite{})
}

func (suite *FollowServerSuite) TestListFriendship() {
	user1 := factories.UserFactory.MustCreate().(*models.User)
	users := make([]*models.User, 12)
	for i := range users {
		users[i] = factories.UserFactory.MustCreate().(*models.User)
		isFriend := i > 7
		if i <= 3 || i > 7 {
			factories.FollowFactory.MustCreateWithOption(map[string]interface{}{
				"FromUID":  user1.UID,
				"ToUID":    users[i].UID,
				"IsFriend": isFriend,
			})
		}
		if i > 3 {
			factories.FollowFactory.MustCreateWithOption(map[string]interface{}{
				"FromUID":  users[i].UID,
				"ToUID":    user1.UID,
				"IsFriend": isFriend,
			})
		}
	}

	suite.T().Run("not found user", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/user/999/friendship").
			Expect().Status(http.StatusNotFound).JSON().Object()
		resp.Value("code").Equal(404000)
	})

	suite.T().Run("list fans", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/user/"+fmt.Sprintf("%d", user1.UID)+"/friendship").WithQuery("relation_type", "fan").
			Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Array().Length().Equal(8)
		resp.Value("data").Array().First().Object().Value("relation_type").Equal("friend")
		resp.Value("data").Array().First().Object().Value("user").Object().Value("uid").Equal(users[11].UID)
		resp.Value("data").Array().Last().Object().Value("user").Object().Value("uid").Equal(users[4].UID)
		resp.Value("data").Array().Last().Object().Value("relation_type").Equal("fan")
		resp.Value("pagination").Object().Value("last_id").Equal("")
	})

	suite.T().Run("list following", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/user/"+fmt.Sprintf("%d", user1.UID)+"/friendship").WithQuery("relation_type", "following").
			Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Array().Length().Equal(8)
		resp.Value("data").Array().First().Object().Value("relation_type").Equal("friend")
		resp.Value("data").Array().First().Object().Value("user").Object().Value("uid").Equal(users[11].UID)
		resp.Value("data").Array().Last().Object().Value("user").Object().Value("uid").Equal(users[0].UID)
		resp.Value("data").Array().Last().Object().Value("relation_type").Equal("following")
	})

	suite.T().Run("list friend", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/user/"+fmt.Sprintf("%d", user1.UID)+"/friendship").WithQuery("relation_type", "friend").
			Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Array().Length().Equal(4)
		resp.Value("data").Array().Last().Object().Value("user").Object().Value("uid").Equal(users[8].UID)
	})

	suite.T().Run("list page", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/user/"+fmt.Sprintf("%d", user1.UID)+"/friendship").WithQuery("relation_type", "fan").
			WithQuery("limit", "3").
			Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Array().Length().Equal(3)
		resp.Value("pagination").Object().Value("limit").Equal(3)
		resp.Value("pagination").Object().Value("last_id").NotEqual("")
	})

}

func (suite *FollowServerSuite) TestFollowUnfollow() {

	user1 := factories.UserFactory.MustCreate().(*models.User)
	user2 := factories.UserFactory.MustCreate().(*models.User)
	user3 := factories.UserFactory.MustCreate().(*models.User)
	factories.FollowFactory.MustCreateWithOption(map[string]interface{}{
		"FromUID":  user3.UID,
		"ToUID":    user1.UID,
		"IsFriend": false,
	})
	token := suite.MockLoginUser("1001:" + user1.Misesid)
	println("token", token)
	suite.T().Run("follow stranger", func(t *testing.T) {
		resp := suite.Expect.POST("/api/v1/user/follow").WithJSON(map[string]interface{}{"to_user_id": user2.UID}).
			WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		f, err := models.GetFollow(context.Background(), user1.UID, user2.UID)
		suite.Nil(err)
		suite.False(f.IsFriend)
		_, err = models.GetFollow(context.Background(), user2.UID, user1.UID)
		suite.Equal(err, mongo.ErrNoDocuments)

		u1, u2 := &models.User{}, &models.User{}
		err = db.ODM(context.Background()).First(u1, bson.M{"_id": user1.UID}).Error
		suite.Nil(err)
		db.ODM(context.Background()).First(u2, bson.M{"_id": user2.UID})
		suite.Nil(err)
		suite.Equal(uint32(0), u1.FansCount)
		suite.Equal(uint32(1), u1.FollowingCount)
		suite.Equal(uint32(1), u2.FansCount)
		suite.Equal(uint32(0), u2.FollowingCount)
	})

	suite.T().Run("follow fans", func(t *testing.T) {
		resp := suite.Expect.POST("/api/v1/user/follow").WithJSON(map[string]interface{}{"to_user_id": user3.UID}).
			WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		f, err := models.GetFollow(context.Background(), user1.UID, user3.UID)
		suite.Nil(err)
		suite.True(f.IsFriend)
		f, err = models.GetFollow(context.Background(), user3.UID, user1.UID)
		suite.Nil(err)
		suite.True(f.IsFriend)
	})

	suite.T().Run("unfollow focus user", func(t *testing.T) {
		resp := suite.Expect.DELETE("/api/v1/user/follow").WithJSON(map[string]interface{}{"to_user_id": user2.UID}).
			WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		_, err := models.GetFollow(context.Background(), user1.UID, user2.UID)
		suite.Equal(err, mongo.ErrNoDocuments)
	})

	suite.T().Run("unfollow friend", func(t *testing.T) {
		resp := suite.Expect.DELETE("/api/v1/user/follow").WithJSON(map[string]interface{}{"to_user_id": user3.UID}).
			WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		_, err := models.GetFollow(context.Background(), user1.UID, user3.UID)
		suite.Equal(err, mongo.ErrNoDocuments)
		f, err := models.GetFollow(context.Background(), user3.UID, user1.UID)
		suite.Nil(err)
		suite.False(f.IsFriend)
	})
}
func (suite *FollowServerSuite) TestLatestFollowing() {

	user1 := factories.UserFactory.MustCreate().(*models.User)
	user2 := factories.UserFactory.MustCreate().(*models.User)
	token1 := suite.MockLoginUser("1:" + user1.Misesid)
	token2 := suite.MockLoginUser("2:" + user2.Misesid)
	suite.T().Run("list latest following", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/user/following/latest").
			WithHeader("Authorization", "Bearer "+token1).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Array().Length().Equal(0)
	})
	suite.T().Run("new fans", func(t *testing.T) {
		resp := suite.Expect.POST("/api/v1/user/follow").WithJSON(map[string]interface{}{"to_user_id": user1.UID}).
			WithHeader("Authorization", "Bearer "+token2).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp = suite.Expect.GET("/api/v1/user/"+fmt.Sprintf("%d", user1.UID)+"/friendship").WithQuery("relation_type", "fan").
			Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Array().Length().Equal(1)
	})

	suite.T().Run("list latest following empty", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/user/following/latest").
			WithHeader("Authorization", "Bearer "+token2).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Array().Length().Equal(1)
		resp.Value("data").Array().First().Object().Value("unread").Equal(false)
	})

	suite.T().Run("list latest following unread", func(t *testing.T) {
		resp := suite.Expect.POST("/api/v1/status").WithJSON(map[string]interface{}{
			"status_type": "text",
			"content":     "post a text status",
		}).WithHeader("Authorization", "Bearer "+token1).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)

		resp = suite.Expect.GET("/api/v1/user/following/latest").
			WithHeader("Authorization", "Bearer "+token2).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Array().Length().Equal(1)
		resp.Value("data").Array().First().Object().Value("unread").Equal(true)
	})

	suite.T().Run("list latest following unread", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/user/1/status").
			WithHeader("Authorization", "Bearer "+token2).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("data").Array()

		resp = suite.Expect.GET("/api/v1/user/following/latest").
			WithHeader("Authorization", "Bearer "+token2).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Array().Length().Equal(1)
		resp.Value("data").Array().First().Object().Value("unread").Equal(false)
	})
}
