package wasmtesting

import (
	"bytes"
	"crypto/sha256"
	"github.com/CosmWasm/wasmd/x/wasm/types"
	wasmvm "github.com/CosmWasm/wasmvm"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/libs/rand"
)

var _ types.WasmerEngine = &MockWasmer{}

// MockWasmer implements types.WasmerEngine for testing purpose. One or multiple messages can be stubbed.
// Without a stub function a panic is thrown.
type MockWasmer struct {
	CreateFn            func(codeID wasmvm.WasmCode) (wasmvm.Checksum, error)
	AnalyzeCodeFn       func(codeID wasmvm.Checksum) (*wasmvmtypes.AnalysisReport, error)
	InstantiateFn       func(codeID wasmvm.Checksum, env wasmvmtypes.Env, info wasmvmtypes.MessageInfo, initMsg []byte, store wasmvm.KVStore, goapi wasmvm.GoAPI, querier wasmvm.Querier, gasMeter wasmvm.GasMeter, gasLimit uint64) (*wasmvmtypes.Response, uint64, error)
	ExecuteFn           func(codeID wasmvm.Checksum, env wasmvmtypes.Env, info wasmvmtypes.MessageInfo, executeMsg []byte, store wasmvm.KVStore, goapi wasmvm.GoAPI, querier wasmvm.Querier, gasMeter wasmvm.GasMeter, gasLimit uint64) (*wasmvmtypes.Response, uint64, error)
	QueryFn             func(codeID wasmvm.Checksum, env wasmvmtypes.Env, queryMsg []byte, store wasmvm.KVStore, goapi wasmvm.GoAPI, querier wasmvm.Querier, gasMeter wasmvm.GasMeter, gasLimit uint64) ([]byte, uint64, error)
	MigrateFn           func(codeID wasmvm.Checksum, env wasmvmtypes.Env, migrateMsg []byte, store wasmvm.KVStore, goapi wasmvm.GoAPI, querier wasmvm.Querier, gasMeter wasmvm.GasMeter, gasLimit uint64) (*wasmvmtypes.Response, uint64, error)
	SudoFn              func(codeID wasmvm.Checksum, env wasmvmtypes.Env, sudoMsg []byte, store wasmvm.KVStore, goapi wasmvm.GoAPI, querier wasmvm.Querier, gasMeter wasmvm.GasMeter, gasLimit uint64) (*wasmvmtypes.Response, uint64, error)
	ReplyFn             func(codeID wasmvm.Checksum, env wasmvmtypes.Env, reply wasmvmtypes.Reply, store wasmvm.KVStore, goapi wasmvm.GoAPI, querier wasmvm.Querier, gasMeter wasmvm.GasMeter, gasLimit uint64) (*wasmvmtypes.Response, uint64, error)
	GetCodeFn           func(codeID wasmvm.Checksum) (wasmvm.WasmCode, error)
	CleanupFn           func()
	IBCChannelOpenFn    func(codeID wasmvm.Checksum, env wasmvmtypes.Env, channel wasmvmtypes.IBCChannel, store wasmvm.KVStore, goapi wasmvm.GoAPI, querier wasmvm.Querier, gasMeter wasmvm.GasMeter, gasLimit uint64) (uint64, error)
	IBCChannelConnectFn func(codeID wasmvm.Checksum, env wasmvmtypes.Env, channel wasmvmtypes.IBCChannel, store wasmvm.KVStore, goapi wasmvm.GoAPI, querier wasmvm.Querier, gasMeter wasmvm.GasMeter, gasLimit uint64) (*wasmvmtypes.IBCBasicResponse, uint64, error)
	IBCChannelCloseFn   func(codeID wasmvm.Checksum, env wasmvmtypes.Env, channel wasmvmtypes.IBCChannel, store wasmvm.KVStore, goapi wasmvm.GoAPI, querier wasmvm.Querier, gasMeter wasmvm.GasMeter, gasLimit uint64) (*wasmvmtypes.IBCBasicResponse, uint64, error)
	IBCPacketReceiveFn  func(codeID wasmvm.Checksum, env wasmvmtypes.Env, packet wasmvmtypes.IBCPacket, store wasmvm.KVStore, goapi wasmvm.GoAPI, querier wasmvm.Querier, gasMeter wasmvm.GasMeter, gasLimit uint64) (*wasmvmtypes.IBCReceiveResponse, uint64, error)
	IBCPacketAckFn      func(codeID wasmvm.Checksum, env wasmvmtypes.Env, ack wasmvmtypes.IBCAcknowledgement, store wasmvm.KVStore, goapi wasmvm.GoAPI, querier wasmvm.Querier, gasMeter wasmvm.GasMeter, gasLimit uint64) (*wasmvmtypes.IBCBasicResponse, uint64, error)
	IBCPacketTimeoutFn  func(codeID wasmvm.Checksum, env wasmvmtypes.Env, packet wasmvmtypes.IBCPacket, store wasmvm.KVStore, goapi wasmvm.GoAPI, querier wasmvm.Querier, gasMeter wasmvm.GasMeter, gasLimit uint64) (*wasmvmtypes.IBCBasicResponse, uint64, error)
	PinFn               func(checksum wasmvm.Checksum) error
	UnpinFn             func(checksum wasmvm.Checksum) error
}

func (m *MockWasmer) IBCChannelOpen(codeID wasmvm.Checksum, env wasmvmtypes.Env, channel wasmvmtypes.IBCChannel, store wasmvm.KVStore, goapi wasmvm.GoAPI, querier wasmvm.Querier, gasMeter wasmvm.GasMeter, gasLimit uint64) (uint64, error) {
	if m.IBCChannelOpenFn == nil {
		panic("not supposed to be called!")
	}
	return m.IBCChannelOpenFn(codeID, env, channel, store, goapi, querier, gasMeter, gasLimit)
}

func (m *MockWasmer) IBCChannelConnect(codeID wasmvm.Checksum, env wasmvmtypes.Env, channel wasmvmtypes.IBCChannel, store wasmvm.KVStore, goapi wasmvm.GoAPI, querier wasmvm.Querier, gasMeter wasmvm.GasMeter, gasLimit uint64) (*wasmvmtypes.IBCBasicResponse, uint64, error) {
	if m.IBCChannelConnectFn == nil {
		panic("not supposed to be called!")
	}
	return m.IBCChannelConnectFn(codeID, env, channel, store, goapi, querier, gasMeter, gasLimit)
}

func (m *MockWasmer) IBCChannelClose(codeID wasmvm.Checksum, env wasmvmtypes.Env, channel wasmvmtypes.IBCChannel, store wasmvm.KVStore, goapi wasmvm.GoAPI, querier wasmvm.Querier, gasMeter wasmvm.GasMeter, gasLimit uint64) (*wasmvmtypes.IBCBasicResponse, uint64, error) {
	if m.IBCChannelCloseFn == nil {
		panic("not supposed to be called!")
	}
	return m.IBCChannelCloseFn(codeID, env, channel, store, goapi, querier, gasMeter, gasLimit)
}

func (m *MockWasmer) IBCPacketReceive(codeID wasmvm.Checksum, env wasmvmtypes.Env, packet wasmvmtypes.IBCPacket, store wasmvm.KVStore, goapi wasmvm.GoAPI, querier wasmvm.Querier, gasMeter wasmvm.GasMeter, gasLimit uint64) (*wasmvmtypes.IBCReceiveResponse, uint64, error) {
	if m.IBCPacketReceiveFn == nil {
		panic("not supposed to be called!")
	}
	return m.IBCPacketReceiveFn(codeID, env, packet, store, goapi, querier, gasMeter, gasLimit)
}

func (m *MockWasmer) IBCPacketAck(codeID wasmvm.Checksum, env wasmvmtypes.Env, ack wasmvmtypes.IBCAcknowledgement, store wasmvm.KVStore, goapi wasmvm.GoAPI, querier wasmvm.Querier, gasMeter wasmvm.GasMeter, gasLimit uint64) (*wasmvmtypes.IBCBasicResponse, uint64, error) {
	if m.IBCPacketAckFn == nil {
		panic("not supposed to be called!")
	}
	return m.IBCPacketAckFn(codeID, env, ack, store, goapi, querier, gasMeter, gasLimit)
}

func (m *MockWasmer) IBCPacketTimeout(codeID wasmvm.Checksum, env wasmvmtypes.Env, packet wasmvmtypes.IBCPacket, store wasmvm.KVStore, goapi wasmvm.GoAPI, querier wasmvm.Querier, gasMeter wasmvm.GasMeter, gasLimit uint64) (*wasmvmtypes.IBCBasicResponse, uint64, error) {
	if m.IBCPacketTimeoutFn == nil {
		panic("not supposed to be called!")
	}
	return m.IBCPacketTimeoutFn(codeID, env, packet, store, goapi, querier, gasMeter, gasLimit)
}

func (m *MockWasmer) Create(codeID wasmvm.WasmCode) (wasmvm.Checksum, error) {
	if m.CreateFn == nil {
		panic("not supposed to be called!")
	}
	return m.CreateFn(codeID)
}

func (m *MockWasmer) AnalyzeCode(codeID wasmvm.Checksum) (*wasmvmtypes.AnalysisReport, error) {
	if m.AnalyzeCodeFn == nil {
		panic("not supposed to be called!")
	}
	return m.AnalyzeCodeFn(codeID)
}

func (m *MockWasmer) Instantiate(codeID wasmvm.Checksum, env wasmvmtypes.Env, info wasmvmtypes.MessageInfo, initMsg []byte, store wasmvm.KVStore, goapi wasmvm.GoAPI, querier wasmvm.Querier, gasMeter wasmvm.GasMeter, gasLimit uint64) (*wasmvmtypes.Response, uint64, error) {
	if m.InstantiateFn == nil {
		panic("not supposed to be called!")
	}
	return m.InstantiateFn(codeID, env, info, initMsg, store, goapi, querier, gasMeter, gasLimit)
}

func (m *MockWasmer) Execute(codeID wasmvm.Checksum, env wasmvmtypes.Env, info wasmvmtypes.MessageInfo, executeMsg []byte, store wasmvm.KVStore, goapi wasmvm.GoAPI, querier wasmvm.Querier, gasMeter wasmvm.GasMeter, gasLimit uint64) (*wasmvmtypes.Response, uint64, error) {
	if m.ExecuteFn == nil {
		panic("not supposed to be called!")
	}
	return m.ExecuteFn(codeID, env, info, executeMsg, store, goapi, querier, gasMeter, gasLimit)
}

func (m *MockWasmer) Query(codeID wasmvm.Checksum, env wasmvmtypes.Env, queryMsg []byte, store wasmvm.KVStore, goapi wasmvm.GoAPI, querier wasmvm.Querier, gasMeter wasmvm.GasMeter, gasLimit uint64) ([]byte, uint64, error) {
	if m.QueryFn == nil {
		panic("not supposed to be called!")
	}
	return m.QueryFn(codeID, env, queryMsg, store, goapi, querier, gasMeter, gasLimit)
}

func (m *MockWasmer) Migrate(codeID wasmvm.Checksum, env wasmvmtypes.Env, migrateMsg []byte, store wasmvm.KVStore, goapi wasmvm.GoAPI, querier wasmvm.Querier, gasMeter wasmvm.GasMeter, gasLimit uint64) (*wasmvmtypes.Response, uint64, error) {
	if m.MigrateFn == nil {
		panic("not supposed to be called!")
	}
	return m.MigrateFn(codeID, env, migrateMsg, store, goapi, querier, gasMeter, gasLimit)
}

func (m *MockWasmer) Sudo(codeID wasmvm.Checksum, env wasmvmtypes.Env, sudoMsg []byte, store wasmvm.KVStore, goapi wasmvm.GoAPI, querier wasmvm.Querier, gasMeter wasmvm.GasMeter, gasLimit uint64) (*wasmvmtypes.Response, uint64, error) {
	if m.SudoFn == nil {
		panic("not supposed to be called!")
	}
	return m.SudoFn(codeID, env, sudoMsg, store, goapi, querier, gasMeter, gasLimit)

}

func (m *MockWasmer) Reply(codeID wasmvm.Checksum, env wasmvmtypes.Env, reply wasmvmtypes.Reply, store wasmvm.KVStore, goapi wasmvm.GoAPI, querier wasmvm.Querier, gasMeter wasmvm.GasMeter, gasLimit uint64) (*wasmvmtypes.Response, uint64, error) {
	if m.ReplyFn == nil {
		panic("not supposed to be called!")
	}
	return m.ReplyFn(codeID, env, reply, store, goapi, querier, gasMeter, gasLimit)
}

func (m *MockWasmer) GetCode(codeID wasmvm.Checksum) (wasmvm.WasmCode, error) {
	if m.GetCodeFn == nil {
		panic("not supposed to be called!")
	}
	return m.GetCodeFn(codeID)
}

func (m *MockWasmer) Cleanup() {
	if m.CleanupFn == nil {
		panic("not supposed to be called!")
	}
	m.CleanupFn()
}

func (m *MockWasmer) Pin(checksum wasmvm.Checksum) error {
	if m.PinFn == nil {
		panic("not supposed to be called!")
	}
	return m.PinFn(checksum)
}

func (m *MockWasmer) Unpin(checksum wasmvm.Checksum) error {
	if m.UnpinFn == nil {
		panic("not supposed to be called!")
	}
	return m.UnpinFn(checksum)
}

var AlwaysPanicMockWasmer = &MockWasmer{}

// selfCallingInstMockWasmer prepares a Wasmer mock that calls itself on instantiation.
func SelfCallingInstMockWasmer(executeCalled *bool) *MockWasmer {
	return &MockWasmer{

		CreateFn: func(code wasmvm.WasmCode) (wasmvm.Checksum, error) {
			anyCodeID := bytes.Repeat([]byte{0x1}, 32)
			return anyCodeID, nil
		},
		InstantiateFn: func(codeID wasmvm.Checksum, env wasmvmtypes.Env, info wasmvmtypes.MessageInfo, initMsg []byte, store wasmvm.KVStore, goapi wasmvm.GoAPI, querier wasmvm.Querier, gasMeter wasmvm.GasMeter, gasLimit uint64) (*wasmvmtypes.Response, uint64, error) {
			return &wasmvmtypes.Response{
				Messages: []wasmvmtypes.CosmosMsg{
					{Wasm: &wasmvmtypes.WasmMsg{Execute: &wasmvmtypes.ExecuteMsg{ContractAddr: env.Contract.Address, Msg: []byte(`{}`)}}},
				},
			}, 1, nil
		},
		AnalyzeCodeFn: WithoutIBCAnalyzeFn,
		ExecuteFn: func(codeID wasmvm.Checksum, env wasmvmtypes.Env, info wasmvmtypes.MessageInfo, executeMsg []byte, store wasmvm.KVStore, goapi wasmvm.GoAPI, querier wasmvm.Querier, gasMeter wasmvm.GasMeter, gasLimit uint64) (*wasmvmtypes.Response, uint64, error) {
			*executeCalled = true
			return &wasmvmtypes.Response{}, 1, nil
		},
	}
}

// IBCContractCallbacks defines the methods from wasmvm to interact with the wasm contract.
// A mock contract would implement the interface to fully simulate a wasm contract's behaviour.
type IBCContractCallbacks interface {
	IBCChannelOpen(
		codeID wasmvm.Checksum,
		env wasmvmtypes.Env,
		channel wasmvmtypes.IBCChannel,
		store wasmvm.KVStore,
		goapi wasmvm.GoAPI,
		querier wasmvm.Querier,
		gasMeter wasmvm.GasMeter,
		gasLimit uint64,
	) (uint64, error)

	IBCChannelConnect(
		codeID wasmvm.Checksum,
		env wasmvmtypes.Env,
		channel wasmvmtypes.IBCChannel,
		store wasmvm.KVStore,
		goapi wasmvm.GoAPI,
		querier wasmvm.Querier,
		gasMeter wasmvm.GasMeter,
		gasLimit uint64,
	) (*wasmvmtypes.IBCBasicResponse, uint64, error)

	IBCChannelClose(
		codeID wasmvm.Checksum,
		env wasmvmtypes.Env,
		channel wasmvmtypes.IBCChannel,
		store wasmvm.KVStore,
		goapi wasmvm.GoAPI,
		querier wasmvm.Querier,
		gasMeter wasmvm.GasMeter,
		gasLimit uint64,
	) (*wasmvmtypes.IBCBasicResponse, uint64, error)

	IBCPacketReceive(
		codeID wasmvm.Checksum,
		env wasmvmtypes.Env,
		packet wasmvmtypes.IBCPacket,
		store wasmvm.KVStore,
		goapi wasmvm.GoAPI,
		querier wasmvm.Querier,
		gasMeter wasmvm.GasMeter,
		gasLimit uint64,
	) (*wasmvmtypes.IBCReceiveResponse, uint64, error)

	IBCPacketAck(
		codeID wasmvm.Checksum,
		env wasmvmtypes.Env,
		ack wasmvmtypes.IBCAcknowledgement,
		store wasmvm.KVStore,
		goapi wasmvm.GoAPI,
		querier wasmvm.Querier,
		gasMeter wasmvm.GasMeter,
		gasLimit uint64,
	) (*wasmvmtypes.IBCBasicResponse, uint64, error)

	IBCPacketTimeout(
		codeID wasmvm.Checksum,
		env wasmvmtypes.Env,
		packet wasmvmtypes.IBCPacket,
		store wasmvm.KVStore,
		goapi wasmvm.GoAPI,
		querier wasmvm.Querier,
		gasMeter wasmvm.GasMeter,
		gasLimit uint64,
	) (*wasmvmtypes.IBCBasicResponse, uint64, error)
}

type contractExecutable interface {
	Execute(
		codeID wasmvm.Checksum,
		env wasmvmtypes.Env,
		info wasmvmtypes.MessageInfo,
		executeMsg []byte,
		store wasmvm.KVStore,
		goapi wasmvm.GoAPI,
		querier wasmvm.Querier,
		gasMeter wasmvm.GasMeter,
		gasLimit uint64,
	) (*wasmvmtypes.Response, uint64, error)
}

//MakeIBCInstantiable adds some noop functions to not fail when contract is used for instantiation
func MakeIBCInstantiable(m *MockWasmer) {
	m.CreateFn = HashOnlyCreateFn
	m.InstantiateFn = NoOpInstantiateFn
	m.AnalyzeCodeFn = HasIBCAnalyzeFn
}

// NewIBCContractMockWasmer prepares a mocked wasm_engine for testing with an IBC contract test type.
// It is safe to use the mock with store code and instantiate functions in keeper as is also prepared
// with stubs. Execute is optional. When implemented by the Go test contract then it can be used with
// the mock.
func NewIBCContractMockWasmer(c IBCContractCallbacks) *MockWasmer {
	m := &MockWasmer{
		IBCChannelOpenFn:    c.IBCChannelOpen,
		IBCChannelConnectFn: c.IBCChannelConnect,
		IBCChannelCloseFn:   c.IBCChannelClose,
		IBCPacketReceiveFn:  c.IBCPacketReceive,
		IBCPacketAckFn:      c.IBCPacketAck,
		IBCPacketTimeoutFn:  c.IBCPacketTimeout,
	}
	MakeIBCInstantiable(m)
	if e, ok := c.(contractExecutable); ok { // optional function
		m.ExecuteFn = e.Execute
	}
	return m
}

func HashOnlyCreateFn(code wasmvm.WasmCode) (wasmvm.Checksum, error) {
	if code == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "wasm code must not be nil")
	}
	hash := sha256.Sum256(code)
	return hash[:], nil
}

func NoOpInstantiateFn(wasmvm.Checksum, wasmvmtypes.Env, wasmvmtypes.MessageInfo, []byte, wasmvm.KVStore, wasmvm.GoAPI, wasmvm.Querier, wasmvm.GasMeter, uint64) (*wasmvmtypes.Response, uint64, error) {
	return &wasmvmtypes.Response{}, 0, nil
}

func NoOpCreateFn(_ wasmvm.WasmCode) (wasmvm.Checksum, error) {
	return rand.Bytes(32), nil
}

func HasIBCAnalyzeFn(wasmvm.Checksum) (*wasmvmtypes.AnalysisReport, error) {
	return &wasmvmtypes.AnalysisReport{
		HasIBCEntryPoints: true,
	}, nil
}

func WithoutIBCAnalyzeFn(wasmvm.Checksum) (*wasmvmtypes.AnalysisReport, error) {
	return &wasmvmtypes.AnalysisReport{}, nil
}
