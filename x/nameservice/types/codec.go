package types

import "github.com/cosmos/cosmos-sdk/codec"

var (
	ModuleCdc = codec.New()
)

func init() {
	RegisterCodec(ModuleCdc)
}

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgBase{}, "NameService/MsgBase", nil)
	cdc.RegisterConcrete(MsgSetName{}, "NameService/SetName", nil)
	cdc.RegisterConcrete(MsgBuyName{}, "NameService/BuyName", nil)
	cdc.RegisterConcrete(MsgDeleteName{}, "NameService/DeleteName", nil)
	cdc.RegisterConcrete(MsgSetCanAuction{}, "NameService/SetCanAuction", nil)
	cdc.RegisterConcrete(MsgOffset{}, "NameService/Offset", nil)
}
