package cli

import (
	types "github.com/ExocoreNetwork/exocore/x/exomint/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	FlagMintDenom       = "mint-denom"
	FlagEpochReward     = "epoch-reward"
	FlagEpochIdentifier = "epoch-identifier"
)

// NewTxCmd returns a root CLI command handler for deposit commands
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "exomint subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		CmdUpdateParams(),
	)
	return txCmd
}

// CmdUpdateParams returns a CLI command handler for creating a MsgUpdateParams transaction.
func CmdUpdateParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-params",
		Short: "updata the parameters of the module",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf, err := tx.NewFactoryCLI(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			msg, err := newBuildUpdateParamsMsg(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			// this calls ValidateBasic internally so we don't need to do that.
			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	f := cmd.Flags()
	f.String(
		FlagMintDenom, "", "The mint denomination",
	)
	f.String(
		FlagEpochReward, "", "The amount of the mint denomination to mint, per epoch (as a string)",
	)
	f.String(
		FlagEpochIdentifier, "", "The identifier of the epoch at which it should be minted",
	)

	// transaction level flags from the SDK
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func newBuildUpdateParamsMsg(
	clientCtx client.Context, fs *pflag.FlagSet,
) (*types.MsgUpdateParams, error) {
	sender := clientCtx.GetFromAddress()
	// #nosec G703 // this only errors if the flag isn't defined.
	mintDenom, _ := fs.GetString(FlagMintDenom)
	// #nosec G703 // this only errors if the flag isn't defined.
	epochIdentifier, _ := fs.GetString(FlagEpochIdentifier)
	// #nosec G703 // this only errors if the flag isn't defined.
	epochRewardStr, _ := fs.GetString(FlagEpochReward)
	res, ok := sdk.NewIntFromString(epochRewardStr)
	if !ok {
		res = sdk.NewInt(0)
	}
	msg := &types.MsgUpdateParams{
		Authority: sender.String(),
		Params: types.Params{
			MintDenom:       mintDenom,
			EpochReward:     res,
			EpochIdentifier: epochIdentifier,
		},
	}
	return msg, nil
}
