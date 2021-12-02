package v1

import (
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/mises-id/apigateway/app/apis/rest"
	"github.com/mises-id/apigateway/lib/codes"

	pb "github.com/mises-id/socialsvc/proto"
)

type ListFriendshipParams struct {
	rest.PageQuickParams
	RelationType string `query:"relation_type"`
}

type FollowParams struct {
	ToUserID uint64 `json:"to_user_id" query:"to_user_id"`
}

type FriendshipResp struct {
	User         *UserResp `json:"user"`
	RelationType string    `json:"relation_type"`
	CreatedAt    time.Time `json:"created_at"`
}

func ListFriendship(c echo.Context) error {
	uidParam := c.Param("uid")
	uid, err := strconv.ParseUint(uidParam, 10, 64)
	if err != nil {
		return codes.ErrInvalidArgument.Newf("invalid uid %s", uidParam)
	}
	params := &ListFriendshipParams{}
	if err := c.Bind(params); err != nil {
		return err
	}
	if len(params.RelationType) == 0 {
		params.RelationType = "fan"
	}

	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.ListRelationship(ctx, &pb.ListRelationshipRequest{
		CurrentUid:   uid,
		RelationType: params.RelationType,
		Paginator: &pb.PageQuick{
			NextId: params.PageQuickParams.NextID,
			Limit:  uint64(params.PageQuickParams.Limit),
		},
	})
	if err != nil {
		return err
	}

	return rest.BuildSuccessRespWithPagination(c, svcresp.Relations, svcresp.Paginator)
}

func Follow(c echo.Context) error {
	params := &FollowParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument
	}

	var currentUID uint64
	if c.Get("CurrentUID") != nil {
		currentUID = c.Get("CurrentUID").(uint64)
	} else {
		return codes.ErrInvalidArgument
	}
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}

	_, err = grpcsvc.Follow(ctx, &pb.FollowRequest{
		CurrentUid: currentUID,
		TargetUid:  params.ToUserID,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, nil)
}

func Unfollow(c echo.Context) error {
	params := &FollowParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument
	}

	var currentUID uint64
	if c.Get("CurrentUID") != nil {
		currentUID = c.Get("CurrentUID").(uint64)
	} else {
		return codes.ErrInvalidArgument
	}
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}

	_, err = grpcsvc.UnFollow(ctx, &pb.UnFollowRequest{
		CurrentUid: currentUID,
		TargetUid:  params.ToUserID,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, nil)
}
