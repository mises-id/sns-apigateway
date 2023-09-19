package v1

import (
	"github.com/labstack/echo/v4"
	miningsvc "github.com/mises-id/mises-miningsvc/proto"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	"github.com/mises-id/sns-apigateway/lib/codes"
)

type MBAirdropUser struct {
	Misesid             string `json:"misesid"`
	ReceiveAddress      string `json:"receive_address"`
	TotalAirdropLimit   uint64 `json:"total_airdrop_limit"`
	CurrentAirdropLimit uint64 `json:"current_airdrop_limit"`
	CurrentAirdrop      uint64 `json:"current_airdrop"`
}

type MBAirdropClaimParams struct {
	//Misesid        string `json:"misesid" query:"misesid"`
	Pubkey         string `json:"pubkey" query:"pubkey"`
	ReceiveAddress string `json:"receive_address" query:"receive_address"`
	Nonce          string `json:"nonce" query:"nonce"`
	Sig            string `json:"sig" query:"sig"`
	TxHash         string `json:"tx_hash" query:"tx_hash"`
}

func ClaimMBAirdrop(c echo.Context) (err error) {

	misesid := GetCurrentMisesID(c)
	params := &MBAirdropClaimParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}
	grpcsvc, ctx, err := rest.GrpcMiningService()
	if err != nil {
		return err
	}
	_, err = grpcsvc.ClaimAirdrop(ctx, &miningsvc.ClaimAirdropRequest{
		Misesid:        misesid,
		Pubkey:         params.Pubkey,
		Nonce:          params.Nonce,
		ReceiveAddress: params.ReceiveAddress,
		Sig:            params.Sig,
		TxHash:         params.TxHash,
	})
	if err != nil {
		return
	}

	return rest.BuildSuccessResp(c, nil)
}

func FindMBAirdropUser(c echo.Context) error {
	misesid := c.Param("misesid")
	grpcsvc, ctx, err := rest.GrpcMiningService()
	if err != nil {
		return err
	}
	svcresp, err := grpcsvc.FindAirdropUser(ctx, &miningsvc.FindAirdropUserRequest{
		Misesid: misesid,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, buildMBAirdropUser(svcresp))
}

func buildMBAirdropUser(in *miningsvc.FindAirdropUserResponse) *MBAirdropUser {

	if in == nil || in.User == nil {
		return nil
	}
	resp := &MBAirdropUser{
		Misesid:             in.User.Misesid,
		ReceiveAddress:      in.User.ReceiveAddress,
		TotalAirdropLimit:   in.User.TotalAirdropLimit,
		CurrentAirdropLimit: in.User.CurrentAirdropLimit,
		CurrentAirdrop:      in.User.CurrentAirdrop,
	}

	return resp
}
