package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	MinNamePrice = sdk.Coins{sdk.NewInt64Coin("nametoken", 1)}
)

type Whois struct {
	Owner                  sdk.AccAddress `json:"owner"`                     // 拥有者
	Price                  sdk.Coins      `json:"price"`                     // 价格
	Value                  string         `json:"value"`                     // 值，可以是IP
	CanAuction             bool           `json:"can_auction"`               // 是否可拍卖
	AuctionUser            sdk.AccAddress `json:"auction_user"`              // 竞价用户
	AuctionPrice           sdk.Coins      `json:"auction_price"`             // 竞价价格
	AuctionDealBlockHeight int            `json:"auction_deal_block_height"` // 竞价成交区块高度
}

func NewWhois() Whois {
	return Whois{
		Price: MinNamePrice,
	}
}

func (object Whois) String() string {
	raw, _ := json.Marshal(&object)
	return string(raw)
}
