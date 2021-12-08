package v1

import (
	"context"
	"strconv"

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

type UserResp struct {
	UID        uint64      `json:"uid"`
	Username   string      `json:"username"`
	Misesid    string      `json:"misesid"`
	Gender     string      `json:"gender"`
	Mobile     string      `json:"mobile"`
	Email      string      `json:"email"`
	Address    string      `json:"address"`
	Avatar     *AvatarResp `json:"avatar"`
	IsFollowed bool        `json:"is_followed"`
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
	return rest.BuildSuccessResp(c, buildUserResp(svcresp.User, false))
}

func FindUser(c echo.Context) error {
	uidParam := c.Param("uid")
	uid, err := strconv.ParseUint(uidParam, 10, 64)
	if err != nil {
		return codes.ErrInvalidArgument.Newf("invalid uid %s", uidParam)
	}
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.FindUser(ctx, &pb.FindUserRequest{Uid: uid})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, buildUserResp(svcresp.User, svcresp.IsFollowed))
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
	AttachmentID uint64 `json:"attachment_id"`
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
			Uid:          uid,
			AttachmentId: params.Avatar.AttachmentID,
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
	return rest.BuildSuccessResp(c, buildUserResp(serverresp.User, false))
}

func buildUserResp(user *pb.UserInfo, followed bool) *UserResp {
	if user == nil {
		return nil
	}
	resp := &UserResp{
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
