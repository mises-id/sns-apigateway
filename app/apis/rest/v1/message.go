package v1

import (
	"github.com/labstack/echo"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	"github.com/mises-id/sns-apigateway/lib/codes"
	pb "github.com/mises-id/sns-socialsvc/proto"
)

type ListMessageParams struct {
	rest.PageQuickParams
}

type ReadMessageParams struct {
	rest.PageQuickParams
	IDs      []string `body:"ids"`
	LatestID string   `body:"latest_id"`
}

func BuildMessageRespSlice(in []*pb.Message) []*pb.Message {
	if in == nil {
		return []*pb.Message{}
	} else {
		return in
	}
}

func ListMessage(c echo.Context) error {
	params := &ListMessageParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.Newf("invalid query params")
	}
	var currentUID uint64
	if c.Get("CurrentUID") != nil {
		currentUID = c.Get("CurrentUID").(uint64)
	}
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.ListMessage(ctx, &pb.ListMessageRequest{
		CurrentUid: currentUID,
		Paginator: &pb.PageQuick{
			NextId: params.PageQuickParams.NextID,
			Limit:  uint64(params.PageQuickParams.Limit),
		},
	})
	if err != nil {
		return err
	}

	return rest.BuildSuccessRespWithPagination(c, BuildMessageRespSlice(svcresp.Messages), svcresp.Paginator)
}

func ReadMessage(c echo.Context) error {
	params := &ReadMessageParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.Newf("invalid query params")
	}
	var currentUID uint64
	if c.Get("CurrentUID") != nil {
		currentUID = c.Get("CurrentUID").(uint64)
	}

	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	_, err = grpcsvc.ReadMessage(ctx, &pb.ReadMessageRequest{
		CurrentUid: currentUID,
		LatestID:   params.LatestID,
		Ids:        params.IDs,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, nil)
}
