package cli

import (
	"context"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	"github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the parent command for all incentives CLI query commands.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        delegationtype.ModuleName,
		Short:                      "Querying commands for the delegation module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		QuerySingleDelegationInfo(),
		QueryDelegationInfo(),
	)
	return cmd
}

// QuerySingleDelegationInfo queries the single delegation info
func QuerySingleDelegationInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "QuerySingleDelegationInfo clientChainID stakerAddr assetAddr operatorAddr",
		Short: "Get single delegation info",
		Long:  "Get single delegation info",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := delegationtype.NewQueryClient(clientCtx)
			clientChainLzID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return errorsmod.Wrap(types.ErrCliCmdInputArg, err.Error())
			}
			stakerID, assetID := types.GetStakeIDAndAssetIDFromStr(clientChainLzID, args[1], args[2])
			req := &delegationtype.SingleDelegationInfoReq{
				StakerID:     stakerID,
				AssetID:      assetID,
				OperatorAddr: args[3],
			}
			res, err := queryClient.QuerySingleDelegationInfo(context.Background(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// QueryDelegationInfo queries delegation info
func QueryDelegationInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "QueryDelegationInfo stakerID assetID",
		Short: "Get delegation info",
		Long:  "Get delegation info",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := delegationtype.NewQueryClient(clientCtx)
			req := &delegationtype.DelegationInfoReq{
				StakerID: args[0],
				AssetID:  args[1],
			}
			res, err := queryClient.QueryDelegationInfo(context.Background(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
