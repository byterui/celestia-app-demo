package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/tokenfactory module sentinel errors
var (
	ErrInvalidToken = sdkerrors.Register(ModuleName, 1, "invalid token")
)
