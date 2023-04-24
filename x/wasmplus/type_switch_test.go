package wasmplus

import (
	"testing"

	sdk "github.com/Finschia/finschia-sdk/types"
	wasmtypes "github.com/Finschia/wasmd/x/wasm/types"
	"github.com/Finschia/wasmd/x/wasmplus/types"
)

func BenchmarkTypeSwitch(b *testing.B) {
	msg := &wasmtypes.MsgStoreCode{
		Sender:       "foo",
		WASMByteCode: []byte("invalid WASM contract"),
	}
	for i := 0; i < b.N; i++ {
		typeSwitch(sdk.Msg(msg))
	}
}

func BenchmarkNoTypeSwitch(b *testing.B) {
	msg := &wasmtypes.MsgStoreCode{
		Sender:       "foo",
		WASMByteCode: []byte("invalid WASM contract"),
	}
	for i := 0; i < b.N; i++ {
		doNothing(sdk.Msg(msg))
	}
}

func typeSwitch(msg sdk.Msg) *types.MsgStoreCodeAndInstantiateContract {
	switch msg := msg.(type) {
	case *types.MsgStoreCodeAndInstantiateContract:
		return msg
	default:
		return nil
	}
}

func doNothing(msg sdk.Msg) *types.MsgStoreCodeAndInstantiateContract {
	return nil
}
