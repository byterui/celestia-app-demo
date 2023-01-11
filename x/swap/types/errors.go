package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidTokens           = sdkerrors.Register(ModuleName, 1, "invalid token")
	ErrPairNotExist            = sdkerrors.Register(ModuleName, 2, "pair not exist")
	ErrPairCreated             = sdkerrors.Register(ModuleName, 3, "pair already created")
	ErrInvalidTokenAmountRange = sdkerrors.Register(ModuleName, 4, "invalid token amount range")
	ErrInsufficientFunds       = sdkerrors.Register(ModuleName, 5, "insufficient token amount")
	ErrMath                    = sdkerrors.Register(ModuleName, 6, "error occor in calculate")
)
