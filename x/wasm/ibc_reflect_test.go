package wasm_test

import (
	"github.com/CosmWasm/wasmd/x/wasm/ibctesting"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	ibcexported "github.com/cosmos/cosmos-sdk/x/ibc/core/exported"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIBCReflectContract(t *testing.T) {
	var (
		coordinator = ibctesting.NewCoordinator(t, 2, nil, nil)
		chainA      = coordinator.GetChain(ibctesting.GetChainID(0))
		chainB      = coordinator.GetChain(ibctesting.GetChainID(1))
	)
	coordinator.CommitBlock(chainA, chainB)

	initMsg := []byte(`{}`)
	codeID := chainA.StoreCodeFile("./keeper/testdata/ibc_reflect_send.wasm").CodeID
	sendContractAddr := chainA.InstantiateContract(codeID, initMsg)

	reflectID := chainB.StoreCodeFile("./keeper/testdata/reflect.wasm").CodeID
	initMsg = wasmkeeper.IBCReflectInitMsg{
		ReflectCodeID: reflectID,
	}.GetBytes(t)
	codeID = chainB.StoreCodeFile("./keeper/testdata/ibc_reflect.wasm").CodeID

	reflectContractAddr := chainB.InstantiateContract(codeID, initMsg)
	var (
		sourcePortID      = chainA.ContractInfo(sendContractAddr).IBCPortID
		counterpartPortID = chainB.ContractInfo(reflectContractAddr).IBCPortID
	)
	clientA, clientB, connA, connB := coordinator.SetupClientConnections(chainA, chainB, ibcexported.Tendermint)
	connA.NextChannelVersion = "ibc-reflect-v1"
	connB.NextChannelVersion = "ibc-reflect-v1"
	// flip instantiation so that we do not run into https://github.com/cosmos/cosmos-sdk/issues/8334
	channelA, channelB := coordinator.CreateChannel(chainA, chainB, connA, connB, sourcePortID, counterpartPortID, channeltypes.ORDERED)

	// TODO: query both contracts directly to ensure they have registered the proper connection
	// (and the chainB has created a reflect contract)

	// there should be one packet to relay back and forth (whoami)
	// TODO: how do I find the packet that was previously sent by the smart contract?
	// Coordinator.RecvPacket requires channeltypes.Packet as input?
	// Given the source (portID, channelID), we should be able to count how many packets are pending, query the data
	// and submit them to the other side (same with acks). This is what the real relayer does. I guess the test framework doesn't?

	// Update: I dug through the code, expecially channel.Keeper.SendPacket, and it only writes a commitment
	// only writes I see: https://github.com/cosmos/cosmos-sdk/blob/31fdee0228bd6f3e787489c8e4434aabc8facb7d/x/ibc/core/04-channel/keeper/packet.go#L115-L116
	// commitment is hashed packet: https://github.com/cosmos/cosmos-sdk/blob/31fdee0228bd6f3e787489c8e4434aabc8facb7d/x/ibc/core/04-channel/types/packet.go#L14-L34
	// how is the relayer supposed to get the original packet data??
	// eg. ibctransfer doesn't store the packet either: https://github.com/cosmos/cosmos-sdk/blob/master/x/ibc/applications/transfer/keeper/relay.go#L145-L162
	// ... or I guess the original packet data is only available in the event logs????
	// https://github.com/cosmos/cosmos-sdk/blob/31fdee0228bd6f3e787489c8e4434aabc8facb7d/x/ibc/core/04-channel/keeper/packet.go#L121-L132

	// ensure the expected packet was prepared, and relay it
	require.Equal(t, 1, len(chainA.PendingSendPackets))
	require.Equal(t, 0, len(chainB.PendingSendPackets))
	err := coordinator.RelayAndAckPendingPackets(chainA, chainB, clientA, clientB)
	require.NoError(t, err)
	require.Equal(t, 0, len(chainA.PendingSendPackets))
	require.Equal(t, 0, len(chainB.PendingSendPackets))

	// let's query the source contract and make sure it registered an address
	query := ReflectSendQueryMsg{Account: &AccountQuery{ChannelID: channelA.ID}}
	var account AccountResponse
	err = chainA.SmartQuery(sendContractAddr.String(), query, &account)
	require.NoError(t, err)
	require.NotEmpty(t, account.RemoteAddr)
	require.Empty(t, account.RemoteBalance)

	// close channel
	coordinator.CloseChannel(chainA, chainB, channelA, channelB)

	// let's query the source contract and make sure it registered an address
	account = AccountResponse{}
	err = chainA.SmartQuery(sendContractAddr.String(), query, &account)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	_ = clientB
}

type ReflectSendQueryMsg struct {
	Admin        *struct{}     `json:"admin,omitempty"`
	ListAccounts *struct{}     `json:"list_accounts,omitempty"`
	Account      *AccountQuery `json:"account,omitempty"`
}

type AccountQuery struct {
	ChannelID string `json:"channel_id"`
}

type AccountResponse struct {
	LastUpdateTime uint64            `json:"last_update_time"`
	RemoteAddr     string            `json:"remote_addr"`
	RemoteBalance  wasmvmtypes.Coins `json:"remote_balance"`
}
