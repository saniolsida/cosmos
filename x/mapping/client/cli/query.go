package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/mapping/types"
	"github.com/spf13/cobra"
)

func CmdGetAddressFromDeviceID() *cobra.Command {
	return &cobra.Command{
		Use:   "get-address [device-id]",
		Short: "Query the address mapped to a given device ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			deviceID := args[0]

			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.GetAddressFromDeviceID(
				context.Background(),
				&types.QueryGetAddressFromDeviceIDRequest{DeviceId: deviceID},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}
