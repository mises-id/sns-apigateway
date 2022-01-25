package v1

import (
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
	LatestMessage      *MessageResp `json:"latest_message"`
	Total              uint64       `json:"total"`
	NotificationsCount uint64       `json:"notifications_count"`
	UsersCount         uint64       `json:"users_count"`
}

type MessageMeta interface {
	isMessageMeta()
}

type NewLikeStatusMeta struct {
	UID            uint64 `json:"uid"`
	StatusID       string `json:"status_id"`
	StatusContent  string `json:"status_content"`
	StatusImageURL string `json:"status_image_url"`
}

func (NewLikeStatusMeta) isMessageMeta() {}

type NewLikeCommentMeta struct {
	UID             uint64 `json:"uid"`
	CommentID       string `json:"comment_id"`
	CommentUsername string `json:"comment_username"`
	CommentContent  string `json:"comment_content"`
}

func (NewLikeCommentMeta) isMessageMeta() {}

type NewCommentMeta struct {
	UID                  uint64 `json:"uid"`
	GroupID              string `json:"group_id"`
	CommentID            string `json:"comment_id"`
	Content              string `json:"content"`
	ParentContent        string `json:"parent_content"`
	ParentUsername       string `json:"parent_username"`
	StatusContentSummary string `json:"status_content_summary"`
	StatusImageURL       string `json:"status_image_url"`
}

func (NewCommentMeta) isMessageMeta() {}

type NewFanMeta struct {
	UID         uint64 `json:"uid"`
	FanUsername string `json:"fan_username"`
}

func (NewFanMeta) isMessageMeta() {}

type NewForwardMeta struct {
	UID            uint64 `json:"uid"`
	StatusID       string `json:"status_id"`
	ForwardContent string `json:"forward_content"`
	ContentSummary string `json:"content_summary"`
	ImageURL       string `json:"image_url"`
}

func (NewForwardMeta) isMessageMeta() {}

type MessageResp struct {
	ID               string           `json:"id"`
	User             *UserSummaryResp `json:"user"`
	MessageType      string           `json:"message_type"`
	MetaData         MessageMeta      `json:"meta_data"`
	State            string           `json:"state"`
	Status           *StatusResp      `json:"status"`
	StatusIsDeleted  bool             `json:"ststus_is_deleted"`
	CommentIsDeleted bool             `json:"comment_is_deleted"`
	CreatedAt        time.Time        `json:"created_at"`
}

func BuildMessageSummaryResp(in *pb.MessageSummary) *MessageSummaryResp {

	if in == nil {
		return &MessageSummaryResp{}
	} else {
		return &MessageSummaryResp{
			LatestMessage:      BuildMessageResp(in.LatestMessage),
			Total:              uint64(in.Total),
			NotificationsCount: uint64(in.NotificationsCount),
			UsersCount:         uint64(in.UsersCount),
		}
	}
}
func BuildMessageResp(in *pb.Message) *MessageResp {
	if in == nil {
		return &MessageResp{}
	} else {
		ret := &MessageResp{
			ID:               in.Id,
			User:             BuildUserSummaryResp(in.FromUser),
			MessageType:      in.MessageType,
			State:            in.State,
			StatusIsDeleted:  in.StatusIsDeleted,
			CommentIsDeleted: in.CommentIsDeleted,
			Status:           BuildStatusResp(in.Status),
			CreatedAt:        time.Unix(int64(in.CreatedAt), 0),
		}
		switch in.MessageType {
		case "new_like_status":
			ret.MetaData = &NewLikeStatusMeta{
				UID:            in.NewLikeStatusMeta.Uid,
				StatusID:       in.NewLikeStatusMeta.StatusId,
				StatusContent:  in.NewLikeStatusMeta.StatusContent,
				StatusImageURL: in.NewLikeStatusMeta.StatusImageUrl,
			}
		case "new_like_comment":
			ret.MetaData = &NewLikeCommentMeta{
				UID:             in.NewLikeCommentMeta.Uid,
				CommentID:       in.NewLikeCommentMeta.CommentId,
				CommentUsername: in.NewLikeCommentMeta.CommentUsername,
				CommentContent:  in.NewLikeCommentMeta.CommentContent,
			}
		case "new_comment":
			ret.MetaData = &NewCommentMeta{
				UID:                  in.NewCommentMeta.Uid,
				GroupID:              in.NewCommentMeta.GroupId,
				CommentID:            in.NewCommentMeta.CommentId,
				Content:              in.NewCommentMeta.Content,
				ParentContent:        in.NewCommentMeta.ParentContent,
				ParentUsername:       in.NewCommentMeta.ParentUserName,
				StatusContentSummary: in.NewCommentMeta.StatusContentSummary,
				StatusImageURL:       in.NewCommentMeta.StatusImageUrl,
			}
		case "new_fans":
			ret.MetaData = &NewFanMeta{
				UID:         in.NewFansMeta.Uid,
				FanUsername: in.NewFansMeta.FanUsername,
			}
		case "new_fowards":
			ret.MetaData = &NewForwardMeta{
				UID:            in.NewForwardMeta.Uid,
				StatusID:       in.NewForwardMeta.StatusId,
				ForwardContent: in.NewForwardMeta.ForwardContent,
				ContentSummary: in.NewForwardMeta.ContentSummary,
				ImageURL:       in.NewForwardMeta.ImageUrl,
			}
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
func MessageSummary(c echo.Context) error {
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.GetMessageSummary(ctx, &pb.GetMessageSummaryRequest{
		CurrentUid: GetCurrentUID(c),
	})
	if err != nil {
		return err
	}

	return rest.BuildSuccessResp(c, BuildMessageSummaryResp(svcresp.Summary))
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
