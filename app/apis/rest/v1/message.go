package v1

import (
	"encoding/json"
	"time"

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
type MessageSummaryResp struct {
	LatestMessage      MessageResp `json:"latest_message"`
	Total              uint64      `json:"total"`
	NotificationsCount uint64      `json:"notifications_count"`
	UsersCount         uint64      `json:"users_count"`
}
type MessageResp struct {
	ID          string    `json:"id"`
	MessageType string    `json:"message_type"`
	MetaData    string    `json:"meta_data"`
	State       string    `json:"state"`
	CreatedAt   time.Time `json:"created_at"`
}

func BuildMessageResp(in *pb.Message) *MessageResp {
	if in == nil {
		return &MessageResp{}
	} else {

		ret := &MessageResp{
			ID:          in.Id,
			MessageType: in.MessageType,
			State:       in.State,
			CreatedAt:   time.Now(),
		}
		if meta, err := json.Marshal(in.MetaData); err == nil {
			ret.MetaData = string(meta)
		}
		return ret
	}
}

func BuildMessageRespSlice(ins []*pb.Message) []*MessageResp {
	resp := []*MessageResp{}
	if ins == nil {
		return []*MessageResp{}
	} else {
		for _, info := range ins {
			resp = append(resp, BuildMessageResp(info))
		}
	}
	return resp
}

func ListMessage(c echo.Context) error {
	params := &ListMessageParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.Newf("invalid query params")
	}
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.ListMessage(ctx, &pb.ListMessageRequest{
		CurrentUid: GetCurrentUID(c),
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
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	_, err = grpcsvc.ReadMessage(ctx, &pb.ReadMessageRequest{
		CurrentUid: GetCurrentUID(c),
		LatestID:   params.LatestID,
		Ids:        params.IDs,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, nil)
}

func MessageSummary(c echo.Context) error {

	return rest.BuildSuccessResp(c, &MessageSummaryResp{
		Total:              1,
		NotificationsCount: 2,
		UsersCount:         3,
	})
}
