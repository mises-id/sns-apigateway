package v1

import (
	"github.com/labstack/echo/v4"
	miningsvc "github.com/mises-id/mises-miningsvc/proto"
	"github.com/mises-id/sns-apigateway/app/apis/rest"
	"github.com/mises-id/sns-apigateway/lib/codes"
)

type AdMiningConfig struct {
	LimitPerDay uint32 `json:"limit_per_day"`
}

type MiningBonusConfig struct {
	BonusToMbRate        float32 `json:"bonus_to_mb_rate"`
	MinRedeemBonusAmount float64 `json:"min_redeem_bonus_amount"`
}

type MBAirdropConfig struct {
	MinReDeemMisAmount float64 `json:"min_redeem_mis_amount"`
	MisRedeemMBFee     float64 `json:"mis_redeem_mb_fee"`
}

type MiningConfigResponse struct {
	Bonus     *MiningBonusConfig `json:"bonus"`
	AdMining  *AdMiningConfig    `json:"ad_mining"`
	MBAirdrop *MBAirdropConfig   `json:"mb_airdrop"`
}

type RedeemBonusRequest struct {
	Bonus float64 `json:"bonus" query:"bonus"`
}

type BonusResponse struct {
	Bonus float64 `json:"bonus"`
}

func GeMiningConfig(c echo.Context) error {

	ethAddress := GetCurrentEthAddress(c)

	grpcsvc, ctx, err := rest.GrpcMiningService()
	if err != nil {
		return err
	}
	resp, err := grpcsvc.GetMiningConfig(ctx, &miningsvc.GetMiningConfigRequest{
		EthAddress: ethAddress,
	})
	if err != nil {
		return err
	}

	return rest.BuildSuccessResp(c, buildMiningConfigResponse(resp))
}

func buildMiningConfigResponse(in *miningsvc.GetMiningConfigResponse) (resp *MiningConfigResponse) {

	if in == nil {
		return
	}

	resp = &MiningConfigResponse{
		MBAirdrop: &MBAirdropConfig{
			MinReDeemMisAmount: in.MinRedeemMisAmount,
			MisRedeemMBFee:     in.MisRedeemMbFee,
		},
	}
	if in.AdMining != nil {
		resp.AdMining = &AdMiningConfig{
			LimitPerDay: in.AdMining.LimitPerDay,
		}
	}
	if in.Bonus != nil {
		resp.Bonus = &MiningBonusConfig{
			BonusToMbRate:        in.Bonus.BonusToMbRate,
			MinRedeemBonusAmount: in.Bonus.MinRedeemBonusAmount,
		}
	}

	return
}

func GetBonus(c echo.Context) error {

	ethAddress := GetCurrentEthAddress(c)

	grpcsvc, ctx, err := rest.GrpcMiningService()
	if err != nil {
		return err
	}
	resp, err := grpcsvc.GetBonus(ctx, &miningsvc.GetBonusRequest{
		EthAddress: ethAddress,
	})
	if err != nil {
		return err
	}

	return rest.BuildSuccessResp(c, buildBonusResponse(resp))
}

func RedeemBonus(c echo.Context) error {

	ethAddress := GetCurrentEthAddress(c)
	params := &RedeemBonusRequest{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}

	grpcsvc, ctx, err := rest.GrpcMiningService()
	if err != nil {
		return err
	}
	_, err = grpcsvc.RedeemBonus(ctx, &miningsvc.RedeemBonusRequest{
		EthAddress:     ethAddress,
		ReceiveAddress: ethAddress,
		Bonus:          params.Bonus,
	})
	if err != nil {
		return err
	}

	return rest.BuildSuccessResp(c, nil)
}

func buildBonusResponse(in *miningsvc.GetBonusResponse) (resp *BonusResponse) {

	if in == nil {
		return
	}
	resp = &BonusResponse{
		Bonus: in.Bonus,
	}

	return
}
