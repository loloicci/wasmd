package simulation

import (
	"fmt"
	"math/rand"

	"github.com/Finschia/finschia-sdk/codec"
	simtypes "github.com/Finschia/finschia-sdk/types/simulation"
	"github.com/Finschia/finschia-sdk/x/simulation"

	"github.com/Finschia/wasmd/x/wasm/types"
)

func ParamChanges(r *rand.Rand, cdc codec.Codec) []simtypes.ParamChange {
	params := types.DefaultParams()
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyUploadAccess),
			func(r *rand.Rand) string {
				jsonBz, err := cdc.MarshalJSON(&params.CodeUploadAccess)
				if err != nil {
					panic(err)
				}
				return string(jsonBz)
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyInstantiateAccess),
			func(r *rand.Rand) string {
				return fmt.Sprintf("%q", params.CodeUploadAccess.Permission.String())
			},
		),
	}
}

func RandomParams(r *rand.Rand) types.Params {
	permissionType := types.AccessType(simtypes.RandIntBetween(r, 1, 3))
	account, _ := simtypes.RandomAcc(r, simtypes.RandomAccounts(r, 10))
	accessConfig := permissionType.With(account.Address)
	return types.Params{
		CodeUploadAccess:             accessConfig,
		InstantiateDefaultPermission: accessConfig.Permission,
	}
}
