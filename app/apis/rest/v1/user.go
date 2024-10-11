package v1

import (
	"context"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	"github.com/mises-id/sns-apigateway/app/middleware"
	"github.com/mises-id/sns-apigateway/lib/codes"
	"github.com/mssola/user_agent"

	airdroppb "github.com/mises-id/mises-airdropsvc/proto"
	pb "github.com/mises-id/sns-socialsvc/proto"
)

type SignInParams struct {
	Provider  string `json:"provider"`
	Referrer  string `json:"referrer"`
	UserAuthz *struct {
		Auth string `json:"auth"`
	} `json:"user_authz"`
}

type AvatarResp struct {
	Small      string `json:"small"`
	Medium     string `json:"medium"`
	Large      string `json:"large"`
	NftAssetId string `json:"nft_asset_id"`
}
type ComplaintParams struct {
	TargetType string `json:"target_type" body:"target_type"`
	TargetID   string `json:"target_id" body:"target_id"`
	Reason     string `json:"reason" body:"reason"`
}

type UserFullResp struct {
	UID            uint64      `json:"uid"`
	Username       string      `json:"username"`
	Misesid        string      `json:"misesid"`
	Gender         string      `json:"gender"`
	Mobile         string      `json:"mobile"`
	Email          string      `json:"email"`
	Address        string      `json:"address"`
	Intro          string      `json:"intro"`
	Avatar         *AvatarResp `json:"avatar"`
	IsFollowed     bool        `json:"is_followed"`
	IsBlocked      bool        `json:"is_blocked"`
	IsLogined      bool        `json:"is_logined"`
	IsAirdropped   bool        `json:"is_airdropped"`
	AirdropStatus  bool        `json:"airdrop_status"`
	FollowingCount uint64      `json:"followings_count"`
	FansCount      uint64      `json:"fans_count"`
	LikedCount     uint64      `json:"liked_count"`
	NewFansCount   uint64      `json:"new_fans_count"`
}

type UserRestrictedResp struct {
	UID            uint64      `json:"uid"`
	Username       string      `json:"username"`
	Misesid        string      `json:"misesid"`
	Gender         string      `json:"gender"`
	Intro          string      `json:"intro"`
	Avatar         *AvatarResp `json:"avatar"`
	IsFollowed     bool        `json:"is_followed"`
	IsBlocked      bool        `json:"is_blocked"`
	IsLogined      bool        `json:"is_logined"`
	IsAirdropped   bool        `json:"is_airdropped"`
	AirdropStatus  bool        `json:"airdrop_status"`
	FollowingCount uint64      `json:"followings_count"`
	FansCount      uint64      `json:"fans_count"`
	LikedCount     uint64      `json:"liked_count"`
	NewFansCount   uint64      `json:"new_fans_count"`
}

type UserSummaryResp struct {
	UID         uint64      `json:"uid"`
	Username    string      `json:"username"`
	Misesid     string      `json:"misesid"`
	Avatar      *AvatarResp `json:"avatar"`
	HelpMisesid string      `json:"help_misesid"`
	IsFollowed  bool        `json:"is_followed"`
}

type (
	ReceiveAirdropParams struct {
		Tweet string `json:"tweet"`
	}
	UserAgent struct {
		ua        string
		ipaddr    string
		os        string
		browser   string
		platform  string
		device_id string
	}
)

func GetCurrentUID(c echo.Context) uint64 {
	var currentUID uint64
	if c.Get("CurrentUID") != nil {
		currentUID = c.Get("CurrentUID").(uint64)
	}
	return currentUID
}
func GetCurrentMisesID(c echo.Context) string {
	var currentMisesID string
	if c.Get("CurrentMisesID") != nil {
		currentMisesID = c.Get("CurrentMisesID").(string)
	}
	return currentMisesID
}
func GetCurrentEthAddress(c echo.Context) string {
	var currentEthAddress string
	if c.Get("CurrentEthAddress") != nil {
		currentEthAddress = c.Get("CurrentEthAddress").(string)
	}
	return currentEthAddress
}
func GetUIDParam(c echo.Context) (uint64, error) {
	uidParam := c.Param("uid")
	uid, err := strconv.ParseUint(uidParam, 10, 64)
	if err != nil {
		return 0, codes.ErrInvalidArgument.Newf("invalid uid %s", uidParam)
	}
	return uid, nil
}

func SignIn(c echo.Context) error {
	params := &SignInParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}
	if params.UserAuthz == nil {
		return codes.ErrInvalidArgument.New("invalid auth params")
	}
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	user_agent := userAgent(c)
	svcresp, err := grpcsvc.SignIn(ctx, &pb.SignInRequest{
		Auth:     params.UserAuthz.Auth,
		Referrer: params.Referrer,
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
		"token":      svcresp.Jwt,
		"is_created": svcresp.IsCreated,
	})
}

func userAgent(c echo.Context) *UserAgent {
	res := &UserAgent{}
	uastr := c.Request().UserAgent()
	ua := user_agent.New(uastr)
	res.ua = uastr
	res.ipaddr = c.RealIP()
	browserName, browserVersion := ua.Browser()
	res.browser = browserName + " " + browserVersion
	res.os = ua.OS()
	res.platform = ua.Platform()
	res.device_id = c.Request().Header.Get("mises-device-id")
	return res
}

func ReceiveAirdrop(c echo.Context) error {
	params := &ReceiveAirdropParams{}
	if err := c.Bind(params); err != nil {
		return err
	}
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	_, err = grpcsvc.ReceiveAirdrop(ctx, &pb.ReceiveAirdropRequest{
		CurrentUid: GetCurrentUID(c),
		Tweet:      params.Tweet,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, nil)
}

func Complaint(c echo.Context) error {
	params := &ComplaintParams{}
	if err := c.Bind(params); err != nil {
		return err
	}
	uid := GetCurrentUID(c)
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	_, err = grpcsvc.Complaint(ctx, &pb.ComplaintRequest{
		CurrentUid: uid,
		TargetType: params.TargetType,
		TargetId:   params.TargetID,
		Reason:     params.Reason,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, nil)
}

func ShareTweetUrl(c echo.Context) error {
	uid := GetCurrentUID(c)
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.ShareTweetUrl(ctx, &pb.ShareTweetUrlRequest{
		CurrentUid: uid,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, echo.Map{
		"url": svcresp.Url,
	})
}

func MyProfile(c echo.Context) error {
	uid := GetCurrentUID(c)
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.FindUser(ctx, &pb.FindUserRequest{
		Uid:        uid,
		CurrentUid: uid,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, BuildUserFullResp(svcresp.User, false))
}

func FindUser(c echo.Context) error {
	currentUID := GetCurrentUID(c)
	uid, err := GetUIDParam(c)
	if err != nil {
		return err
	}
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.FindUser(ctx, &pb.FindUserRequest{
		Uid:        uid,
		CurrentUid: currentUID,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, BuildUserRestrictedResp(svcresp.User, svcresp.IsFollowed))
}

func FindMisesUser(c echo.Context) error {
	currentUID := GetCurrentUID(c)
	misesidParam := c.Param("misesid")

	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.FindMisesUser(ctx, &pb.FindMisesUserRequest{
		Misesid:    misesidParam,
		CurrentUid: currentUID,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, BuildUserRestrictedResp(svcresp.User, svcresp.IsFollowed))
}

type UserProfileParams struct {
	Gender  string `json:"gender"`
	Mobile  string `json:"mobile"`
	Eamil   string `json:"email"`
	Address string `json:"address"`
	Intro   string `json:"intro"`
}
type UpdateUserConfigParams struct {
	NftState bool `json:"nft_state"`
}

type UserConfig struct {
	NftState bool `json:"nft_state"`
}
type UserNameParams struct {
	Username string `json:"username"`
}

type UserAvatarParams struct {
	AttachmentPath string `json:"attachment_path"`
	NftAssetId     string `json:"nft_asset_id"`
}

type UserUpdateParams struct {
	By       string             `json:"by"`
	Profile  *UserProfileParams `json:"profile"`
	Username *UserNameParams    `json:"username"`
	Avatar   *UserAvatarParams  `json:"avatar"`
}

func UpdateUserConfig(c echo.Context) error {
	params := &UpdateUserConfigParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument
	}
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.UpdateUserConfig(ctx, &pb.UpdateUserConfigRequest{
		CurrentUid: GetCurrentUID(c),
		NftState:   params.NftState,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, BuildUserConfig(svcresp.Config))
}

func GetUserConfig(c echo.Context) error {
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	uid, err := GetUIDParam(c)
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.GetUserConfig(ctx, &pb.GetUserConfigRequest{
		CurrentUid: GetCurrentUID(c),
		Uid:        uid,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, BuildUserConfig(svcresp.Config))
}

func BuildUserConfig(config *pb.UserConfig) *UserConfig {
	if config == nil {
		return nil
	}
	return &UserConfig{
		NftState: config.NftState,
	}
}

func UpdateUser(c echo.Context) error {
	params := &UserUpdateParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument
	}

	uid := c.Get("CurrentUser").(*middleware.UserSession).UID
	var grpcsvc pb.SocialServer
	var ctx context.Context
	var err error
	var serverresp *pb.UpdateUserResponse
	if grpcsvc, ctx, err = rest.GrpcSocialService(); err != nil {
		return err
	}
	switch params.By {
	default:
		return codes.ErrInvalidArgument
	case "profile":
		serverresp, err = grpcsvc.UpdateUserProfile(ctx, &pb.UpdateUserProfileRequest{
			Uid:     uid,
			Gender:  params.Profile.Gender,
			Mobile:  params.Profile.Mobile,
			Email:   params.Profile.Eamil,
			Address: params.Profile.Address,
			Intro:   params.Profile.Intro,
		})
	case "avatar":
		serverresp, err = grpcsvc.UpdateUserAvatar(ctx, &pb.UpdateUserAvatarRequest{
			Uid:            uid,
			AttachmentPath: params.Avatar.AttachmentPath,
			NftAssetId:     params.Avatar.NftAssetId,
		})
	case "username":
		serverresp, err = grpcsvc.UpdateUserName(ctx, &pb.UpdateUserNameRequest{
			Uid:      uid,
			Username: params.Username.Username,
		})
	}
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, BuildUserFullResp(serverresp.User, false))
}

func BuildUserFullResp(user *pb.UserInfo, followed bool) *UserFullResp {
	if user == nil {
		return nil
	}
	resp := &UserFullResp{
		UID:            user.Uid,
		Username:       user.Username,
		Misesid:        user.Misesid,
		Gender:         user.Gender,
		Mobile:         user.Mobile,
		Email:          user.Email,
		Address:        user.Address,
		IsFollowed:     followed,
		IsAirdropped:   user.IsAirdropped,
		AirdropStatus:  user.AirdropStatus,
		IsBlocked:      user.IsBlocked,
		IsLogined:      user.IsLogined,
		FollowingCount: uint64(user.FollowingsCount),
		FansCount:      uint64(user.FansCount),
		LikedCount:     uint64(user.LikedCount),
		NewFansCount:   uint64(user.NewFansCount),
		Intro:          user.Intro,
	}
	if user.AvatarUrl != nil {
		resp.Avatar = &AvatarResp{
			Small:      user.AvatarUrl.Small,
			Medium:     user.AvatarUrl.Medium,
			Large:      user.AvatarUrl.Large,
			NftAssetId: user.AvatarUrl.NftAssetId,
		}
	}
	return resp
}

func BuildUserRestrictedResp(user *pb.UserInfo, followed bool) *UserRestrictedResp {
	if user == nil {
		return nil
	}
	resp := &UserRestrictedResp{
		UID:            user.Uid,
		Username:       user.Username,
		Misesid:        user.Misesid,
		Gender:         user.Gender,
		IsFollowed:     followed,
		IsAirdropped:   user.IsAirdropped,
		AirdropStatus:  user.AirdropStatus,
		IsBlocked:      user.IsBlocked,
		IsLogined:      user.IsLogined,
		FollowingCount: uint64(user.FollowingsCount),
		FansCount:      uint64(user.FansCount),
		LikedCount:     uint64(user.LikedCount),
		NewFansCount:   uint64(user.NewFansCount),
		Intro:          user.Intro,
	}
	if user.AvatarUrl != nil {
		resp.Avatar = &AvatarResp{
			Small:      user.AvatarUrl.Small,
			Medium:     user.AvatarUrl.Medium,
			Large:      user.AvatarUrl.Large,
			NftAssetId: user.AvatarUrl.NftAssetId,
		}
	}
	return resp
}

func BuildUserSummaryResp(user *pb.UserInfo) *UserSummaryResp {
	if user == nil {
		return nil
	}
	resp := &UserSummaryResp{
		UID:         user.Uid,
		Username:    user.Username,
		Misesid:     user.Misesid,
		IsFollowed:  user.IsFollowed,
		HelpMisesid: user.HelpMisesid,
	}
	if user.AvatarUrl != nil {
		resp.Avatar = &AvatarResp{
			Small:      user.AvatarUrl.Small,
			Medium:     user.AvatarUrl.Medium,
			Large:      user.AvatarUrl.Large,
			NftAssetId: user.AvatarUrl.NftAssetId,
		}
	}
	return resp
}
func BuildAirdropSvcUserSummaryResp(user *airdroppb.UserInfo) *UserSummaryResp {
	if user == nil {
		return nil
	}
	resp := &UserSummaryResp{
		UID:         user.Uid,
		Username:    user.Username,
		Misesid:     user.Misesid,
		IsFollowed:  user.IsFollowed,
		HelpMisesid: user.HelpMisesid,
	}
	if user.AvatarUrl != nil {
		resp.Avatar = &AvatarResp{
			Small:      user.AvatarUrl.Small,
			Medium:     user.AvatarUrl.Medium,
			Large:      user.AvatarUrl.Large,
			NftAssetId: user.AvatarUrl.NftAssetId,
		}
	}
	return resp
}

type ListUserLikeParams struct {
	rest.PageQuickParams
}
type UserLikeResp struct {
	Status    *StatusResp `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
}

func BuildUserLikeResplice(in []*pb.StatusLike) []*UserLikeResp {
	resp := []*UserLikeResp{}
	for _, i := range in {
		resp = append(resp, &UserLikeResp{
			Status:    BuildStatusResp(i.Status),
			CreatedAt: time.Unix(int64(i.CreatedAt), 0),
		})
	}

	return resp
}

func PageUserNftAsset(c echo.Context) error {
	uid, err := GetUIDParam(c)
	if err != nil {
		return err
	}
	params := &PageNftAssetParams{}
	if err = c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}
	params.UID = uid

	return PageNftAsset(c, params)
}
func PageNftAsset(c echo.Context, params *PageNftAssetParams) error {

	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.PageNftAsset(ctx, &pb.PageNftAssetRequest{
		CurrentUid: GetCurrentUID(c),
		Uid:        params.UID,
		SortBy:     params.SortBy,
		Scene:      params.Scene,
		Paginator: &pb.PageQuick{
			NextId: params.PageQuickParams.NextID,
			Limit:  uint64(params.PageQuickParams.Limit),
		},
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessRespWithPagination(c, BuildNftAssetRespSlice(svcresp.Assets), svcresp.Paginator)
}

func MyNftAsset(c echo.Context) error {
	params := &PageNftAssetParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}
	params.UID = GetCurrentUID(c)
	return PageNftAsset(c, params)
}

func ListUserLike(c echo.Context) error {
	params := &ListUserLikeParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.Newf("invalid query params")
	}
	uid, err := GetUIDParam(c)
	if err != nil {
		return err
	}
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	paginator := &pb.PageQuick{
		NextId: params.PageQuickParams.NextID,
		Limit:  uint64(params.PageQuickParams.Limit),
	}
	svcresp, err := grpcsvc.ListLikeStatus(ctx, &pb.ListLikeRequest{
		Uid:        uid,
		CurrentUid: GetCurrentUID(c),
		Paginator:  paginator,
	})
	if err != nil {
		return err
	}

	return rest.BuildSuccessRespWithPagination(c, BuildUserLikeResplice(svcresp.Statuses), svcresp.Paginator)
}
