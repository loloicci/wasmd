package keeper

import (
	"fmt"
	"strings"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	wasmvmtypes "github.com/line/wasmvm/types"

	wasmkeeper "github.com/line/wasmd/x/wasm/keeper"
	"github.com/line/wasmd/x/wasm/types"
	wasmplustypes "github.com/line/wasmd/x/wasmplus/types"
)

func newCallablePointEvent(customAttributes []wasmvmtypes.EventAttribute, contractAddr sdk.AccAddress, callstack []byte) (sdk.Events, error) {
	attrs, err := callablePointSDKEventAttributes(customAttributes, contractAddr, callstack)

	if err != nil {
		return nil, err
	}

	return sdk.Events{sdk.NewEvent(wasmplustypes.CallablePointEventType, attrs...)}, nil
}

func newCustomCallablePointEvents(evts wasmvmtypes.Events, contractAddr sdk.AccAddress, callstack []byte) (sdk.Events, error) {
	events := make(sdk.Events, 0, len(evts))
	for _, e := range evts {
		typ := strings.TrimSpace(e.Type)
		if len(typ) <= wasmkeeper.EventTypeMinLength {
			return nil, sdkerrors.Wrap(types.ErrInvalidEvent, fmt.Sprintf("Event type too short: '%s'", typ))
		}
		attributes, err := callablePointSDKEventAttributes(e.Attributes, contractAddr, callstack)
		if err != nil {
			return nil, err
		}
		events = append(events, sdk.NewEvent(fmt.Sprintf("%s%s", wasmplustypes.CustomCallablePointEventPrefix, typ), attributes...))
	}
	return events, nil
}

func callablePointSDKEventAttributes(customAttributes []wasmvmtypes.EventAttribute, contractAddr sdk.AccAddress, callstack []byte) ([]sdk.Attribute, error) {
	attrs, err := wasmkeeper.ContractSDKEventAttributes(customAttributes, contractAddr)
	if err != nil {
		return nil, err
	}
	// attrs[0] is addr
	attrs = append([]sdk.Attribute{attrs[0], sdk.NewAttribute(wasmplustypes.AttributeKeyCallstack, string(callstack))}, attrs[1:]...)
	return attrs, nil
}
