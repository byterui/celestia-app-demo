package keeper

import (
	context "context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"celestia-app-demo/x/swap/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// CreatePair implements types.MsgServer
func (m msgServer) CreatePair(goCtx context.Context, msg *types.MsgCreatePair) (*types.MsgCreateaPairResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	pair, err := m.Keeper.createPair(ctx, msg.Token0, msg.Token1)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateaPairResponse{
		NewLpToken: pair.LpToken.Denom,
	}, nil
}
