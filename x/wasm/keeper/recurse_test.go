package keeper

import (
	"encoding/json"
	"github.com/CosmWasm/wasmd/x/wasm/keeper/wasmtesting"
	"github.com/CosmWasm/wasmd/x/wasm/types"
	"testing"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type Recurse struct {
	Depth    uint32         `json:"depth"`
	Work     uint32         `json:"work"`
	Contract sdk.AccAddress `json:"contract"`
}

type recurseWrapper struct {
	Recurse Recurse `json:"recurse"`
}

func buildRecurseQuery(t *testing.T, msg Recurse) []byte {
	wrapper := recurseWrapper{Recurse: msg}
	bz, err := json.Marshal(wrapper)
	require.NoError(t, err)
	return bz
}

type recurseResponse struct {
	Hashed []byte `json:"hashed"`
}

// number os wasm queries called from a contract
var totalWasmQueryCounter int

func initRecurseContract(t *testing.T) (contract sdk.AccAddress, creator sdk.AccAddress, ctx sdk.Context, keeper *Keeper) {
	// we do one basic setup before all test cases (which are read-only and don't change state)
	var realWasmQuerier func(ctx sdk.Context, request *wasmvmtypes.WasmQuery) ([]byte, error)
	countingQuerier := &wasmtesting.MockQueryHandler{
		HandleQueryFn: func(ctx sdk.Context, request wasmvmtypes.QueryRequest, caller sdk.AccAddress) ([]byte, error) {
			totalWasmQueryCounter++
			return realWasmQuerier(ctx, request.Wasm)
		}}

	ctx, keepers := CreateTestInput(t, false, SupportedFeatures, WithQueryHandler(countingQuerier))
	keeper = keepers.WasmKeeper
	realWasmQuerier = WasmQuerier(keeper)

	exampleContract := InstantiateHackatomExampleContract(t, ctx, keepers)
	return exampleContract.Contract, exampleContract.CreatorAddr, ctx, keeper
}

func TestGasCostOnQuery(t *testing.T) {
	const (
		GasNoWork uint64 = 44_074
		// Note: about 100 SDK gas (10k wasmer gas) for each round of sha256
		GasWork50 uint64 = 49_744 // this is a little shy of 50k gas - to keep an eye on the limit

		GasReturnUnhashed uint64 = 287
		GasReturnHashed   uint64 = 262
	)

	cases := map[string]struct {
		gasLimit    uint64
		msg         Recurse
		expectedGas uint64
	}{
		"no recursion, no work": {
			gasLimit:    400_000,
			msg:         Recurse{},
			expectedGas: GasNoWork,
		},
		"no recursion, some work": {
			gasLimit: 400_000,
			msg: Recurse{
				Work: 50, // 50 rounds of sha256 inside the contract
			},
			expectedGas: GasWork50 + 1,
		},
		"recursion 1, no work": {
			gasLimit: 400_000,
			msg: Recurse{
				Depth: 1,
			},
			expectedGas: 2*GasNoWork + GasReturnUnhashed - 2,
		},
		"recursion 1, some work": {
			gasLimit: 400_000,
			msg: Recurse{
				Depth: 1,
				Work:  50,
			},
			expectedGas: 2*GasWork50 + GasReturnHashed,
		},
		"recursion 4, some work": {
			gasLimit: 400_000,
			msg: Recurse{
				Depth: 4,
				Work:  50,
			},
			// FIXME: why -6... confused a bit by calculations, seems like rounding issues
			expectedGas: 5*GasWork50 + 4*GasReturnHashed - 6,
		},
	}

	contractAddr, creator, ctx, keeper := initRecurseContract(t)

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// external limit has no effect (we get a panic if this is enforced)
			keeper.queryGasLimit = 1000

			// make sure we set a limit before calling
			ctx = ctx.WithGasMeter(sdk.NewGasMeter(tc.gasLimit))
			require.Equal(t, uint64(0), ctx.GasMeter().GasConsumed())

			// do the query
			recurse := tc.msg
			recurse.Contract = contractAddr
			msg := buildRecurseQuery(t, recurse)
			data, err := keeper.QuerySmart(ctx, contractAddr, msg)
			require.NoError(t, err)

			// check the gas is what we expected
			if types.EnableGasVerification {
				assert.Equal(t, tc.expectedGas, ctx.GasMeter().GasConsumed())
			}
			// assert result is 32 byte sha256 hash (if hashed), or contractAddr if not
			var resp recurseResponse
			err = json.Unmarshal(data, &resp)
			require.NoError(t, err)
			if recurse.Work == 0 {
				assert.Equal(t, len(resp.Hashed), len(creator.String()))
			} else {
				assert.Equal(t, len(resp.Hashed), 32)
			}
		})
	}
}

func TestGasOnExternalQuery(t *testing.T) {
	const (
		GasWork50 uint64 = InstanceCost + 8_464
	)

	cases := map[string]struct {
		gasLimit    uint64
		msg         Recurse
		expectPanic bool
	}{
		"no recursion, plenty gas": {
			gasLimit: 400_000,
			msg: Recurse{
				Work: 50, // 50 rounds of sha256 inside the contract
			},
		},
		"recursion 4, plenty gas": {
			// this uses 244708 gas
			gasLimit: 400_000,
			msg: Recurse{
				Depth: 4,
				Work:  50,
			},
		},
		"no recursion, external gas limit": {
			gasLimit: 5000, // this is not enough
			msg: Recurse{
				Work: 50,
			},
			expectPanic: true,
		},
		"recursion 4, external gas limit": {
			// this uses 244708 gas but give less
			gasLimit: 4 * GasWork50,
			msg: Recurse{
				Depth: 4,
				Work:  50,
			},
			expectPanic: true,
		},
	}

	contractAddr, _, ctx, keeper := initRecurseContract(t)

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// set the external gas limit (normally from config file)
			keeper.queryGasLimit = tc.gasLimit

			recurse := tc.msg
			recurse.Contract = contractAddr
			msg := buildRecurseQuery(t, recurse)

			// do the query
			path := []string{QueryGetContractState, contractAddr.String(), QueryMethodContractStateSmart}
			req := abci.RequestQuery{Data: msg}
			if tc.expectPanic {
				require.Panics(t, func() {
					// this should run out of gas
					_, err := NewLegacyQuerier(keeper)(ctx, path, req)
					t.Logf("%v", err)
				})
			} else {
				// otherwise, make sure we get a good success
				_, err := NewLegacyQuerier(keeper)(ctx, path, req)
				require.NoError(t, err)
			}
		})
	}
}

func TestLimitRecursiveQueryGas(t *testing.T) {
	// The point of this test from https://github.com/CosmWasm/cosmwasm/issues/456
	// Basically, if I burn 90% of gas in CPU loop, then query out (to my self)
	// the sub-query will have all the original gas (minus the 40k instance charge)
	// and can burn 90% and call a sub-contract again...
	// This attack would allow us to use far more than the provided gas before
	// eventually hitting an OutOfGas panic.

	const (
		// Note: about 100 SDK gas (10k wasmer gas) for each round of sha256
		GasWork2k uint64 = 272_797 // = InstanceCost + x // we have 6x gas used in cpu than in the instance
		// This is overhead for calling into a sub-contract
		GasReturnHashed uint64 = 265
	)

	cases := map[string]struct {
		gasLimit                  uint64
		msg                       Recurse
		expectQueriesFromContract int
		expectedGas               uint64
		expectOutOfGas            bool
	}{
		"no recursion, lots of work": {
			gasLimit: 4_000_000,
			msg: Recurse{
				Depth: 0,
				Work:  2000,
			},
			expectQueriesFromContract: 0,
			expectedGas:               GasWork2k,
		},
		"recursion 5, lots of work": {
			gasLimit: 4_000_000,
			msg: Recurse{
				Depth: 5,
				Work:  2000,
			},
			expectQueriesFromContract: 5,
			// FIXME: why +1 ... confused a bit by calculations, seems like rounding issues
			expectedGas: GasWork2k + 5*(GasWork2k+GasReturnHashed) + 1,
		},
		// this is where we expect an error...
		// it has enough gas to run 4 times and die on the 5th (4th time dispatching to sub-contract)
		// however, if we don't charge the cpu gas before sub-dispatching, we can recurse over 20 times
		// TODO: figure out how to asset how deep it went
		"deep recursion, should die on 5th level": {
			gasLimit: 1_200_000,
			msg: Recurse{
				Depth: 50,
				Work:  2000,
			},
			expectQueriesFromContract: 4,
			expectOutOfGas:            true,
		},
	}

	contractAddr, _, ctx, keeper := initRecurseContract(t)

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// reset the counter before test
			totalWasmQueryCounter = 0

			// make sure we set a limit before calling
			ctx = ctx.WithGasMeter(sdk.NewGasMeter(tc.gasLimit))
			require.Equal(t, uint64(0), ctx.GasMeter().GasConsumed())

			// prepare the query
			recurse := tc.msg
			recurse.Contract = contractAddr
			msg := buildRecurseQuery(t, recurse)

			// if we expect out of gas, make sure this panics
			if tc.expectOutOfGas {
				require.Panics(t, func() {
					_, err := keeper.QuerySmart(ctx, contractAddr, msg)
					t.Logf("Got error not panic: %#v", err)
				})
				assert.Equal(t, tc.expectQueriesFromContract, totalWasmQueryCounter)
				return
			}

			// otherwise, we expect a successful call
			_, err := keeper.QuerySmart(ctx, contractAddr, msg)
			require.NoError(t, err)
			if types.EnableGasVerification {
				assert.Equal(t, tc.expectedGas, ctx.GasMeter().GasConsumed())
			}
			assert.Equal(t, tc.expectQueriesFromContract, totalWasmQueryCounter)
		})
	}
}
