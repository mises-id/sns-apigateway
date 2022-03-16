package v1

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mises-id/sns-apigateway/app/apis/rest"

	cosmosrest "github.com/cosmos/cosmos-sdk/types/rest"
	tmjson "github.com/tendermint/tendermint/libs/json"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	types "github.com/tendermint/tendermint/rpc/jsonrpc/types"
)

type GasPricesResp struct {
	SafeGasPrice    float64 `json:"safe_gasprice"`
	ProposeGasPrice float64 `json:"propose_gasprice"`
	FastGasPrice    float64 `json:"fast_gasprice"`
}

type ChainInfoResp struct {
	BlockHeight int64    `json:"block_height"`
	BlockHash   string   `json:"block_hash"`
	ChainID     string   `json:"chain_id"`
	TrustNodes  []string `json:"trust_nodes"`
}

func GasPrices(c echo.Context) error {

	return rest.BuildSuccessResp(c, &GasPricesResp{
		SafeGasPrice:    0.001,
		ProposeGasPrice: 0.001,
		FastGasPrice:    0.001,
	})
}

type LastBlockInfo struct {
	blockTime time.Time
	chainID   string
	hash      string
	height    int64
}

func ChainInfo(c echo.Context) error {

	store := rest.InMemoryStore()
	var info *LastBlockInfo
	if cached, exist := store.Load("LastBlockInfo"); exist {
		blockCache := cached.(*LastBlockInfo)
		if blockCache != nil && time.Since(blockCache.blockTime) < time.Hour {
			info = blockCache
		}

	}
	if info == nil {
		resp, err := cosmosrest.GetRequest(fmt.Sprintf("%s/block", "http://127.0.0.1:26657"))
		if err != nil {
			return err
		}
		rpcResponse := &types.RPCResponse{}
		err = tmjson.Unmarshal(
			resp,
			rpcResponse,
		)
		if err != nil {
			return err
		}
		resultBlock := &ctypes.ResultBlock{}
		err = tmjson.Unmarshal(
			rpcResponse.Result,
			resultBlock,
		)
		if err != nil {
			return err
		}
		info = &LastBlockInfo{
			hash:      resultBlock.BlockID.Hash.String(),
			height:    resultBlock.Block.Height,
			chainID:   resultBlock.Block.ChainID,
			blockTime: time.Now(),
		}
		store.Store("LastBlockInfo", info)
	}

	return rest.BuildSuccessResp(c, &ChainInfoResp{
		BlockHeight: info.height,
		BlockHash:   info.hash,
		ChainID:     info.chainID,
		TrustNodes: []string{
			"http://e1.mises.site:26657",
			"http://e2.mises.site:26657",
			"http://w1.mises.site:26657",
			"http://w2.mises.site:26657",
		},
	})
}
