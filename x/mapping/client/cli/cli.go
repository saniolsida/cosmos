package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/x/mapping/types"
	"github.com/spf13/cobra"
)

func CmdRegisterDevice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-device [device-id] [address]",
		Short: "Register a new device mapping",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			deviceID := args[0]
			address := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgRegisterDevice{
				Creator:  clientCtx.GetFromAddress().String(),
				DeviceId: deviceID,
				Address:  address,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	// ✅ 이 부분이 있어야 --from, --gas 등 사용 가능
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetQueryCmd(moduleName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        moduleName,
		Short:                      "Query commands for the mapping module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdGetAddressFromDeviceID(),
	)

	return cmd
}
