package cli

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"fmt"
	restakingtype "github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

// NewTxCmd returns a root CLI command handler for deposit commands
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        restakingtype.ModuleName,
		Short:                      "restaking subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		RegisterClientChain(),
		RegisterAsset(),
	)
	return txCmd
}

// RegisterClientChain register client chain
// todo: this function should be controlled by governance in the future
func RegisterClientChain() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "RegisterClientChain Name MetaInfo LZChainId AddressLength",
		Short: "register client chain",
		Args:  cobra.MinimumNArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := cliCtx.GetFromAddress()
			msg := &restakingtype.RegisterClientChainReq{
				FromAddress: sender.String(),
				Info: &restakingtype.ClientChainInfo{
					Name:     args[0],
					MetaInfo: args[1],
				},
			}
			lzChainId, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return errorsmod.Wrap(restakingtype.ErrCliCmdInputArg, fmt.Sprintf("error arg is:%v", args[2]))
			}
			addressLength, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return errorsmod.Wrap(restakingtype.ErrCliCmdInputArg, fmt.Sprintf("error arg is:%v", args[3]))
			}
			msg.Info.LayerZeroChainId = lzChainId
			msg.Info.AddressLength = uint32(addressLength)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// RegisterAsset register the asset on the client chain
// todo: this function should be controlled by governance in the future
func RegisterAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "RegisterAsset Name Symbol Address MetaInfo TotalSupply LZChainId Decimals",
		Short: "register asset",
		Args:  cobra.MinimumNArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := cliCtx.GetFromAddress()
			msg := &restakingtype.RegisterAssetReq{
				FromAddress: sender.String(),
				Info: &restakingtype.AssetInfo{
					Name:     args[0],
					Symbol:   args[1],
					Address:  strings.ToLower(args[2]),
					MetaInfo: args[3],
				},
			}
			totalSupply, ok := sdkmath.NewIntFromString(args[4])
			if !ok {
				return errorsmod.Wrap(restakingtype.ErrCliCmdInputArg, fmt.Sprintf("error arg is:%v", args[4]))
			}

			lzChainId, err := strconv.ParseUint(args[5], 10, 64)
			if err != nil {
				return errorsmod.Wrap(restakingtype.ErrCliCmdInputArg, fmt.Sprintf("error arg is:%v", args[5]))
			}
			decimal, err := strconv.ParseUint(args[6], 10, 64)
			if err != nil {
				return errorsmod.Wrap(restakingtype.ErrCliCmdInputArg, fmt.Sprintf("error arg is:%v", args[6]))
			}

			msg.Info.TotalSupply = totalSupply
			msg.Info.LayerZeroChainId = lzChainId
			msg.Info.Decimals = uint32(decimal)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
