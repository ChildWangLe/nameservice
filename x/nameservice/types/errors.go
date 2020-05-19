package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrNameDoesNotExists = sdkerrors.Register(ModuleName, 1, "name does not exists")
)
