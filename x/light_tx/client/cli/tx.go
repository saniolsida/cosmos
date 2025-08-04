package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	lighttype "github.com/cosmos/cosmos-sdk/x/light_tx/types"
	"github.com/spf13/cobra"
)

func CmdSendLightTx() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-light-tx [...fields]",
		Short: "Send a MsgSendLightTx with either SolarData or RECMeta",
		Args:  cobra.MinimumNArgs(8),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var payload interface{}

			if len(args) == 6 {
				// SolarData 전송
				totalEnergy, err := strconv.ParseFloat(args[2], 64)
				if err != nil {
					return fmt.Errorf("invalid total_energy: %v", err)
				}

				lat, err := strconv.ParseFloat(args[3], 64)
				if err != nil {
					return fmt.Errorf("invalid latitude: %v", err)
				}
				lon, err := strconv.ParseFloat(args[4], 64)
				if err != nil {
					return fmt.Errorf("invalid longitude: %v", err)
				}

				payload = &lighttype.MsgSendLightTx_Original{
					Original: &lighttype.SolarData{
						DeviceId:    args[0],
						Timestamp:   args[1],
						TotalEnergy: totalEnergy,
						Location: &lighttype.Location{
							Latitude:   lat,
							Longitutde: lon,
						},
					},
				}
			} else if len(args) >= 13 {
				// RECMeta 전송
				payload = &lighttype.MsgSendLightTx_Rec{
					Rec: &lighttype.RECMeta{
						FacilityId:          args[0],
						FacilityName:        args[1],
						Location:            args[2],
						TechnologyType:      args[3],
						CapacityMw:          args[4],
						RegistrationDate:    args[5],
						CertifiedId:         args[6],
						IssueData:           args[7],
						GenerationStartDate: args[8],
						GenerationEndDate:   args[9],
						MeasuredVolume_MWh:  args[10],
						RetiredDate:         args[11],
						RetirementPurpose:   args[12],
						Status:              args[13],
						Timestamp:           args[14],
					},
				}

			} else {
				return fmt.Errorf("insufficient arguments: need 8 for SolarData or 13+ for RECMeta")
			}

			hash := args[len(args)-3]
			signature := args[len(args)-2]
			pubkey := args[len(args)-1]

			msg := &lighttype.MsgSendLightTx{
				Creator:   clientCtx.GetFromAddress().String(),
				Hash:      hash,
				Signature: signature,
				Pubkey:    pubkey,
			}

			switch p := payload.(type) {
			case *lighttype.MsgSendLightTx_Original:
				msg.Payload = p
			case *lighttype.MsgSendLightTx_Rec:
				msg.Payload = p
			default:
				return fmt.Errorf("invalid payload type")
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
