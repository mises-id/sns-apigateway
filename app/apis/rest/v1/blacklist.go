package v1

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	"github.com/mises-id/sns-apigateway/lib/codes"

	pb "github.com/mises-id/sns-socialsvc/proto"
)

type BlacklistParams struct {
	UID uint64 `json:"uid"`
}

type ListBlacklistParams struct {
	rest.PageQuickParams
}

type BlackResp struct {
	User      *UserSummaryResp `json:"user"`
	CreatedAt time.Time        `json:"created_at"`
}

func BuildBlackResp(in *pb.Blacklist) *BlackResp {
	return &BlackResp{
		User:      BuildUserSummaryResp(in.User),
		CreatedAt: time.Unix(int64(in.CreatedAt), 0),
	}
}

func BuildBlacklistRespSlice(infos []*pb.Blacklist) []*BlackResp {
	resp := []*BlackResp{}
	for _, info := range infos {
		resp = append(resp, BuildBlackResp(info))
	}

	return resp
}

func ListBlacklist(c echo.Context) error {
	params := &ListBlacklistParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.Newf("invalid query params")
	}
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.ListBlacklist(ctx, &pb.ListBlacklistRequest{
		Uid: GetCurrentUID(c),
		Paginator: &pb.PageQuick{
			NextId: params.PageQuickParams.NextID,
			Limit:  uint64(params.PageQuickParams.Limit),
		},
	},
	)
	if err != nil {
		return err
	}
	return rest.BuildSuccessRespWithPagination(c, BuildBlacklistRespSlice(svcresp.Blacklists), svcresp.Paginator)
}

func CreateBlacklist(c echo.Context) error {
	params := &BlacklistParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.Newf("invalid query params")
	}
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	_, err = grpcsvc.CreateBlacklist(ctx, &pb.CreateBlacklistRequest{
		Uid:       GetCurrentUID(c),
		TargetUid: params.UID,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, nil)
}

func DeleteBlacklist(c echo.Context) error {
	uid, err := GetUIDParam(c)
	if err != nil {
		return err
	}
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	_, err = grpcsvc.DeleteBlacklist(ctx, &pb.DeleteBlacklistRequest{
		Uid:       GetCurrentUID(c),
		TargetUid: uid,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, nil)
}
