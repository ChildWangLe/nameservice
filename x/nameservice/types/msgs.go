package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// 消息基础结构
type MsgBase struct {
	Name string `json:"name"`
}

func (object MsgBase) Route() string {
	return RouterKey
}

// 设置Whois的Value
type MsgSetName struct {
	MsgBase
	Value string         `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
}

func NewMsgSetName(name, value string, owner sdk.AccAddress) MsgSetName {
	return MsgSetName{
		MsgBase: MsgBase{Name: name},
		Value:   value,
		Owner:   owner,
	}
}

func (object MsgSetName) Type() string {
	return "set_name"
}

func (object MsgSetName) ValidateBasic() error {
	if object.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, object.Owner.String())
	}
	if 0 >= len(object.Name) || 0 >= len(object.Value) {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Name and/or Value cannot be empty")
	}
	return nil
}

func (object MsgSetName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(object))
}

func (object MsgSetName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{object.Owner}
}

//买Whois
type MsgBuyName struct {
	MsgBase
	Bid   sdk.Coins      `json:"bid"`
	Buyer sdk.AccAddress `json:"buyer"`
}

func NewMsgBuyName(name string, bid sdk.Coins, buyer sdk.AccAddress) MsgBuyName {
	return MsgBuyName{
		MsgBase: MsgBase{Name: name},
		Bid:     bid,
		Buyer:   buyer,
	}
}

func (object MsgBuyName) Type() string {
	return "buy_name"
}

func (object MsgBuyName) ValidateBasic() error {
	if object.Buyer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, object.Buyer.String())
	}
	if 0 >= len(object.Name) {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Name cannot be empty")
	}
	if !object.Bid.IsAllPositive() {
		return sdkerrors.ErrInsufficientFunds
	}
	return nil
}

func (object MsgBuyName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(object))
}

func (object MsgBuyName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{object.Buyer}
}

// 删除Whois
type MsgDeleteName struct {
	MsgBase
	Owner sdk.AccAddress `json:"owner"`
}

func NewMsgDeleteName(name string, owner sdk.AccAddress) MsgDeleteName {
	return MsgDeleteName{
		MsgBase: MsgBase{Name: name},
		Owner:   owner,
	}
}

func (object MsgDeleteName) Type() string {
	return "delete_name"
}

func (object MsgDeleteName) ValidateBasic() error {
	if object.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, object.Owner.String())
	}
	if 0 >= len(object.Name) {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Name cannot be empty")
	}
	return nil
}

func (object MsgDeleteName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(object))
}

func (object MsgDeleteName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{object.Owner}
}

// 设置竞拍
type MsgSetCanAuction struct {
	MsgBase
	Owner      sdk.AccAddress `json:"owner"`
	CanAuction bool           `json:"can_auction"`
}

func NewMsgSetCanAuction(name string, owner sdk.AccAddress, auction bool) MsgSetCanAuction {
	return MsgSetCanAuction{
		MsgBase:    MsgBase{Name: name},
		Owner:      owner,
		CanAuction: auction,
	}
}

func (object MsgSetCanAuction) Type() string {
	return "set_can_auction"
}

func (object MsgSetCanAuction) ValidateBasic() error {
	if object.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, object.Owner.String())
	}
	if 0 >= len(object.Name) {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Name cannot be empty")
	}
	return nil
}

func (object MsgSetCanAuction) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(object))
}

func (object MsgSetCanAuction) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{object.Owner}
}

// 报价
type MsgOffset struct {
	MsgBase
	Bid   sdk.Coins      `json:"bid"`
	Buyer sdk.AccAddress `json:"buyer"`
}

func NewMsgOffset(name string, bid sdk.Coins, buyer sdk.AccAddress) MsgOffset {
	return MsgOffset{
		MsgBase: MsgBase{Name: name},
		Bid:     bid,
		Buyer:   buyer,
	}
}

func (object MsgOffset) Type() string {
	return "offset"
}

func (object MsgOffset) ValidateBasic() error {
	if object.Buyer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, object.Buyer.String())
	}
	if 0 >= len(object.Name) {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Name cannot be empty")
	}
	if !object.Bid.IsAllPositive() {
		return sdkerrors.ErrInsufficientFunds
	}
	return nil
}

func (object MsgOffset) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(object))
}

func (object MsgOffset) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{object.Buyer}
}
