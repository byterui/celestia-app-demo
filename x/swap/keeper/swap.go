package keeper

import (
	"celestia-app-demo/x/swap/types"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) swapExactTokensForTokens(
	ctx sdk.Context,
	sender sdk.AccAddress,
	amountIn sdk.Int,
	amountOutMin sdk.Int,
	path []string,
	recipient sdk.AccAddress,
) ([]sdk.Int, error) {
	amounts, err := k.getAmountOutByPath(ctx, amountIn, path)
	if err != nil {
		return nil, err
	}

	if amounts[len(amounts)-1].LT(amountOutMin) {
		return nil, types.ErrInsufficientTargetAmount
	}

	firstPairAccount, err := k.getPairAccountFromTokens(ctx, path[0], path[1])
	if err != nil {
		return nil, err
	}

	err = k.bankKeeper.SendCoins(ctx, sender, firstPairAccount, sdk.Coins{sdk.NewCoin(path[0], amounts[0])})
	if err != nil {
		return nil, err
	}

	err = k.swap(ctx, amounts, path, recipient)
	if err != nil {
		return nil, err
	}

	return amounts, nil
}

func (k Keeper) getAmountOutByPath(
	ctx sdk.Context,
	amountIn sdk.Int,
	path []string,
) ([]sdk.Int, error) {
	if len(path) <= 1 {
		return nil, types.ErrInvalidPath
	}

	pathLastIndex := len(path) - 1
	outVec := []sdk.Int{amountIn}

	for i := 0; i < pathLastIndex; i++ {
		pairAccount, err := k.getPairAccountFromTokens(ctx, path[i], path[i+1])
		if pairAccount == nil {
			return nil, types.ErrPairNotExist
		}

		reserve0 := k.bankKeeper.GetBalance(ctx, pairAccount, path[i])
		reserve1 := k.bankKeeper.GetBalance(ctx, pairAccount, path[i+1])

		if reserve0.IsZero() || reserve1.IsZero() {
			return nil, types.ErrInvalidPath
		}

		amount, err := getAmountOut(outVec[i], reserve0.Amount, reserve1.Amount)
		if amount.IsZero() || err != nil {
			return nil, types.ErrInvalidPath
		}

		// todo check k

		outVec = append(outVec, amount)
	}

	return outVec, nil
}

func (k Keeper) getPairIdFromUnsortTokens(ctx sdk.Context, token0 string, token1 string) (uint64, error) {
	sortedToken0, sortedToken1 := types.SortToken(token0, token1)
	return k.GetPoolIdFromTokens(ctx, sortedToken0, sortedToken1)
}

func getAmountOut(amountIn, reserveIn, reserveOut sdk.Int) (sdk.Int, error) {
	if amountIn.IsZero() || reserveIn.IsZero() || reserveOut.IsZero() {
		return sdk.ZeroInt(), types.ErrMath
	}

	inputAmountWithFee := amountIn.ModRaw(997)

	numerator := inputAmountWithFee.Mul(reserveOut)

	denominator := reserveIn.MulRaw(1000).Add(inputAmountWithFee)

	amountOut := big.NewInt(0).Div(numerator.BigInt(), denominator.BigInt())

	return sdk.NewIntFromBigInt(amountOut), nil
}

func (k Keeper) getPairAccountFromTokens(ctx sdk.Context, token0 string, token1 string) (sdk.AccAddress, error) {
	id, err := k.getPairIdFromUnsortTokens(ctx, token0, token1)
	if err != nil {
		return nil, nil
	}

	pair := k.GetPairFromId(ctx, id)
	if pair == nil {
		return nil, types.ErrPairNotExist
	}

	pairAccount := sdk.MustAccAddressFromBech32(pair.Account)
	return pairAccount, nil
}

func (k Keeper) swap(ctx sdk.Context, amounts []sdk.Int, path []string, recipient sdk.AccAddress) error {
	if len(amounts) != len(path) {
		return types.ErrInsufficientFunds
	}
	pathLen := len(path)

	for i := 0; i < pathLen-1; i++ {
		input := path[i]
		output := path[i+1]
		amount0Out := sdk.ZeroInt()
		amount1Out := amounts[i+1]

		token0, token1 := types.SortToken(input, output)

		if input != token0 {
			amount0Out = amounts[i+1]
			amount1Out = sdk.ZeroInt()
		}

		if i < pathLen-2 {
			midAccount, err := k.getPairAccountFromTokens(ctx, output, path[i+1])
			if err != nil {
				return err
			}

			err = k.pairSwap(ctx, token0, token1, amount0Out, amount1Out, midAccount)
			if err != nil {
				return err
			}
		} else {
			err := k.pairSwap(ctx, token0, token1, amount0Out, amount1Out, recipient)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (k Keeper) pairSwap(
	ctx sdk.Context,
	token0 string,
	token1 string,
	amount0 sdk.Int,
	amount1 sdk.Int,
	recipient sdk.AccAddress,
) error {
	pairAccount, err := k.getPairAccountFromTokens(ctx, token0, token1)
	if err != nil {
		return err
	}

	reserve0 := k.bankKeeper.GetBalance(ctx, pairAccount, token0)
	reserve1 := k.bankKeeper.GetBalance(ctx, pairAccount, token1)

	if amount0.GT(reserve0.Amount) || amount1.GT(reserve1.Amount) {
		return types.ErrInsufficientPairReserve
	}

	if !amount0.IsZero() {
		err := k.bankKeeper.SendCoins(ctx, pairAccount, recipient, sdk.Coins{sdk.NewCoin(token0, amount0)})
		if err != nil {
			return err
		}
	}
	if !amount1.IsZero() {
		err := k.bankKeeper.SendCoins(ctx, pairAccount, recipient, sdk.Coins{sdk.NewCoin(token1, amount1)})
		if err != nil {
			return err
		}
	}

	return nil
}
