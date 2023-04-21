package keeper

import (
	"fmt"

	sdk "github.com/line/lbm-sdk/types"
	wasmkeeper "github.com/line/wasmd/x/wasm/keeper"
	"github.com/line/wasmd/x/wasm/types"
	wasmplustypes "github.com/line/wasmd/x/wasmplus/types"
	wasmvm "github.com/line/wasmvm"
)

type cosmwasmAPIImpl struct {
	keeper *Keeper
	ctx    *sdk.Context
}

func (a cosmwasmAPIImpl) callCallablePoint(contractAddrStr string, name []byte, args []byte, isReadonly bool, callstack []byte, gasLimit uint64) ([]byte, uint64, error) {
	contractAddr := sdk.MustAccAddressFromBech32(contractAddrStr)
	contractInfo, codeInfo, prefixStore, err := a.keeper.ContractInstance(*a.ctx, contractAddr)
	if err != nil {
		return nil, 0, err
	}

	if a.keeper.IsInactiveContract(*a.ctx, contractAddr) {
		return nil, 0, fmt.Errorf("called contract cannot be executed")
	}

	env := types.NewEnv(*a.ctx, contractAddr)
	wasmStore := wasmplustypes.NewWasmStore(prefixStore)
	gasRegister := a.keeper.GetGasRegister()
	querier := wasmkeeper.NewQueryHandler(*a.ctx, a.keeper.GetWasmVMQueryHandler(), contractAddr, gasRegister)
	gasMeter := a.keeper.GasMeter(*a.ctx)
	api := a.keeper.CosmwasmAPI(*a.ctx)

	instantiateCost := gasRegister.ToWasmVMGas(gasRegister.InstantiateContractCosts(a.keeper.IsPinnedCode(*a.ctx, contractInfo.CodeID), len(args)))
	if gasLimit < instantiateCost {
		return nil, 0, fmt.Errorf("lack of gas for calling callable point")
	}
	wasmGasLimit := gasLimit - instantiateCost

	result, events, attrs, gas, err := a.keeper.GetWasmVM().CallCallablePoint(name, codeInfo.CodeHash, isReadonly, callstack, env, args, wasmStore, api, querier, gasMeter, wasmGasLimit, wasmkeeper.CostJSONDeserialization)
	gas += instantiateCost

	if err != nil {
		return nil, gas, err
	}

	if !isReadonly {
		// issue events and attrs
		if len(attrs) != 0 {
			eventsByAttr, err := newCallablePointEvent(attrs, contractAddr, callstack)
			if err != nil {
				return nil, gas, err
			}
			a.ctx.EventManager().EmitEvents(eventsByAttr)
		}

		if len(events) != 0 {
			customEvents, err := newCustomCallablePointEvents(events, contractAddr, callstack)
			if err != nil {
				return nil, gas, err
			}
			a.ctx.EventManager().EmitEvents(customEvents)
		}
	}

	return result, gas, err
}

// returns result, gas used, error
func (a cosmwasmAPIImpl) validateInterface(contractAddrStr string, expectedInterface []byte) ([]byte, uint64, error) {
	contractAddr := sdk.MustAccAddressFromBech32(contractAddrStr)

	if a.keeper.IsInactiveContract(*a.ctx, contractAddr) {
		return nil, 0, fmt.Errorf("try to validate a contract cannot be executed")
	}

	_, codeInfo, _, err := a.keeper.ContractInstance(*a.ctx, contractAddr)
	if err != nil {
		return nil, 0, err
	}

	result, err := a.keeper.GetWasmVM().ValidateDynamicLinkInterface(codeInfo.CodeHash, expectedInterface)

	return result, 0, err
}

func (k *Keeper) CosmwasmAPI(ctx sdk.Context) wasmvm.GoAPI {
	x := cosmwasmAPIImpl{
		keeper: k,
		ctx:    &ctx,
	}
	return wasmvm.GoAPI{
		HumanAddress:      wasmkeeper.HumanAddress,
		CanonicalAddress:  wasmkeeper.CanonicalAddress,
		CallCallablePoint: x.callCallablePoint,
		ValidateInterface: x.validateInterface,
	}
}
