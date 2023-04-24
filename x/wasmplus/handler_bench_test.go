package wasmplus

import (
	"testing"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/wasmd/x/wasm"
	wasmkeeper "github.com/Finschia/wasmd/x/wasm/keeper"
	wasmtypes "github.com/Finschia/wasmd/x/wasm/types"
	"github.com/Finschia/wasmd/x/wasmplus/keeper"
	"github.com/stretchr/testify/require"
)

func setupHandlerBenches(b *testing.B) (sdk.Context, keeper.TestKeepers, wasmtypes.MsgStoreCode) {
	ctx, keepers := keeper.CreateTestInput(b, false, "iterator,staking,stargate,cosmwasm_1_1")
	creator := keepers.Faucet.NewFundedRandomAccount(ctx, sdk.NewInt64Coin("denom", 100000))

	msg := wasmtypes.MsgStoreCode{
		Sender:       creator.String(),
		WASMByteCode: testContract,
	}

	return ctx, keepers, msg
}

func BenchmarkWasmHandler(b *testing.B) {
	ctx, keepers, msg := setupHandlerBenches(b)

	// prepare test wasm handler
	handler := wasm.NewHandler(keeper.NewPermissionedKeeper(*wasmkeeper.NewDefaultPermissionKeeper(keepers.WasmKeeper), keepers.WasmKeeper))

	for i := 0; i < b.N; i++ {
		_, err := handler(ctx, sdk.Msg(&msg))
		require.NoError(b, err)
	}
}

func BenchmarkWasmPlusHandler(b *testing.B) {
	ctx, keepers, msg := setupHandlerBenches(b)

	// prepare test wasm handler
	handler := NewHandler(keeper.NewPermissionedKeeper(*wasmkeeper.NewDefaultPermissionKeeper(keepers.WasmKeeper), keepers.WasmKeeper))

	for i := 0; i < b.N; i++ {
		_, err := handler(ctx, sdk.Msg(&msg))
		require.NoError(b, err)
	}
}
