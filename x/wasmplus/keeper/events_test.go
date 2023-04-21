package keeper

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdk "github.com/line/lbm-sdk/types"
	wasmvmtypes "github.com/line/wasmvm/types"
)

func TestNewCustomCallablePointEvents(t *testing.T) {
	myContract := RandomAccountAddress(t)
	myCallstack := []sdk.AccAddress{RandomAccountAddress(t), RandomAccountAddress(t)}
	myCallstackBinary, err := json.Marshal(myCallstack)
	require.NoError(t, err)
	specs := map[string]struct {
		src     wasmvmtypes.Events
		exp     sdk.Events
		isError bool
	}{
		"all good": {
			src: wasmvmtypes.Events{{
				Type:       "foo",
				Attributes: []wasmvmtypes.EventAttribute{{Key: "myKey", Value: "myVal"}},
			}},
			exp: sdk.Events{sdk.NewEvent("wasm-callablepoint-foo",
				sdk.NewAttribute("_contract_address", myContract.String()),
				sdk.NewAttribute("_callstack", string(myCallstackBinary)),
				sdk.NewAttribute("myKey", "myVal"))},
		},
		"multiple attributes": {
			src: wasmvmtypes.Events{{
				Type: "foo",
				Attributes: []wasmvmtypes.EventAttribute{{Key: "myKey", Value: "myVal"},
					{Key: "myOtherKey", Value: "myOtherVal"}},
			}},
			exp: sdk.Events{sdk.NewEvent("wasm-callablepoint-foo",
				sdk.NewAttribute("_contract_address", myContract.String()),
				sdk.NewAttribute("_callstack", string(myCallstackBinary)),
				sdk.NewAttribute("myKey", "myVal"),
				sdk.NewAttribute("myOtherKey", "myOtherVal"))},
		},
		"multiple events": {
			src: wasmvmtypes.Events{{
				Type:       "foo",
				Attributes: []wasmvmtypes.EventAttribute{{Key: "myKey", Value: "myVal"}},
			}, {
				Type:       "bar",
				Attributes: []wasmvmtypes.EventAttribute{{Key: "otherKey", Value: "otherVal"}},
			}},
			exp: sdk.Events{
				sdk.NewEvent("wasm-callablepoint-foo",
					sdk.NewAttribute("_contract_address", myContract.String()),
					sdk.NewAttribute("_callstack", string(myCallstackBinary)),
					sdk.NewAttribute("myKey", "myVal")),
				sdk.NewEvent("wasm-callablepoint-bar",
					sdk.NewAttribute("_contract_address", myContract.String()),
					sdk.NewAttribute("_callstack", string(myCallstackBinary)),
					sdk.NewAttribute("otherKey", "otherVal")),
			},
		},
		"without attributes": {
			src: wasmvmtypes.Events{{
				Type: "foo",
			}},
			exp: sdk.Events{sdk.NewEvent("wasm-callablepoint-foo",
				sdk.NewAttribute("_contract_address", myContract.String()),
				sdk.NewAttribute("_callstack", string(myCallstackBinary))),
			},
		},
		"error on short event type": {
			src: wasmvmtypes.Events{{
				Type: "f",
			}},
			isError: true,
		},
		"error on _contract_address": {
			src: wasmvmtypes.Events{{
				Type:       "foo",
				Attributes: []wasmvmtypes.EventAttribute{{Key: "_contract_address", Value: RandomBech32AccountAddress(t)}},
			}},
			isError: true,
		},
		"error on reserved prefix": {
			src: wasmvmtypes.Events{{
				Type: "wasm",
				Attributes: []wasmvmtypes.EventAttribute{
					{Key: "_reserved", Value: "is skipped"},
					{Key: "normal", Value: "is used"}},
			}},
			isError: true,
		},
		"error on empty value": {
			src: wasmvmtypes.Events{{
				Type: "boom",
				Attributes: []wasmvmtypes.EventAttribute{
					{Key: "some", Value: "data"},
					{Key: "key", Value: ""},
				},
			}},
			isError: true,
		},
		"error on empty key": {
			src: wasmvmtypes.Events{{
				Type: "boom",
				Attributes: []wasmvmtypes.EventAttribute{
					{Key: "some", Value: "data"},
					{Key: "", Value: "value"},
				},
			}},
			isError: true,
		},
		"error on whitespace type": {
			src: wasmvmtypes.Events{{
				Type: "    f   ",
				Attributes: []wasmvmtypes.EventAttribute{
					{Key: "some", Value: "data"},
				},
			}},
			isError: true,
		},
		"error on only whitespace key": {
			src: wasmvmtypes.Events{{
				Type: "boom",
				Attributes: []wasmvmtypes.EventAttribute{
					{Key: "some", Value: "data"},
					{Key: "\n\n\n\n", Value: "value"},
				},
			}},
			isError: true,
		},
		"error on only whitespace value": {
			src: wasmvmtypes.Events{{
				Type: "boom",
				Attributes: []wasmvmtypes.EventAttribute{
					{Key: "some", Value: "data"},
					{Key: "myKey", Value: " \t\r\n"},
				},
			}},
			isError: true,
		},
		"strip out whitespace": {
			src: wasmvmtypes.Events{{
				Type:       "  food\n",
				Attributes: []wasmvmtypes.EventAttribute{{Key: "my Key", Value: "\tmyVal"}},
			}},
			exp: sdk.Events{sdk.NewEvent("wasm-callablepoint-food",
				sdk.NewAttribute("_contract_address", myContract.String()),
				sdk.NewAttribute("_callstack", string(myCallstackBinary)),

				sdk.NewAttribute("my Key", "myVal"))},
		},
		"empty event elements": {
			src:     make(wasmvmtypes.Events, 10),
			isError: true,
		},
		"nil": {
			exp: sdk.Events{},
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			gotEvent, err := newCustomCallablePointEvents(spec.src, myContract, myCallstackBinary)
			if spec.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, spec.exp, gotEvent)
			}
		})
	}
}

func TestNewCallablePointEvent(t *testing.T) {
	myContract := RandomAccountAddress(t)
	myCallstack := []sdk.AccAddress{RandomAccountAddress(t), RandomAccountAddress(t)}
	myCallstackBinary, err := json.Marshal(myCallstack)
	require.NoError(t, err)
	specs := map[string]struct {
		src     []wasmvmtypes.EventAttribute
		exp     sdk.Events
		isError bool
	}{
		"all good": {
			src: []wasmvmtypes.EventAttribute{{Key: "myKey", Value: "myVal"}},
			exp: sdk.Events{sdk.NewEvent("wasm-callablepoint",
				sdk.NewAttribute("_contract_address", myContract.String()),
				sdk.NewAttribute("_callstack", string(myCallstackBinary)),
				sdk.NewAttribute("myKey", "myVal"))},
		},
		"multiple attributes": {
			src: []wasmvmtypes.EventAttribute{{Key: "myKey", Value: "myVal"},
				{Key: "myOtherKey", Value: "myOtherVal"}},
			exp: sdk.Events{sdk.NewEvent("wasm-callablepoint",
				sdk.NewAttribute("_contract_address", myContract.String()),
				sdk.NewAttribute("_callstack", string(myCallstackBinary)),
				sdk.NewAttribute("myKey", "myVal"),
				sdk.NewAttribute("myOtherKey", "myOtherVal"))},
		},
		"without attributes": {
			exp: sdk.Events{sdk.NewEvent("wasm-callablepoint",
				sdk.NewAttribute("_contract_address", myContract.String()), sdk.NewAttribute("_callstack", string(myCallstackBinary))),
			},
		},
		"error on _contract_address": {
			src:     []wasmvmtypes.EventAttribute{{Key: "_contract_address", Value: RandomBech32AccountAddress(t)}},
			isError: true,
		},
		"error on whitespace key": {
			src:     []wasmvmtypes.EventAttribute{{Key: "  ", Value: "value"}},
			isError: true,
		},
		"error on whitespace value": {
			src:     []wasmvmtypes.EventAttribute{{Key: "key", Value: "\n\n\n"}},
			isError: true,
		},
		"strip whitespace": {
			src: []wasmvmtypes.EventAttribute{{Key: "   my-real-key    ", Value: "\n\n\nsome-val\t\t\t"}},
			exp: sdk.Events{sdk.NewEvent("wasm-callablepoint",
				sdk.NewAttribute("_contract_address", myContract.String()),
				sdk.NewAttribute("_callstack", string(myCallstackBinary)),
				sdk.NewAttribute("my-real-key", "some-val"))},
		},
		"empty elements": {
			src:     make([]wasmvmtypes.EventAttribute, 10),
			isError: true,
		},
		"nil": {
			exp: sdk.Events{sdk.NewEvent("wasm-callablepoint",
				sdk.NewAttribute("_contract_address", myContract.String()),
				sdk.NewAttribute("_callstack", string(myCallstackBinary)),
			)},
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			gotEvent, err := newCallablePointEvent(spec.src, myContract, myCallstackBinary)
			if spec.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, spec.exp, gotEvent)
			}
		})
	}
}
