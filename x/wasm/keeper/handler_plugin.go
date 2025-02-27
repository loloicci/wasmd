package keeper

import (
	"errors"
	"fmt"
	"github.com/CosmWasm/wasmd/x/wasm/types"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	host "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
)

// msgEncoder is an extension point to customize encodings
type msgEncoder interface {
	// Encode converts wasmvm message to n cosmos message types
	Encode(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) ([]sdk.Msg, error)
}

// SDKMessageHandler can handles messages that can be encoded into sdk.Message types and routed.
type SDKMessageHandler struct {
	router   sdk.Router
	encoders msgEncoder
}

func NewDefaultMessageHandler(router sdk.Router, channelKeeper types.ChannelKeeper, capabilityKeeper types.CapabilityKeeper, unpacker codectypes.AnyUnpacker, portSource types.ICS20TransferPortSource, customEncoders ...*MessageEncoders) messenger {
	encoders := DefaultEncoders(unpacker, portSource)
	for _, e := range customEncoders {
		encoders = encoders.Merge(e)
	}
	return NewMessageHandlerChain(
		NewSDKMessageHandler(router, encoders),
		NewIBCRawPacketHandler(channelKeeper, capabilityKeeper),
	)
}

func NewSDKMessageHandler(router sdk.Router, encoders msgEncoder) SDKMessageHandler {
	return SDKMessageHandler{
		router:   router,
		encoders: encoders,
	}
}

func (h SDKMessageHandler) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) (events []sdk.Event, data [][]byte, err error) {
	sdkMsgs, err := h.encoders.Encode(ctx, contractAddr, contractIBCPortID, msg)
	if err != nil {
		return nil, nil, err
	}
	for _, sdkMsg := range sdkMsgs {
		res, err := h.handleSdkMessage(ctx, contractAddr, sdkMsg)
		if err != nil {
			return nil, nil, err
		}
		// append data
		data = append(data, res.Data)
		// append events
		sdkEvents := make([]sdk.Event, len(res.Events))
		for i := range res.Events {
			sdkEvents[i] = sdk.Event(res.Events[i])
		}
		events = append(events, sdkEvents...)
	}
	return
}

func (h SDKMessageHandler) handleSdkMessage(ctx sdk.Context, contractAddr sdk.Address, msg sdk.Msg) (*sdk.Result, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	// make sure this account can send it
	for _, acct := range msg.GetSigners() {
		if !acct.Equals(contractAddr) {
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "contract doesn't have permission")
		}
	}

	// find the handler and execute it
	handler := h.router.Route(ctx, msg.Route())
	if handler == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, msg.Route())
	}
	res, err := handler(ctx, msg)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// MessageHandlerChain defines a chain of handlers that are called one by one until it can be handled.
type MessageHandlerChain struct {
	handlers []messenger
}

func NewMessageHandlerChain(first messenger, others ...messenger) *MessageHandlerChain {
	r := &MessageHandlerChain{handlers: append([]messenger{first}, others...)}
	for i := range r.handlers {
		if r.handlers[i] == nil {
			panic(fmt.Sprintf("handler must not be nil at position : %d", i))
		}
	}
	return r
}

// DispatchMsg dispatch message to handlers.
func (m MessageHandlerChain) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) ([]sdk.Event, [][]byte, error) {
	for _, h := range m.handlers {
		events, data, err := h.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
		switch {
		case err == nil:
			return events, data, err
		case errors.Is(err, types.ErrUnknownMsg):
			continue
		default:
			return events, data, err
		}
	}
	return nil, nil, sdkerrors.Wrap(types.ErrUnknownMsg, "no handler found")
}

// IBCRawPacketHandler handels IBC.SendPacket messages which are published to an IBC channel.
type IBCRawPacketHandler struct {
	channelKeeper    types.ChannelKeeper
	capabilityKeeper types.CapabilityKeeper
}

func NewIBCRawPacketHandler(chk types.ChannelKeeper, cak types.CapabilityKeeper) *IBCRawPacketHandler {
	return &IBCRawPacketHandler{channelKeeper: chk, capabilityKeeper: cak}
}

// DispatchMsg publishes a raw IBC packet onto the channel.
func (h IBCRawPacketHandler) DispatchMsg(ctx sdk.Context, _ sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) (events []sdk.Event, data [][]byte, err error) {
	if msg.IBC == nil || msg.IBC.SendPacket == nil {
		return nil, nil, types.ErrUnknownMsg
	}
	if contractIBCPortID == "" {
		return nil, nil, sdkerrors.Wrapf(types.ErrUnsupportedForContract, "ibc not supported")
	}
	contractIBCChannelID := msg.IBC.SendPacket.ChannelID
	if contractIBCChannelID == "" {
		return nil, nil, sdkerrors.Wrapf(types.ErrEmpty, "ibc channel")
	}

	sequence, found := h.channelKeeper.GetNextSequenceSend(ctx, contractIBCPortID, contractIBCChannelID)
	if !found {
		return nil, nil, sdkerrors.Wrapf(channeltypes.ErrSequenceSendNotFound,
			"source port: %s, source channel: %s", contractIBCPortID, contractIBCChannelID,
		)
	}

	channelInfo, ok := h.channelKeeper.GetChannel(ctx, contractIBCPortID, contractIBCChannelID)
	if !ok {
		return nil, nil, sdkerrors.Wrap(channeltypes.ErrInvalidChannel, "not found")
	}
	channelCap, ok := h.capabilityKeeper.GetCapability(ctx, host.ChannelCapabilityPath(contractIBCPortID, contractIBCChannelID))
	if !ok {
		return nil, nil, sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}
	packet := channeltypes.NewPacket(
		msg.IBC.SendPacket.Data,
		sequence,
		contractIBCPortID,
		contractIBCChannelID,
		channelInfo.Counterparty.PortId,
		channelInfo.Counterparty.ChannelId,
		convertWasmIBCTimeoutHeightToCosmosHeight(msg.IBC.SendPacket.TimeoutBlock),
		convertWasmIBCTimeoutTimestampToCosmosTimestamp(msg.IBC.SendPacket.TimeoutTimestamp),
	)
	return nil, nil, h.channelKeeper.SendPacket(ctx, channelCap, packet)
}
