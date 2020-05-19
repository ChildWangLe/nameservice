package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/neo.wang/nameservice/x/nameservice/types"
)

// 持久层
type Keeper struct {
	CoinKeeper types.BankKeeper
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
}

// 工厂方法
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		CoinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

// 名字是否存在
func (object Keeper) IsNamePresent(ctx sdk.Context, name string) bool {
	if 0 >= len(name) {
		return false
	}
	store := ctx.KVStore(object.storeKey)
	return store.Has([]byte(name))
}

// 设置Whois
func (object Keeper) SetWhois(ctx sdk.Context, name string, whois types.Whois) {
	if whois.Owner.Empty() {
		return
	}
	store := ctx.KVStore(object.storeKey)
	store.Set([]byte(name), object.cdc.MustMarshalBinaryBare(whois))
}

// 获取Whois
func (object Keeper) GetWhois(ctx sdk.Context, name string) types.Whois {
	if !object.IsNamePresent(ctx, name) {
		return types.Whois{}
	}
	store := ctx.KVStore(object.storeKey)
	var whois types.Whois
	bz := store.Get([]byte(name))
	object.cdc.MustUnmarshalBinaryBare(bz, &whois)
	return whois
}

// 删除Whois
func (object Keeper) DeleteWhois(ctx sdk.Context, name string) {
	store := ctx.KVStore(object.storeKey)
	store.Delete([]byte(name))
}

// 决定Name
func (object Keeper) ResolveName(ctx sdk.Context, name string) string {
	return object.GetWhois(ctx, name).Value
}

// 设置Name
func (object Keeper) SetName(ctx sdk.Context, name, value string) {
	whois := object.GetWhois(ctx, name)
	whois.Value = value
	object.SetWhois(ctx, name, whois)
}

// 是否拥有Owner
func (object Keeper) HasOwner(ctx sdk.Context, name string) bool {
	return !object.GetWhois(ctx, name).Owner.Empty()
}

// 获取Owner
func (object Keeper) GetOwner(ctx sdk.Context, name string) sdk.AccAddress {
	return object.GetWhois(ctx, name).Owner
}

// 设置Owner
func (object Keeper) SetOwner(ctx sdk.Context, name string, owner sdk.AccAddress) {
	whois := object.GetWhois(ctx, name)
	whois.Owner = owner
	object.SetWhois(ctx, name, whois)
}

// 获取Price
func (object Keeper) GetPrice(ctx sdk.Context, name string) sdk.Coins {
	return object.GetWhois(ctx, name).Price
}

// 设置Price
func (object Keeper) SetPrice(ctx sdk.Context, name string, price sdk.Coins) {
	whois := object.GetWhois(ctx, name)
	whois.Price = price
	object.SetWhois(ctx, name, whois)
}

// 设置是否可以竞拍
func (object Keeper) SetCanAuction(ctx sdk.Context, name string, canAuction bool) {
	var arr []types.Whois
	for it := object.GetNamesIterator(ctx); it.Valid(); it.Next() {
		theName := string(it.Key())
		if theName != name {
			continue
		}
		whois := object.GetWhois(ctx, name)
		whois.CanAuction = canAuction
		arr = append(arr, whois)
	}
	for _, e := range arr {
		object.SetWhois(ctx, name, e)
	}
}

// 是否拥有竞拍信息
func (object Keeper) HasAuctionInfo(ctx sdk.Context, name string) bool {
	whois := object.GetWhois(ctx, name)
	return whois.CanAuction && !whois.AuctionUser.Empty()
}

// 获取竞拍信息
func (object Keeper) GetAuctionInfo(ctx sdk.Context, name string) (user sdk.AccAddress, price sdk.Coins, dealHeight int) {
	whois := object.GetWhois(ctx, name)
	user = whois.AuctionUser
	price = whois.AuctionPrice
	dealHeight = whois.AuctionDealBlockHeight
	return
}

// 设置竞拍信息
func (object Keeper) SetAuction(ctx sdk.Context, name string, user sdk.AccAddress, price sdk.Coins, dealHeight int) {
	whois := object.GetWhois(ctx, name)
	if !whois.Owner.Empty() {
		whois.AuctionUser = user
		whois.AuctionPrice = price
		whois.AuctionDealBlockHeight = dealHeight
		object.SetWhois(ctx, name, whois)
	}
}

// 获取竞拍倒计时
func (object Keeper) GetAuctionCountdown(ctx sdk.Context, name string) int {
	whois := object.GetWhois(ctx, name)
	if whois.Owner.Empty() || whois.AuctionUser.Empty() {
		return -1
	}
	return whois.AuctionDealBlockHeight
}

// 设置竞拍倒计时
func (object Keeper) SetAuctionCountdown(ctx sdk.Context, name string, dealHeight int) {
	whois := object.GetWhois(ctx, name)
	if whois.Owner.Empty() || whois.AuctionUser.Empty() {
		return
	}
	whois.AuctionDealBlockHeight = dealHeight
	object.SetWhois(ctx, name, whois)
}

// 倒计时所有竞拍
func (object Keeper) CountdownAllCanAuction(ctx sdk.Context,
	currHeight int64,
	callback func(m map[string]types.Whois)) {
	m := make(map[string]types.Whois, 0)
	for it := object.GetNamesIterator(ctx); it.Valid(); it.Next() {
		name := string(it.Key())
		whois := object.GetWhois(ctx, name)
		if !whois.CanAuction ||
			0 >= whois.AuctionDealBlockHeight ||
			int64(whois.AuctionDealBlockHeight) > currHeight {
			continue
		}
		m[name] = whois
	}
	if 0 < len(m) {
		callback(m)
	}
}

// 获取名字迭代器
func (object Keeper) GetNamesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(object.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte{})
}
