package cli

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the wasm module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	queryCmd.AddCommand(
		GetCmdListCode(),
		GetCmdListContractByCode(),
		GetCmdQueryCode(),
		GetCmdGetContractInfo(),
		GetCmdGetContractHistory(),
		GetCmdGetContractState(),
	)
	return queryCmd
}

// GetCmdListCode lists all wasm code uploaded
func GetCmdListCode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-code",
		Short: "List all wasm bytecode on the chain",
		Long:  "List all wasm bytecode on the chain",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(withPageKeyDecoded(cmd.Flags()))
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Codes(
				context.Background(),
				&types.QueryCodesRequest{
					Pagination: pageReq,
				},
			)
			if err != nil {
				return err
			}
			return clientCtx.WithJSONMarshaler(&VanillaStdJsonMarshaller{}).PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "list codes")
	return cmd
}

// GetCmdListContractByCode lists all wasm code uploaded for given code id
func GetCmdListContractByCode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-contract-by-code [code_id]",
		Short: "List wasm all bytecode on the chain for given code id",
		Long:  "List wasm all bytecode on the chain for given code id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			codeID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(withPageKeyDecoded(cmd.Flags()))
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ContractsByCode(
				context.Background(),
				&types.QueryContractsByCodeRequest{
					CodeId:     codeID,
					Pagination: pageReq,
				},
			)
			if err != nil {
				return err
			}
			return clientCtx.WithJSONMarshaler(&VanillaStdJsonMarshaller{}).PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "list contracts by code")
	return cmd
}

// GetCmdQueryCode returns the bytecode for a given contract
func GetCmdQueryCode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "code [code_id] [output filename]",
		Short: "Downloads wasm bytecode for given code id",
		Long:  "Downloads wasm bytecode for given code id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			codeID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Code(
				context.Background(),
				&types.QueryCodeRequest{
					CodeId: codeID,
				},
			)
			if err != nil {
				return err
			}
			if len(res.Data) == 0 {
				return fmt.Errorf("contract not found")
			}

			fmt.Printf("Downloading wasm code to %s\n", args[1])
			return ioutil.WriteFile(args[1], res.Data, 0644)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdGetContractInfo gets details about a given contract
func GetCmdGetContractInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract [bech32_address]",
		Short: "Prints out metadata of a contract given its address",
		Long:  "Prints out metadata of a contract given its address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			_, err = sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ContractInfo(
				context.Background(),
				&types.QueryContractInfoRequest{
					Address: args[0],
				},
			)
			if err != nil {
				return err
			}
			return clientCtx.WithJSONMarshaler(&VanillaStdJsonMarshaller{}).PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdGetContractState dumps full internal state of a given contract
func GetCmdGetContractState() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "contract-state",
		Short:                      "Querying commands for the wasm module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		GetCmdGetContractStateAll(),
		GetCmdGetContractStateRaw(),
		GetCmdGetContractStateSmart(),
	)
	return cmd

}

func GetCmdGetContractStateAll() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all [bech32_address]",
		Short: "Prints out all internal state of a contract given its address",
		Long:  "Prints out all internal state of a contract given its address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			_, err = sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(withPageKeyDecoded(cmd.Flags()))
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.AllContractState(
				context.Background(),
				&types.QueryAllContractStateRequest{
					Address:    args[0],
					Pagination: pageReq,
				},
			)
			if err != nil {
				return err
			}
			return clientCtx.WithJSONMarshaler(&VanillaStdJsonMarshaller{}).PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "contract state")
	return cmd
}

func GetCmdGetContractStateRaw() *cobra.Command {
	decoder := newArgDecoder(hex.DecodeString)
	cmd := &cobra.Command{
		Use:   "raw [bech32_address] [key]",
		Short: "Prints out internal state for key of a contract given its address",
		Long:  "Prints out internal state for of a contract given its address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			_, err = sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			queryData, err := decoder.DecodeString(args[1])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.RawContractState(
				context.Background(),
				&types.QueryRawContractStateRequest{
					Address:   args[0],
					QueryData: queryData,
				},
			)
			if err != nil {
				return err
			}
			return clientCtx.WithJSONMarshaler(&VanillaStdJsonMarshaller{}).PrintProto(res)
		},
	}
	decoder.RegisterFlags(cmd.PersistentFlags(), "key argument")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetCmdGetContractStateSmart() *cobra.Command {
	decoder := newArgDecoder(asciiDecodeString)
	cmd := &cobra.Command{
		Use:   "smart [bech32_address] [query]",
		Short: "Calls contract with given address with query data and prints the returned result",
		Long:  "Calls contract with given address with query data and prints the returned result",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			_, err = sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			if args[1] == "" {
				return errors.New("query data must not be empty")
			}

			queryData, err := decoder.DecodeString(args[1])
			if err != nil {
				return fmt.Errorf("decode query: %s", err)
			}
			if !json.Valid(queryData) {
				return errors.New("query data must be json")
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.SmartContractState(
				context.Background(),
				&types.QuerySmartContractStateRequest{
					Address:   args[0],
					QueryData: queryData,
				},
			)
			if err != nil {
				return err
			}
			return clientCtx.WithJSONMarshaler(&VanillaStdJsonMarshaller{}).PrintProto(res)
		},
	}
	decoder.RegisterFlags(cmd.PersistentFlags(), "query argument")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdGetContractHistory prints the code history for a given contract
func GetCmdGetContractHistory() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract-history [bech32_address]",
		Short: "Prints out the code history for a contract given its address",
		Long:  "Prints out the code history for a contract given its address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			_, err = sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(withPageKeyDecoded(cmd.Flags()))
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ContractHistory(
				context.Background(),
				&types.QueryContractHistoryRequest{
					Address:    args[0],
					Pagination: pageReq,
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.WithJSONMarshaler(&VanillaStdJsonMarshaller{}).PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "contract history")
	return cmd
}

type argumentDecoder struct {
	// dec is the default decoder
	dec                func(string) ([]byte, error)
	asciiF, hexF, b64F bool
}

func newArgDecoder(def func(string) ([]byte, error)) *argumentDecoder {
	return &argumentDecoder{dec: def}
}

func (a *argumentDecoder) RegisterFlags(f *flag.FlagSet, argName string) {
	f.BoolVar(&a.asciiF, "ascii", false, "ascii encoded "+argName)
	f.BoolVar(&a.hexF, "hex", false, "hex encoded  "+argName)
	f.BoolVar(&a.b64F, "b64", false, "base64 encoded "+argName)
}

func (a *argumentDecoder) DecodeString(s string) ([]byte, error) {
	found := -1
	for i, v := range []*bool{&a.asciiF, &a.hexF, &a.b64F} {
		if !*v {
			continue
		}
		if found != -1 {
			return nil, errors.New("multiple decoding flags used")
		}
		found = i
	}
	switch found {
	case 0:
		return asciiDecodeString(s)
	case 1:
		return hex.DecodeString(s)
	case 2:
		return base64.StdEncoding.DecodeString(s)
	default:
		return a.dec(s)
	}
}

func asciiDecodeString(s string) ([]byte, error) {
	return []byte(s), nil
}

type VanillaStdJsonMarshaller struct {
}

func (x VanillaStdJsonMarshaller) MarshalInterfaceJSON(i proto.Message) ([]byte, error) {
	any, err := codectypes.NewAnyWithValue(i)
	if err != nil {
		return nil, err
	}
	return x.MarshalJSON(any)
}

func (x VanillaStdJsonMarshaller) MarshalJSON(o proto.Message) ([]byte, error) {
	return json.MarshalIndent(o, "", " ")
}

func (x VanillaStdJsonMarshaller) MustMarshalJSON(o proto.Message) []byte {
	b, err := x.MarshalJSON(o)
	if err != nil {
		panic(err)
	}
	return b
}

func (x VanillaStdJsonMarshaller) UnmarshalInterfaceJSON(bz []byte, ptr interface{}) error {
	panic("not supported")
}

func (x VanillaStdJsonMarshaller) UnmarshalJSON(bz []byte, ptr proto.Message) error {
	panic("not supported")
}

func (x VanillaStdJsonMarshaller) MustUnmarshalJSON(bz []byte, ptr proto.Message) {
	panic("not supported")
}

// sdk ReadPageRequest expects binary but we encoded to base64 in our marshaller
func withPageKeyDecoded(flagSet *flag.FlagSet) *flag.FlagSet {
	encoded, err := flagSet.GetString(flags.FlagPageKey)
	if err != nil {
		panic(err.Error())
	}
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		panic(err.Error())
	}
	flagSet.Set(flags.FlagPageKey, string(raw))
	return flagSet
}
