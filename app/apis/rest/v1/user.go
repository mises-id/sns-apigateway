package v1

import (
	"context"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	"github.com/mises-id/sns-apigateway/app/middleware"
	"github.com/mises-id/sns-apigateway/lib/codes"

	pb "github.com/mises-id/sns-socialsvc/proto"
)

type SignInParams struct {
	Provider  string `json:"provider"`
	UserAuthz *struct {
		Auth string `json:"auth"`
	} `json:"user_authz"`
}

type AvatarResp struct {
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

type UserFullResp struct {
	UID            uint64      `json:"uid"`
	Username       string      `json:"username"`
	Misesid        string      `json:"misesid"`
	Gender         string      `json:"gender"`
	Mobile         string      `json:"mobile"`
	Email          string      `json:"email"`
	Address        string      `json:"address"`
	Avatar         *AvatarResp `json:"avatar"`
	IsFollowed     bool        `json:"is_followed"`
	FollowingCount uint64      `json:"followings_count"`
	FansCount      uint64      `json:"fans_count"`
	LikedCount     uint64      `json:"liked_count"`
	NewFansCount   uint64      `json:"new_fans_count"`
}

type UserSummaryResp struct {
	UID        uint64      `json:"uid"`
	Username   string      `json:"username"`
	Misesid    string      `json:"misesid"`
	Avatar     *AvatarResp `json:"avatar"`
	IsFollowed bool        `json:"is_followed"`
}

func GetCurrentUID(c echo.Context) uint64 {
	var currentUID uint64
	if c.Get("CurrentUID") != nil {
		currentUID = c.Get("CurrentUID").(uint64)
	}
	return currentUID
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
		return err
	}
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.SignIn(ctx, &pb.SignInRequest{Auth: params.UserAuthz.Auth})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, echo.Map{
		"token": svcresp.Jwt,
	})
}

func MyProfile(c echo.Context) error {
	uid := c.Get("CurrentUser").(*middleware.UserSession).UID
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.FindUser(ctx, &pb.FindUserRequest{Uid: uid})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, BuildUserFullResp(svcresp.User, false))
}

func FindUser(c echo.Context) error {
	uid, err := GetUIDParam(c)
	if err != nil {
		return err
	}
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.FindUser(ctx, &pb.FindUserRequest{Uid: uid})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, BuildUserFullResp(svcresp.User, svcresp.IsFollowed))
}

type UserProfileParams struct {
	Gender  string `json:"gender"`
	Mobile  string `json:"mobile"`
	Eamil   string `json:"email"`
	Address string `json:"address"`
}

type UserNameParams struct {
	Username string `json:"username"`
}

type UserAvatarParams struct {
	AttachmentPath string `json:"attachment_path"`
}

type UserUpdateParams struct {
	By       string             `json:"by"`
	Profile  *UserProfileParams `json:"profile"`
	Username *UserNameParams    `json:"username"`
	Avatar   *UserAvatarParams  `json:"avatar"`
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
		})
	case "avatar":
		serverresp, err = grpcsvc.UpdateUserAvatar(ctx, &pb.UpdateUserAvatarRequest{
			Uid:            uid,
			AttachmentPath: params.Avatar.AttachmentPath,
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
		UID:        user.Uid,
		Username:   user.Username,
		Misesid:    user.Misesid,
		Gender:     user.Gender,
		Mobile:     user.Mobile,
		Email:      user.Email,
		Address:    user.Address,
		IsFollowed: followed,
	}
	if len(user.Avatar) > 0 {
		resp.Avatar = &AvatarResp{
			// TODO support multiple sizes avatar
			Small:  user.Avatar,
			Medium: user.Avatar,
			Large:  user.Avatar,
		}
	}
	return resp
}

func BuildUserSummaryResp(user *pb.UserInfo) *UserSummaryResp {
	if user == nil {
		return nil
	}
	resp := &UserSummaryResp{
		UID:        user.Uid,
		Username:   user.Username,
		Misesid:    user.Misesid,
		IsFollowed: user.IsFollowed,
	}
	if len(user.Avatar) > 0 {
		resp.Avatar = &AvatarResp{
			// TODO support multiple sizes avatar
			Small:  user.Avatar,
			Medium: user.Avatar,
			Large:  user.Avatar,
		}
	}
	return resp
}

type BlackParams struct {
	UID string `json:"uid"`
}

type ListBlackParams struct {
	rest.PageQuickParams
}

type BlackResp struct {
	User      *UserSummaryResp `json:"user"`
	CreatedAt time.Time        `json:"created_at"`
}

func BuildBlackRespSlice() []*BlackResp {
	return []*BlackResp{
		{
			User:      BuildUserSummaryResp(&pb.UserInfo{}),
			CreatedAt: time.Now(),
		},
	}
}

func ListBlack(c echo.Context) error {
	params := &ListBlackParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.Newf("invalid query params")
	}
	// grpcsvc, ctx, err := rest.GrpcSocialService()
	// if err != nil {
	// 	return err
	// }
	paginator := &pb.PageQuick{
		NextId: params.PageQuickParams.NextID,
		Limit:  uint64(params.PageQuickParams.Limit),
	}

	return rest.BuildSuccessRespWithPagination(c, BuildBlackRespSlice(), paginator)
}

func Black(c echo.Context) error {
	params := &BlackParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.Newf("invalid query params")
	}
	// GetCurrentUID(c)

	return rest.BuildSuccessResp(c, nil)
}

func UnBlack(c echo.Context) error {
	// uid, err := GetUIDParam(c)
	// if err != nil {
	// 	return err
	// }
	// GetCurrentUID(c)

	return rest.BuildSuccessResp(c, nil)
}

type ListUserLikeParams struct {
	rest.PageQuickParams
}

type UserLikeResp struct {
	Status    *StatusResp `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
}

func BuildUserLikeResplice() []*UserLikeResp {
	return []*UserLikeResp{
		{
			Status:    BuildStatusResp(&pb.StatusInfo{}),
			CreatedAt: time.Now(),
		},
	}
}

func ListUserLike(c echo.Context) error {
	params := &ListUserLikeParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.Newf("invalid query params")
	}
	// grpcsvc, ctx, err := rest.GrpcSocialService()
	// if err != nil {
	// 	return err
	// }
	paginator := &pb.PageQuick{
		NextId: params.PageQuickParams.NextID,
		Limit:  uint64(params.PageQuickParams.Limit),
	}

	return rest.BuildSuccessRespWithPagination(c, BuildUserLikeResplice(), paginator)
}
