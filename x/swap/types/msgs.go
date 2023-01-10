package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreatePair = "create_pair"
)

var _ sdk.Msg = &MsgCreatePair{}

func (m MsgCreatePair) Route() string { return RouterKey }

func (m MsgCreatePair) Type() string { return TypeMsgCreatePair }

func (m MsgCreatePair) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	if len(m.Token0) == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "token0:(%s) is nil", m.Token0)
	}

	if m.Token0 == m.Token1 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "token0 == token1")
	}

	return nil
}

func (m MsgCreatePair) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgCreatePair) GetSigners() []sdk.AccAddress {
	sender, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{sender}
}
