package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"celestia-app-demo/x/swap/types"
)

func (k Keeper) createPair(ctx sdk.Context, token0, token1 string) (*types.Pair, error) {
	if !k.bankKeeper.HasSupply(ctx, token0) {
		return nil, fmt.Errorf("%s has't supply", token0)
	}

	if !k.bankKeeper.HasSupply(ctx, token1) {
		return nil, fmt.Errorf("%s has't supply", token1)
	}

	poolId := k.getNextPoolIdAndIncrement(ctx)

	pair, err := types.CreatePair(poolId, token0, token1)
	if err != nil {
		return nil, err
	}

	bz, err := k.cdc.Marshal(pair)
	if err != nil {
		return nil, err
	}
	store := ctx.KVStore(k.storeKey)

	poolKey := types.GetKeyPrefixPools(poolId)
	store.Set(poolKey, bz)

	k.bankKeeper.SetDenomMetaData(ctx, banktypes.Metadata{
		Description: fmt.Sprintf("the lp token of the swap pair%d", poolId),
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    pair.LpToken.Denom,
				Exponent: 0,
			},
		},
		Base:    pair.LpToken.Denom,
		Display: pair.LpToken.Denom,
	})

	return pair, nil
}

func (k Keeper) getNextPoolIdAndIncrement(ctx sdk.Context) uint64 {
	nextPoolId := k.GetNextPairId(ctx)
	k.SetNextPairId(ctx, nextPoolId+1)
	return nextPoolId
}
