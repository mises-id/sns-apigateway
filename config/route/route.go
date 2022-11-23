package route

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	v1 "github.com/mises-id/sns-apigateway/app/apis/rest/v1"
	appmw "github.com/mises-id/sns-apigateway/app/middleware"
	mw "github.com/mises-id/sns-apigateway/lib/middleware"
)

// SetRoutes sets the routes of echo http server
func SetRoutes(e *echo.Echo) {
	e.Static("/", "assets")
	e.GET("/", rest.Probe)
	e.GET("/healthz", rest.Probe)
	groupV1 := e.Group("/api/v1", mw.ErrorResponseMiddleware, appmw.SetCurrentUserMiddleware)
	groupOpensea := e.Group("/api/v1", middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(4)), mw.ErrorResponseMiddleware, appmw.SetCurrentUserMiddleware, appmw.RequireCurrentUserMiddleware)
	groupV1.GET("/user/:uid", v1.FindUser)
	groupV1.GET("/mises_user/:misesid", v1.FindMisesUser)
	groupV1.GET("/channel_user/:misesid", v1.GetChannelUser)
	groupV1.GET("/channel/info", v1.ChannelInfo)
	groupV1.GET("/channel_user/page", v1.PageChannelUser)
	groupOpensea.GET("/opensea/single_asset", v1.GetOpenseaAsset)
	groupOpensea.GET("/opensea/assets", v1.ListOpenseaAsset)
	groupOpensea.GET("/opensea/assets_contract", v1.GetOpenseaAssetContract)
	groupV1.POST("/signin", v1.SignIn)
	groupV1.POST("/complaint", v1.Complaint)
	groupV1.GET("/twitter/callback", v1.TwitterCallback)
	groupV1.GET("/user/:uid/friendship", v1.ListFriendship)
	//website
	groupV1.GET("/website_category/list", v1.ListWebsiteCategory)
	groupV1.GET("/website/page", v1.PageWebsite)
	//extension
	groupV1.GET("/extensions_category/list", v1.ListExtensionsCategory)
	groupV1.GET("/extensions/page", v1.PageExtensions)
	userGroup := e.Group("/api/v1", mw.ErrorResponseMiddleware, appmw.SetCurrentUserMiddleware, appmw.RequireCurrentUserMiddleware)

	userGroup.POST("/upload", v1.UploadFile, middleware.BodyLimitWithConfig(middleware.BodyLimitConfig{
		Skipper: middleware.DefaultSkipper,
		Limit:   "8M",
	}))
	userGroup.GET("/user/me", v1.MyProfile)
	userGroup.GET("/user/:uid/config", v1.GetUserConfig)
	userGroup.GET("/share/twitter", v1.ShareTweetUrl)
	userGroup.PATCH("/user/me", v1.UpdateUser)
	userGroup.PATCH("/user/config", v1.UpdateUserConfig)
	userGroup.POST("/user/follow", v1.Follow)
	userGroup.DELETE("/user/follow", v1.Unfollow)
	userGroup.GET("/user/following/latest", v1.LatestFollowing)
	userGroup.GET("/user/blacklist", v1.ListBlacklist)
	userGroup.POST("/user/blacklist", v1.CreateBlacklist)
	userGroup.DELETE("/user/blacklist/:uid", v1.DeleteBlacklist)

	groupV1.GET("/user/:uid/like", v1.ListUserLike)
	groupV1.GET("/user/:uid/status", v1.ListUserStatus)
	groupV1.GET("/status/recommend", v1.RecommendStatus)
	groupV1.GET("/status/list", v1.ListStatus)
	groupV1.GET("/status/recent", v1.RecentStatus)
	userGroup.GET("/timeline/me", v1.Timeline)
	userGroup.POST("/status", v1.CreateStatus)
	userGroup.PATCH("/status/:id", v1.UpdateStatus)
	groupV1.GET("/status/:id", v1.GetStatus)
	userGroup.DELETE("/status/:id", v1.DeleteStatus)
	userGroup.POST("/status/:id/like", v1.LikeStatus)
	userGroup.DELETE("/status/:id/like", v1.UnlikeStatus)

	groupV1.GET("/nft_asset/:id", v1.GetNftAsset)
	groupV1.GET("/nft_asset/:id/event", v1.PageNftEvent)
	groupV1.GET("/nft_asset/:id/like", v1.ListNftAssetLike)
	userGroup.POST("/nft_asset/:id/like", v1.LikeNftAsset)
	userGroup.DELETE("/nft_asset/:id/like", v1.UnlikeNftAsset)
	groupV1.GET("/user/:uid/nft_asset", v1.PageUserNftAsset)
	userGroup.GET("/user/nft_asset", v1.MyNftAsset)

	groupV1.GET("/comment", v1.ListComment)
	userGroup.POST("/comment", v1.CreateComment)
	groupV1.GET("/comment/:id", v1.GetComment)
	userGroup.DELETE("/comment/:id", v1.DeleteComment)
	userGroup.POST("/comment/:id/like", v1.LikeComment)
	userGroup.DELETE("/comment/:id/like", v1.UnlikeComment)

	userGroup.GET("/user/message", v1.ListMessage)
	userGroup.GET("/user/message/summary", v1.MessageSummary)
	userGroup.PUT("/message/read", v1.ReadMessage)

	groupV1.GET("/user/recommend", v1.RecommendUser)

	groupV1.GET("/mises/gasprices", v1.GasPrices)
	groupV1.GET("/mises/chaininfo", v1.ChainInfo)

	storeC := middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
		Rate:      0.01,
		Burst:     5,
		ExpiresIn: 8 * time.Hour,
	})
	rateConfig := middleware.RateLimiterConfig{
		Store:       storeC,
		DenyHandler: mw.ErrTooManyRequestFunc,
	}
	userGroup.GET("/twitter/auth_url", v1.TwitterAuthUrl, middleware.RateLimiterWithConfig(rateConfig))
	//userGroup.GET("/twitter/auth_url", v1.TwitterAuthUrl)
	userGroup.GET("/airdrop/info", v1.AirdropInfo)
	userGroup.POST("/airdrop/receive", v1.ReceiveAirdrop)
}
