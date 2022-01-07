package v1

import (
	"time"

	"github.com/labstack/echo"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	"github.com/mises-id/sns-apigateway/lib/codes"

	pb "github.com/mises-id/sns-socialsvc/proto"
)

type ListFriendshipParams struct {
	rest.PageQuickParams
	RelationType string `query:"relation_type"`
}

type FollowParams struct {
	ToUserID uint64 `json:"to_user_id" query:"to_user_id"`
}

type FriendshipResp struct {
	User         *UserSummaryResp `json:"user"`
	RelationType string           `json:"relation_type"`
	CreatedAt    time.Time        `json:"created_at"`
}

type FollowingResp struct {
	User   *UserSummaryResp `json:"user"`
	Unread bool             `json:"unread"`
}

func BuildFriendshipResp(info *pb.RelationInfo) *FriendshipResp {
	return &FriendshipResp{
		User:         BuildUserSummaryResp(info.User),
		RelationType: info.RelationType,
		CreatedAt:    time.Unix(int64(info.CreatedAt), 0),
	}
}

func BuildFriendshipRespSlice(infos []*pb.RelationInfo) []*FriendshipResp {
	resp := []*FriendshipResp{}
	for _, info := range infos {
		resp = append(resp, BuildFriendshipResp(info))
	}

	return resp
}
func BuildFollowingsRespSlice(infos []*pb.Following) []*FollowingResp {
	resp := []*FollowingResp{}
	if infos == nil {
		return []*FollowingResp{}
	} else {
		for _, info := range infos {
			resp = append(resp, &FollowingResp{BuildUserSummaryResp(info.User), info.Unread})
		}
	}
	return resp
}

func BuildRecommendUserRespSlice(infos []*pb.UserInfo) []*UserSummaryResp {
	resp := []*UserSummaryResp{}
	if infos == nil {
		return []*UserSummaryResp{}
	} else {
		for _, info := range infos {
			resp = append(resp, BuildUserSummaryResp(info))
		}
	}
	return resp
}

func LatestFollowing(c echo.Context) error {
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.LatestFollowing(ctx, &pb.LatestFollowingRequest{
		CurrentUid: GetCurrentUID(c),
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, BuildFollowingsRespSlice(svcresp.Followings))
}

func ListFriendship(c echo.Context) error {
	uid, err := GetUIDParam(c)
	if err != nil {
		return err
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
		CurrentUid:   GetCurrentUID(c),
		Uid:          uid,
		RelationType: params.RelationType,
		Paginator: &pb.PageQuick{
			NextId: params.PageQuickParams.NextID,
			Limit:  uint64(params.PageQuickParams.Limit),
		},
	})
	if err != nil {
		return err
	}

	return rest.BuildSuccessRespWithPagination(c, BuildFriendshipRespSlice(svcresp.Relations), svcresp.Paginator)
}

func Follow(c echo.Context) error {
	params := &FollowParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument
	}

	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}

	_, err = grpcsvc.Follow(ctx, &pb.FollowRequest{
		CurrentUid: GetCurrentUID(c),
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

	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}

	_, err = grpcsvc.UnFollow(ctx, &pb.UnFollowRequest{
		CurrentUid: GetCurrentUID(c),
		TargetUid:  params.ToUserID,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, nil)
}

func RecommendUser(c echo.Context) error {

	return rest.BuildSuccessResp(c, BuildFollowingsRespSlice(nil))

}
