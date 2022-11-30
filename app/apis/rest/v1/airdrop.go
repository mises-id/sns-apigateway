package v1

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	pb "github.com/mises-id/mises-airdropsvc/proto"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
)

type (
	TwitterCallbackParam struct {
		UID           uint64 `json:"uid" query:"uid"`
		State         string `json:"state" query:"state"`
		OauthToken    string `json:"oauth_token" query:"oauth_token"`
		OauthVerifier string `json:"oauth_verifier" query:"oauth_verifier"`
	}
	UserTwitterAuthResp struct {
		TwitterUserId    string    `json:"twitter_user_id" bson:"twitter_user_id"`
		Misesid          string    `json:"misesid"`
		Name             string    `json:"name" bson:"name"`
		Username         string    `json:"username" bson:"username"`
		CheckState       string    `json:"check_state" bson:"check_state"`
		InvalidCode      string    `json:"invalid_code" bson:"invalid_code"`
		Reason           string    `json:"reason" bson:"reason"`
		FollowersCount   uint64    `json:"followers_count" bson:"followers_count"`
		TweetCount       uint64    `json:"tweet_count" bson:"tweet_count"`
		TwitterCreatedAt time.Time `json:"twitter_created_at" bson:"twitter_created_at"`
		Amount           float32   `json:"amount" bson:"amount"`
		CreatedAt        time.Time `json:"created_at" bson:"created_at"`
	}
	AirdropResp struct {
		Coin      float32   `json:"coin" bson:"coin"`
		CreatedAt time.Time `json:"created_at" bson:"created_at"`
		FinishAt  time.Time `json:"finish_at" bson:"finish_at"`
		Status    string    `json:"status" bson:"status"`
	}
	AirdropInfoResp struct {
		Twitter *UserTwitterAuthResp `json:"twitter"`
		Airdrop *AirdropResp         `json:"airdrop"`
	}
)

func TwitterAuthUrl(c echo.Context) error {
	uid := GetCurrentUID(c)
	grpcsvc, ctx, err := rest.GrpcAirdropService()
	if err != nil {
		return err
	}
	user_agent := userAgent(c)
	svcresp, err := grpcsvc.GetTwitterAuthUrl(ctx, &pb.GetTwitterAuthUrlRequest{
		CurrentUid: uid,
		UserAgent: &pb.UserAgent{
			Ua:       user_agent.ua,
			Ipaddr:   user_agent.ipaddr,
			Os:       user_agent.os,
			Browser:  user_agent.browser,
			Platform: user_agent.platform,
			DeviceId: user_agent.device_id,
		},
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, echo.Map{
		"url": svcresp.Url,
	})
}

func TwitterCallback(c echo.Context) error {
	params := &TwitterCallbackParam{}
	if err := c.Bind(params); err != nil {
		return err
	}
	grpcsvc, ctx, err := rest.GrpcAirdropService()
	if err != nil {
		return err
	}

	user_agent := userAgent(c)
	svcresp, err := grpcsvc.TwitterCallback(ctx, &pb.TwitterCallbackRequest{
		CurrentUid:    params.UID,
		OauthToken:    params.OauthToken,
		OauthVerifier: params.OauthVerifier,
		State:         params.State,
		UserAgent: &pb.UserAgent{
			Ua:       user_agent.ua,
			Ipaddr:   user_agent.ipaddr,
			Os:       user_agent.os,
			Browser:  user_agent.browser,
			Platform: user_agent.platform,
			DeviceId: user_agent.device_id,
		},
	})
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusMovedPermanently, svcresp.Url)

}

func AirdropInfo(c echo.Context) error {
	uid := GetCurrentUID(c)
	grpcsvc, ctx, err := rest.GrpcAirdropService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.GetAirdropInfo(ctx, &pb.GetAirdropInfoRequest{
		CurrentUid: uid,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, BuildAirdropInfoResp(svcresp))
}

func BuildAirdropInfoResp(in *pb.GetAirdropInfoResponse) *AirdropInfoResp {
	if in == nil {
		return nil
	}
	out := &AirdropInfoResp{
		Airdrop: BuildAirdropResp(in.Airdrop),
		Twitter: BuildUserTwitterAuthResp(in.Twitter),
	}
	return out
}

func BuildAirdropResp(in *pb.Airdrop) *AirdropResp {
	if in == nil {
		return nil
	}
	out := &AirdropResp{
		Coin:      in.Coin,
		CreatedAt: time.Unix(int64(in.CreatedAt), 0),
		FinishAt:  time.Unix(int64(in.FinishAt), 0),
		Status:    in.Status,
	}
	return out
}
func BuildUserTwitterAuthResp(in *pb.UserTwitterAuth) *UserTwitterAuthResp {
	if in == nil {
		return nil
	}
	out := &UserTwitterAuthResp{
		TwitterUserId:    in.TwitterUserId,
		Misesid:          in.Misesid,
		Name:             in.Name,
		Username:         in.Username,
		Reason:           in.Reason,
		CheckState:       in.CheckState,
		InvalidCode:      in.InvalidCode,
		FollowersCount:   in.FollowersCount,
		TweetCount:       in.TweetCount,
		Amount:           in.Amount,
		CreatedAt:        time.Unix(int64(in.CreatedAt), 0),
		TwitterCreatedAt: time.Unix(int64(in.TwitterCreatedAt), 0),
	}
	return out
}
