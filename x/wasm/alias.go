// nolint
// autogenerated code using github.com/rigelrozanski/multitool
// aliases generated for the following subdirectories:
// ALIASGEN: github.com/cosmwasm/wasmd/x/wasm/types
// ALIASGEN: github.com/cosmwasm/wasmd/x/wasm/keeper
package wasm

import (
	"github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/types"
)

const (
	firstCodeID                     = 1
	DefaultParamspace               = types.DefaultParamspace
	ModuleName                      = types.ModuleName
	StoreKey                        = types.StoreKey
	TStoreKey                       = types.TStoreKey
	QuerierRoute                    = types.QuerierRoute
	RouterKey                       = types.RouterKey
	MaxWasmSize                     = types.MaxWasmSize
	MaxLabelSize                    = types.MaxLabelSize
	BuildTagRegexp                  = types.BuildTagRegexp
	MaxBuildTagSize                 = types.MaxBuildTagSize
	CustomEventType                 = types.CustomEventType
	AttributeKeyContractAddr        = types.AttributeKeyContractAddr
	ProposalTypeStoreCode           = types.ProposalTypeStoreCode
	ProposalTypeInstantiateContract = types.ProposalTypeInstantiateContract
	ProposalTypeMigrateContract     = types.ProposalTypeMigrateContract
	ProposalTypeUpdateAdmin         = types.ProposalTypeUpdateAdmin
	ProposalTypeClearAdmin          = types.ProposalTypeClearAdmin
	GasMultiplier                   = keeper.GasMultiplier
	MaxGas                          = keeper.MaxGas
	QueryListContractByCode         = keeper.QueryListContractByCode
	QueryGetContract                = keeper.QueryGetContract
	QueryGetContractState           = keeper.QueryGetContractState
	QueryGetCode                    = keeper.QueryGetCode
	QueryListCode                   = keeper.QueryListCode
	QueryMethodContractStateSmart   = keeper.QueryMethodContractStateSmart
	QueryMethodContractStateAll     = keeper.QueryMethodContractStateAll
	QueryMethodContractStateRaw     = keeper.QueryMethodContractStateRaw
)

var (
	// functions aliases
	RegisterCodec             = types.RegisterLegacyAminoCodec
	RegisterInterfaces        = types.RegisterInterfaces
	ValidateGenesis           = types.ValidateGenesis
	ConvertToProposals        = types.ConvertToProposals
	GetCodeKey                = types.GetCodeKey
	GetContractAddressKey     = types.GetContractAddressKey
	GetContractStorePrefixKey = types.GetContractStorePrefix
	NewCodeInfo               = types.NewCodeInfo
	NewAbsoluteTxPosition     = types.NewAbsoluteTxPosition
	NewContractInfo           = types.NewContractInfo
	NewEnv                    = types.NewEnv
	NewWasmCoins              = types.NewWasmCoins
	ParseEvents               = types.ParseEvents
	DefaultWasmConfig         = types.DefaultWasmConfig
	DefaultParams             = types.DefaultParams
	InitGenesis               = keeper.InitGenesis
	ExportGenesis             = keeper.ExportGenesis
	NewMessageHandler         = keeper.NewDefaultMessageHandler
	DefaultEncoders           = keeper.DefaultEncoders
	EncodeBankMsg             = keeper.EncodeBankMsg
	NoCustomMsg               = keeper.NoCustomMsg
	EncodeStakingMsg          = keeper.EncodeStakingMsg
	EncodeWasmMsg             = keeper.EncodeWasmMsg
	NewKeeper                 = keeper.NewKeeper
	NewLegacyQuerier          = keeper.NewLegacyQuerier
	DefaultQueryPlugins       = keeper.DefaultQueryPlugins
	BankQuerier               = keeper.BankQuerier
	NoCustomQuerier           = keeper.NoCustomQuerier
	StakingQuerier            = keeper.StakingQuerier
	WasmQuerier               = keeper.WasmQuerier
	CreateTestInput           = keeper.CreateTestInput
	TestHandler               = keeper.TestHandler
	NewWasmProposalHandler    = keeper.NewWasmProposalHandler
	NewQuerier                = keeper.NewQuerier
	ContractFromPortID        = keeper.ContractFromPortID
	WithWasmEngine            = keeper.WithWasmEngine

	// variable aliases
	ModuleCdc            = types.ModuleCdc
	DefaultCodespace     = types.DefaultCodespace
	ErrCreateFailed      = types.ErrCreateFailed
	ErrAccountExists     = types.ErrAccountExists
	ErrInstantiateFailed = types.ErrInstantiateFailed
	ErrExecuteFailed     = types.ErrExecuteFailed
	ErrGasLimit          = types.ErrGasLimit
	ErrInvalidGenesis    = types.ErrInvalidGenesis
	ErrNotFound          = types.ErrNotFound
	ErrQueryFailed       = types.ErrQueryFailed
	ErrInvalidMsg        = types.ErrInvalidMsg
	KeyLastCodeID        = types.KeyLastCodeID
	KeyLastInstanceID    = types.KeyLastInstanceID
	CodeKeyPrefix        = types.CodeKeyPrefix
	ContractKeyPrefix    = types.ContractKeyPrefix
	ContractStorePrefix  = types.ContractStorePrefix
	EnableAllProposals   = types.EnableAllProposals
	DisableAllProposals  = types.DisableAllProposals
)

type (
	ProposalType                   = types.ProposalType
	GenesisState                   = types.GenesisState
	Code                           = types.Code
	Contract                       = types.Contract
	MsgStoreCode                   = types.MsgStoreCode
	MsgStoreCodeResponse           = types.MsgStoreCodeResponse
	MsgInstantiateContract         = types.MsgInstantiateContract
	MsgInstantiateContractResponse = types.MsgInstantiateContractResponse
	MsgExecuteContract             = types.MsgExecuteContract
	MsgExecuteContractResponse     = types.MsgExecuteContractResponse
	MsgMigrateContract             = types.MsgMigrateContract
	MsgMigrateContractResponse     = types.MsgMigrateContractResponse
	MsgUpdateAdmin                 = types.MsgUpdateAdmin
	MsgUpdateAdminResponse         = types.MsgUpdateAdminResponse
	MsgClearAdmin                  = types.MsgClearAdmin
	MsgWasmIBCCall                 = types.MsgIBCSend
	MsgClearAdminResponse          = types.MsgClearAdminResponse
	MsgServer                      = types.MsgServer
	Model                          = types.Model
	CodeInfo                       = types.CodeInfo
	ContractInfo                   = types.ContractInfo
	CreatedAt                      = types.AbsoluteTxPosition
	Config                         = types.WasmConfig
	ContractInfoWithAddress        = types.ContractInfoWithAddress
	CodeInfoResponse               = types.CodeInfoResponse
	MessageHandler                 = keeper.SDKMessageHandler
	BankEncoder                    = keeper.BankEncoder
	CustomEncoder                  = keeper.CustomEncoder
	StakingEncoder                 = keeper.StakingEncoder
	WasmEncoder                    = keeper.WasmEncoder
	MessageEncoders                = keeper.MessageEncoders
	Keeper                         = keeper.Keeper
	QueryHandler                   = keeper.QueryHandler
	CustomQuerier                  = keeper.CustomQuerier
	QueryPlugins                   = keeper.QueryPlugins
	Option                         = keeper.Option
)
