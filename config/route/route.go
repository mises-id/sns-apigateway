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
	//e.Static("/", "assets")
	e.GET("/", rest.Probe)
	e.GET("/healthz", rest.Probe)
	e.GET("/health/swap", v1.SwapHealth)
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
	groupV1.GET("/website/search", v1.SearchWebsite)
	groupV1.GET("/website/internal_search", v1.WebsiteInternalSearch)
	//extension
	groupV1.GET("/extensions_category/list", v1.ListExtensionsCategory)
	groupV1.GET("/extensions/page", v1.PageExtensions)
	//phishing
	groupV1.POST("/phishing_site/check", v1.PhishingCheck)
	groupV1.GET("/web3safe/verify_contract", v1.VerifyContract)
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
	groupV1.GET("/news", v1.ListNews)
	groupV1.GET("/news/:id", v1.GetNews)

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
		Rate:      0.0001,
		Burst:     20,
		ExpiresIn: 1 * time.Hour,
	})
	rateConfig := middleware.RateLimiterConfig{
		Store:       storeC,
		DenyHandler: mw.ErrTooManyRequestFunc,
	}
	//swap
	swapRateConfigCommon := getSwapRateConfigCommon()
	swapRateConfiWithUserWalletAddress := getSwapRateConfigWithUserWalletAddress()
	swapCommonRateLimiter := middleware.RateLimiterWithConfig(swapRateConfigCommon)
	swapRateLimiterWithUserWalletAddress := middleware.RateLimiterWithConfig(swapRateConfiWithUserWalletAddress)
	swapGroup := e.Group("/api/v1", swapCommonRateLimiter, swapRateLimiterWithUserWalletAddress, mw.ErrorResponseMiddleware, appmw.SetCurrentUserMiddleware)
	swapGroup.GET("/swap/order/:from_address", v1.PageSwapOrder)
	swapGroup.GET("/swap/order/:from_address/:tx_hash", v1.FindSwapOrder)
	swapGroup.GET("/swap/approve/allowance", v1.GetSwapApproveAllowance)
	swapGroup.GET("/swap/approve/transaction", v1.ApproveSwapTransaction)
	swapGroup.GET("/swap/trade", v1.SwapTrade)
	swapGroup.GET("/swap/quote", v1.SwapQuote)
	swapGroup.POST("/swap/wallets_and_tokens", v1.WalletsAndTokens)
	swapGroup.GET("/swap/token/list", v1.ListTokens, middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	// bridge
	bridgeRateLimiterWithIP := middleware.RateLimiterWithConfig(getBridgeRateLimiterWithIPConfig())
	//bridgeGroup := e.Group("/api/v1", bridgeRateLimiterWithIP, mw.ErrorResponseMiddleware, appmw.SetCurrentUserMiddleware, appmw.RequireCurrentUserMiddleware)
	bridgeGroup := e.Group("/api/v1", bridgeRateLimiterWithIP, mw.ErrorResponseMiddleware)
	bridgeGroup.POST("/bridge/get_currencies", v1.BridgeGetCurrencies)
	bridgeGroup.POST("/bridge/get_pairs_params", v1.BridgeGetPairsParams)
	bridgeGroup.POST("/bridge/get_exchange_amount", v1.BridgeGetExchangeAmount)
	bridgeGroup.POST("/bridge/create_transaction", v1.BridgeCreateTransaction, appmw.SetCurrentUserMiddleware, appmw.RequireCurrentUserMiddleware)
	bridgeGroup.POST("/bridge/get_transaction_info", v1.BridgeGetTransactionInfo)
	bridgeGroup.POST("/bridge/get_transaction_status", v1.BridgeGetTransactionStatus)
	bridgeGroup.POST("/bridge/validate_address", v1.BridgeValidateAddress)
	bridgeGroup.POST("/bridge/get_fix_rate_for_amount", v1.BridgeGetFixRateForAmount)
	bridgeGroup.POST("/bridge/create_fix_transaction", v1.BridgeCreateFixTransaction, appmw.SetCurrentUserMiddleware, appmw.RequireCurrentUserMiddleware)
	bridgeGroup.POST("/bridge/history_list", v1.BridgeHistoryList, appmw.SetCurrentUserMiddleware, appmw.RequireCurrentUserMiddleware)

	userGroup.GET("/twitter/auth_url", v1.TwitterAuthUrl, middleware.RateLimiterWithConfig(rateConfig))
	//userGroup.GET("/twitter/auth_url", v1.TwitterAuthUrl)
	userGroup.GET("/airdrop/info", v1.AirdropInfo)
	userGroup.POST("/airdrop/receive", v1.ReceiveAirdrop)

	// mining
	redeemBonusRateConfigWithUser := middleware.RateLimiterWithConfig(getRedeemBonusRateConfigWithUser())
	groupV1.GET("/admob/ssv", v1.ADMobSSV)
	groupV1.GET("/adcallback/mintegral", v1.MintegralCallback)
	groupV1.GET("/ad_mining/estimate_bonus", v1.EstimateAdBonus)
	groupV1.GET("/mb_airdrop/user/:misesid", v1.FindMBAirdropUser)
	userGroup.GET("/mb_airdrop/claim", v1.ClaimMBAirdrop, redeemBonusRateConfigWithUser)
	groupV1.GET("/mining/config", v1.GeMiningConfig)
	userGroup.GET("/mining/bonus", v1.GetBonus)
	userGroup.GET("/ad_mining/me", v1.MyAdMining)
	userGroup.POST("/mining/redeem_bonus", v1.RedeemBonus, redeemBonusRateConfigWithUser)
	userGroup.POST("/ad_mining/log", v1.AdMiningLog, redeemBonusRateConfigWithUser)
}

func getBridgeRateLimiterWithIPConfig() middleware.RateLimiterConfig {
	bridgeRateLimitStore := middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
		Rate:      10,
		Burst:     10,
		ExpiresIn: 1 * time.Minute,
	})
	bridgeRateConfigWithIP := middleware.RateLimiterConfig{
		Store:       bridgeRateLimitStore,
		DenyHandler: mw.ErrTooManyRequestFunc,
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP() + ctx.Path()
			return id, nil
		},
	}
	return bridgeRateConfigWithIP
}

func getRedeemBonusRateConfigWithUser() middleware.RateLimiterConfig {

	redeemBonusRateLimitStore := middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
		Rate:      1,
		Burst:     1,
		ExpiresIn: 1 * time.Second,
	})
	redeemBonusRateConfigWithUser := middleware.RateLimiterConfig{
		Store:       redeemBonusRateLimitStore,
		DenyHandler: mw.ErrTooManyRequestFunc,
		IdentifierExtractor: func(c echo.Context) (string, error) {
			id := getCurrentEthAddress(c)
			return id, nil
		},
	}

	return redeemBonusRateConfigWithUser
}

func getSwapRateConfigCommon() middleware.RateLimiterConfig {
	swapCommonStore := middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
		Rate:      100,
		Burst:     1000,
		ExpiresIn: 1 * time.Minute,
	})
	swapRateConfigCommon := middleware.RateLimiterConfig{
		Store:       swapCommonStore,
		DenyHandler: mw.ErrTooManyRequestFunc,
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP() + ctx.Path()
			return id, nil
		},
		Skipper: func(c echo.Context) bool {
			userWalletAddress := getgUserWalletAddress(c)
			if userWalletAddress != "" {
				return true
			}
			return false
		},
	}
	return swapRateConfigCommon
}

func getSwapRateConfigWithUserWalletAddress() middleware.RateLimiterConfig {
	swapStoreWithUserWalletAddress := middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
		Rate:      5,
		Burst:     60,
		ExpiresIn: 1 * time.Minute,
	})
	swapRateConfigWithUserWalletAddress := middleware.RateLimiterConfig{
		Store:       swapStoreWithUserWalletAddress,
		DenyHandler: mw.ErrTooManyRequestFunc,
		IdentifierExtractor: func(c echo.Context) (string, error) {
			id := getgUserWalletAddress(c) + c.Path()
			return id, nil
		},
		Skipper: func(c echo.Context) bool {
			userWalletAddress := getgUserWalletAddress(c)
			if userWalletAddress == "" {
				return true
			}
			return false
		},
	}
	return swapRateConfigWithUserWalletAddress
}

func getgUserWalletAddress(c echo.Context) string {
	return c.Request().Header.Get("User-Wallet-Address")
}
func getCurrentEthAddress(c echo.Context) string {
	var currentEthAddress string
	if c.Get("CurrentEthAddress") != nil {
		currentEthAddress = c.Get("CurrentEthAddress").(string)
	}
	return currentEthAddress
}
