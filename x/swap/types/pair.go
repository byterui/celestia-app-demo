package types

import (
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	ModuleLpPrefix = "swap"
)

func CreatePair(poolId uint64, token0, token1 string) (*Pair, error) {
	sortToken0, sortToken1 := SortToken(token0, token1)
	lp_token, err := GetPairLpToken(sortToken0, sortToken1)
	if err != nil {
		return nil, err
	}

	return &Pair{
		Account: NewPoolAddress(poolId).String(),
		Token0: sdk.Coin{
			Denom:  token0,
			Amount: math.Int{},
		},
		Token1: sdk.Coin{
			Denom:  token1,
			Amount: math.Int{},
		},
		LpToken: sdk.Coin{
			Denom:  lp_token,
			Amount: math.Int{},
		},
	}, nil
}

// GetPairLpToken constructs a lp token
// The pair lp token constructed is swap/{token0}/{token1}
func GetPairLpToken(token0, token1 string) (string, error) {
	lptoken := strings.Join([]string{ModuleLpPrefix, token0, token1}, "/")
	return lptoken, sdk.ValidateDenom(lptoken)
}

func SortToken(token0, token1 string) (string, string) {
	if token0 > token1 {
		return token1, token0
	}
	return token0, token1
}

// DeconstructLpToken deconstruct a lp token to tokens in it.
func DeconstructLpToken(lptoken string) (string, string, error) {
	err := sdk.ValidateDenom(lptoken)
	if err != nil {
		return "", "", err
	}

	strParts := strings.Split(lptoken, "/")
	if len(strParts) < 3 {
		return "", "", sdkerrors.Wrapf(ErrInvalidToken, "not enough parts of lptoken %s", lptoken)
	}

	err = sdk.ValidateDenom(strParts[1])
	if err != nil {
		return "", "", err
	}

	err = sdk.ValidateDenom(strParts[2])
	if err != nil {
		return "", "", err
	}

	return strParts[1], strParts[1], nil
}

func NewPoolAddress(poolId uint64) sdk.AccAddress {
	key := append([]byte("pair"), sdk.Uint64ToBigEndian(poolId)...)
	return address.Module(ModuleName, key)
}
