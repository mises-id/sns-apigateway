package route

import (
	"github.com/labstack/echo"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	v1 "github.com/mises-id/sns-apigateway/app/apis/rest/v1"
	appmw "github.com/mises-id/sns-apigateway/app/middleware"
	mw "github.com/mises-id/sns-apigateway/lib/middleware"
)

// SetRoutes sets the routes of echo http server
func SetRoutes(e *echo.Echo) {
	e.GET("/", rest.Probe)
	e.GET("/healthz", rest.Probe)
	groupV1 := e.Group("/api/v1", mw.ErrorResponseMiddleware, appmw.SetCurrentUserMiddleware)
	groupV1.GET("/user/:uid", v1.FindUser)
	groupV1.POST("/signin", v1.SignIn)
	groupV1.GET("/user/:uid/friendship", v1.ListFriendship)

	userGroup := e.Group("/api/v1", mw.ErrorResponseMiddleware, appmw.SetCurrentUserMiddleware, appmw.RequireCurrentUserMiddleware)
	userGroup.POST("/upload", v1.UploadFile)
	userGroup.GET("/user/me", v1.MyProfile)
	userGroup.PATCH("/user/me", v1.UpdateUser)
	userGroup.POST("/user/follow", v1.Follow)
	userGroup.DELETE("/user/follow", v1.Unfollow)
	userGroup.GET("/user/following/latest", v1.LatestFollowing)
	userGroup.GET("/user/blacklist", v1.ListBlacklist)
	userGroup.POST("/user/blacklist", v1.CreateBlacklist)
	userGroup.DELETE("/user/blacklist", v1.DeleteBlacklist)

	groupV1.GET("/user/:uid/status", v1.ListUserStatus)
	groupV1.GET("/status/recommend", v1.RecommendStatus)
	userGroup.GET("/timeline/me", v1.Timeline)
	userGroup.POST("/status", v1.CreateStatus)
	groupV1.GET("/status/:id", v1.GetStatus)
	userGroup.DELETE("/status/:id", v1.DeleteStatus)
	userGroup.POST("/status/:id/like", v1.LikeStatus)
	userGroup.DELETE("/status/:id/like", v1.UnlikeStatus)

	groupV1.GET("/comment", v1.ListComment)
	userGroup.POST("/comment", v1.CreateComment)

	userGroup.GET("/message", v1.ListMessage)
	userGroup.GET("/message/summary", v1.MessageSummary)
	userGroup.PUT("/message/read", v1.ReadMessage)

	groupV1.GET("/user/recommend", v1.RecommendUser)
}
