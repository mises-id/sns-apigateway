package v1

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	"github.com/mises-id/sns-apigateway/lib/codes"

	pb "github.com/mises-id/sns-socialsvc/proto"
)

type (
	PageChannelUserInput struct {
		rest.PageParams
		Misesid string `json:"misesid" query:"misesid"`
	}
	GetChannelUserInput struct {
		Misesid string `json:"misesid" query:"misesid"`
	}

	GetChannelUserOutput struct {
		ID             string           `json:"id"`
		ChannelMisesid string           `json:"channel_misesid"`
		AirdropState   int32            `json:"airdrop_state"`
		ValidState     int32            `json:"valid_state"`
		User           *UserSummaryResp `json:"user"`
		CreatedAt      time.Time        `json:"created_at"`
	}
	PageChannelUserOutput struct {
		ID           string           `json:"id"`
		Channel_id   string           `json:"channel_id"`
		Amount       uint64           `json:"amount"`
		TxID         string           `json:"tx_id"`
		AirdropState int32            `json:"airdrop_state"`
		ValidState   int32            `json:"valid_state"`
		User         *UserSummaryResp `json:"user"`
		AirdropTime  time.Time        `json:"airdrop_time"`
		CreatedAt    time.Time        `json:"created_at"`
	}
	ChannelUrlRequest struct {
		Misesid string `json:"misesid" query:"misesid"`
		Type    string `json:"type" query:"type"`
		Medium  string `json:"medium" query:"medium"`
	}
)

func BuildChannelUserSliceResp(channel_users []*pb.ChannelUserInfo) []*PageChannelUserOutput {
	resp := []*PageChannelUserOutput{}
	for _, channel_user := range channel_users {
		resp = append(resp, BuildChannelUserResp(channel_user))
	}

	return resp
}

func BuildChannelUserResp(channel_user *pb.ChannelUserInfo) *PageChannelUserOutput {
	if channel_user == nil {
		return nil
	}
	resp := &PageChannelUserOutput{
		ID:           channel_user.Id,
		User:         BuildUserSummaryResp(channel_user.User),
		Channel_id:   channel_user.ChannelId,
		Amount:       channel_user.Amount,
		TxID:         channel_user.TxId,
		AirdropState: channel_user.AirdropState,
		ValidState:   channel_user.ValidState,
		AirdropTime:  time.Unix(int64(channel_user.AirdropTime), 0),
		CreatedAt:    time.Unix(int64(channel_user.CreatedAt), 0),
	}
	return resp
}
func BuildGetChannelUserResp(channel_user *pb.ChannelUserInfo) *GetChannelUserOutput {
	if channel_user == nil {
		return nil
	}
	resp := &GetChannelUserOutput{
		ID:             channel_user.Id,
		ChannelMisesid: channel_user.ChannelMisesid,
		User:           BuildUserSummaryResp(channel_user.User),
		AirdropState:   channel_user.AirdropState,
		ValidState:     channel_user.ValidState,
		CreatedAt:      time.Unix(int64(channel_user.CreatedAt), 0),
	}
	return resp
}

func PageChannelUser(c echo.Context) error {

	params := &PageChannelUserInput{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}

	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.PageChannelUser(ctx, &pb.PageChannelUserRequest{
		Misesid: params.Misesid,
		Paginator: &pb.Page{
			PageNum:  uint64(params.PageParams.PageNum),
			PageSize: uint64(params.PageParams.PageSize),
		},
	})
	if err != nil {
		return err
	}

	return rest.BuildSuccessRespWithPage(c, BuildChannelUserSliceResp(svcresp.ChannelUsers), svcresp.Paginator)
}

func GetChannelUser(c echo.Context) error {
	/* params := &GetChannelUserInput{}
	if err := c.Bind(params); err != nil {
		return err
	} */
	misesid := c.Param("misesid")
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.GetChannelUser(ctx, &pb.GetChannelUserRequest{
		Misesid: misesid,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, BuildGetChannelUserResp(svcresp.ChanelUser))
}
func ChannelInfo(c echo.Context) error {
	params := &ChannelUrlRequest{}
	if err := c.Bind(params); err != nil {
		return err
	}
	grpcsvc, ctx, err := rest.GrpcSocialService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.ChannelInfo(ctx, &pb.ChannelInfoRequest{
		Misesid: params.Misesid,
		Type:    params.Type,
		Medium:  params.Medium,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, echo.Map{
		"url":                svcresp.Url,
		"total_channel_user": svcresp.TotalChannelUser,
		"airdrop_amount":     svcresp.AirdropAmount,
		"medium_url":         svcresp.MediumUrl,
	})
}
