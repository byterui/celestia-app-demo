package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/gogo/protobuf/proto"
	gogotypes "github.com/gogo/protobuf/types"

	"github.com/tendermint/tendermint/libs/log"

	"celestia-app-demo/x/swap/types"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	memKey     storetypes.StoreKey
	paramSpace paramtypes.Subspace

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
) *Keeper {
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramSpace: ps,

		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) CreateModuleAccount(ctx sdk.Context) {
	moduleAcc := authtypes.NewEmptyModuleAccount(types.ModuleName, authtypes.Minter, authtypes.Burner)
	k.accountKeeper.SetModuleAccount(ctx, moduleAcc)
}

func (k Keeper) GetNextPairId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	nextPoolId := gogotypes.UInt64Value{}

	b := store.Get(types.KeyNextGlobalPoolId)
	if b == nil {
		panic(fmt.Errorf("getting at key (%v) should not have been nil", types.KeyNextGlobalPoolId))
	}
	if err := proto.Unmarshal(b, &nextPoolId); err != nil {
		panic(err)
	}

	return nextPoolId.Value
}

func (k Keeper) SetNextPairId(ctx sdk.Context, poolId uint64) {
	store := ctx.KVStore(k.storeKey)
	bz, err := proto.Marshal(&gogotypes.UInt64Value{Value: poolId})
	if err != nil {
		panic(err)
	}

	store.Set(types.KeyNextGlobalPoolId, bz)
}
