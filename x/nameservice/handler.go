package nameservice

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/neo.wang/nameservice/x/nameservice/keeper"
	"github.com/neo.wang/nameservice/x/nameservice/types"
)

func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgSetName:
			return handleMsgSetName(ctx, keeper, msg)
		case types.MsgBuyName:
			return handleMsgBuyName(ctx, keeper, msg)
		case types.MsgDeleteName:
			return handleMsgDeleteName(ctx, keeper, msg)
		case types.MsgSetCanAuction:
			return handleMsgSetCanAuction(ctx, keeper, msg)
		case types.MsgOffset:
			return handleMsgOffset(ctx, keeper, msg)
		}
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized NameService Msg type: %v", msg))
	}
}

// 处理设置Name
func handleMsgSetName(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgSetName) (*sdk.Result, error) {
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
	}
	keeper.SetName(ctx, msg.Name, msg.Value)
	return &sdk.Result{}, nil
}

// 处理买Name
func handleMsgBuyName(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgBuyName) (*sdk.Result, error) {
	if keeper.GetPrice(ctx, msg.Name).IsAllGT(msg.Bid) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "Bid not high enough")
	}
	if keeper.HasOwner(ctx, msg.Name) {
		if keeper.GetOwner(ctx, msg.Name).Equals(msg.Buyer) {
			return &sdk.Result{}, nil
		}
		err := keeper.CoinKeeper.SendCoins(ctx, msg.Buyer, keeper.GetOwner(ctx, msg.Name), msg.Bid)
		if nil != err {
			return nil, err
		}
	} else {
		_, err := keeper.CoinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Bid)
		if nil != err {
			return nil, err
		}
	}
	keeper.SetOwner(ctx, msg.Name, msg.Buyer)
	keeper.SetPrice(ctx, msg.Name, msg.Bid)
	return &sdk.Result{}, nil
}

// 处理删除Name
func handleMsgDeleteName(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgDeleteName) (*sdk.Result, error) {
	if !keeper.IsNamePresent(ctx, msg.Name) {
		return nil, sdkerrors.Wrap(types.ErrNameDoesNotExists, msg.Name)
	}
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
	}
	keeper.DeleteWhois(ctx, msg.Name)
	return &sdk.Result{}, nil
}

// 处理开启｜关闭竞拍
func handleMsgSetCanAuction(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgSetCanAuction) (*sdk.Result, error) {
	if !keeper.IsNamePresent(ctx, msg.Name) {
		return nil, sdkerrors.Wrap(types.ErrNameDoesNotExists, msg.Name)
	}
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
	}
	keeper.SetCanAuction(ctx, msg.Name, msg.CanAuction)
	return &sdk.Result{}, nil
}

// 处理竞价
func handleMsgOffset(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgOffset) (*sdk.Result, error) {
	// 买家为空
	if msg.Buyer.Empty() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Buyer cannot be empty")
	}
	// Name不存在
	if !keeper.IsNamePresent(ctx, msg.Name) {
		return nil, sdkerrors.Wrap(types.ErrNameDoesNotExists, msg.Name)
	}
	// 判断竞价是否高于
	_, price, _ := keeper.GetAuctionInfo(ctx, msg.Name)
	// 已有的竞价比本次出价高
	if nil != price && price.IsAllGTE(msg.Bid) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "Bid not high enough")
	}
	// 设置新的竞价信息
	keeper.SetAuction(ctx, msg.Name, msg.Buyer, msg.Bid, int(ctx.BlockHeight()+types.DefaultAuctionCountdown))
	return &sdk.Result{}, nil
}
