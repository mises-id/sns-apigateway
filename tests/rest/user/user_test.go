//go:build tests
// +build tests

package user

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
)

type UserServerSuite struct {
	rest.RestBaseTestSuite
	collections []string
}

func (suite *UserServerSuite) SetupSuite() {
	suite.RestBaseTestSuite.SetupSuite()
	suite.collections = []string{"counters", "attachments", "users", "blacklists"}

}

func (suite *UserServerSuite) TearDownSuite() {
	suite.RestBaseTestSuite.TearDownSuite()
}

func (suite *UserServerSuite) SetupTest() {
	suite.Clean(suite.collections...)
	suite.Acquire(suite.collections...)

	factories.InitUsers(&models.User{
		UID:        1001,
		Gender:     enum.GenderFemale,
		AvatarPath: "DummyAvatarPath",
		Misesid:    "1001",
		Username:   "",
	}, &models.User{
		UID:     1002,
		Gender:  enum.GenderFemale,
		Misesid: "1002",
	})

}

func (suite *UserServerSuite) TearDownTest() {
	suite.Clean(suite.collections...)
}

func TestUserServer(t *testing.T) {
	suite.Run(t, &UserServerSuite{})
}

func (suite *UserServerSuite) TestFindUser() {

	suite.T().Run("not found user", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/user/999").
			Expect().Status(http.StatusNotFound).JSON().Object()
		resp.Value("code").Equal(404000)
	})

	suite.T().Run("find user with avatar", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/user/1001").
			Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Object().Value("uid").Equal(1001)
		url := fmt.Sprintf("http://localhost/DummyAvatarPath?")
		resp.Value("data").Object().Value("avatar").Object().Value("small").Equal(url)
	})

	suite.T().Run("find user without avatar", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/user/1002").
			Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Object().Value("uid").Equal(1002)
		resp.Value("data").Object().Value("avatar").Null()
	})
}

func (suite *UserServerSuite) TestSignin() {

	suite.T().Run("user signin success", func(t *testing.T) {
		resp := suite.Expect.POST("/api/v1/signin").WithJSON(map[string]interface{}{
			"provider": "mises",
			"user_authz": map[string]interface{}{
				"auth": "1001:1001",
			},
		}).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
	})
	suite.T().Run("create new success", func(t *testing.T) {
		resp := suite.Expect.POST("/api/v1/signin").WithJSON(map[string]interface{}{
			"provider": "mises",
			"user_authz": map[string]interface{}{
				"auth": "1002:1002",
			},
		}).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
	})
}

func (suite *UserServerSuite) TestUpdateUser() {
	token := suite.MockLoginUser("1001:1001")
	suite.T().Run("update username success", func(t *testing.T) {
		resp := suite.Expect.PATCH("/api/v1/user/me").WithJSON(map[string]interface{}{
			"by": "username",
			"username": map[string]interface{}{
				"username": "Hahaha",
			},
		}).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		u := &models.User{}
		err := db.ODM(context.Background()).First(u, bson.M{"_id": 1001}).Error
		suite.Nil(err)
		suite.Equal("Hahaha", u.Username)

		resp = suite.Expect.PATCH("/api/v1/user/me").WithJSON(map[string]interface{}{
			"by": "username",
			"username": map[string]interface{}{
				"username": "..Hello",
			},
		}).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusUnprocessableEntity).JSON().Object()
		resp.Value("code").Equal(codes.UsernameExistedCode)
	})
	suite.T().Run("update user avatar success", func(t *testing.T) {
		resp := suite.Expect.PATCH("/api/v1/user/me").WithJSON(map[string]interface{}{
			"by": "avatar",
			"avatar": map[string]interface{}{
				"attachment_path": "1",
			},
		}).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		u := &models.User{}
		err := db.ODM(context.Background()).First(u, bson.M{"_id": 1001}).Error
		suite.Nil(err)
		suite.Equal("1", u.AvatarPath)
	})
	suite.T().Run("update user profile success", func(t *testing.T) {
		suite.Expect.GET("/api/v1/user/me").WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK)

		resp := suite.Expect.PATCH("/api/v1/user/me").WithJSON(map[string]interface{}{
			"by": "profile",
			"profile": map[string]interface{}{
				"email":   "test@t.com",
				"gender":  "female",
				"mobile":  "123456",
				"address": "xxxx",
			},
		}).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		u := &models.User{}
		err := db.ODM(context.Background()).First(u, bson.M{"_id": 1001}).Error
		suite.Nil(err)
		suite.Equal("test@t.com", u.Email)
		suite.Equal(enum.GenderFemale, u.Gender)
		suite.Equal("123456", u.Mobile)
		suite.Equal("xxxx", u.Address)

		resp = suite.Expect.GET("/api/v1/user/me").WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("data").Object().Value("email").Equal("test@t.com")
		resp.Value("data").Object().Value("mobile").Equal("123456")
		resp.Value("data").Object().Value("address").Equal("xxxx")
	})

}

func (suite *UserServerSuite) TestBlackUser() {

	token := suite.MockLoginUser("1001:1001")

	suite.T().Run("empty blacklist", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/user/blacklist").WithHeader("Authorization", "Bearer "+token).
			Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Array().Length().Equal(0)
	})

	suite.T().Run("add to black list", func(t *testing.T) {
		resp := suite.Expect.POST("/api/v1/user/blacklist").WithJSON(map[string]interface{}{
			"uid": 1002,
		}).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)

	})

	suite.T().Run("one blacklist", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/user/blacklist").WithHeader("Authorization", "Bearer "+token).
			Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Array().Length().Equal(1)
		resp.Value("data").Array().First().Object().Value("user").Object().Value("uid").Equal(1002)
	})

	suite.T().Run("remove from black list", func(t *testing.T) {
		resp := suite.Expect.DELETE("/api/v1/user/blacklist/1002").
			WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)

	})

	suite.T().Run("cleared blacklist", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/user/blacklist").WithHeader("Authorization", "Bearer "+token).
			Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Array().Length().Equal(0)
	})
}
