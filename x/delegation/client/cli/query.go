package cli

import (
	"context"
	"strconv"
	"strings"

	errorsmod "cosmossdk.io/errors"
	"github.com/ExocoreNetwork/exocore/x/assets/types"
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

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
		QueryUndelegations(),
		QueryWaitCompleteUndelegations(),
		QueryUndelegationHoldCount(),
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
				return errorsmod.Wrap(types.ErrInvalidCliCmdArg, err.Error())
			}
			stakerID, assetID := types.GetStakeIDAndAssetIDFromStr(clientChainLzID, args[1], args[2])
			accAddr, err := sdk.AccAddressFromBech32(args[3])
			if err != nil {
				return errorsmod.Wrap(types.ErrInvalidCliCmdArg, err.Error())
			}
			req := &delegationtype.SingleDelegationInfoReq{
				StakerID:     stakerID,         // already lowercase
				AssetID:      assetID,          // already lowercase
				OperatorAddr: accAddr.String(), // already lowercase
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
			stakerID := strings.ToLower(args[0])
			if _, _, err := types.ValidateID(stakerID, false, false); err != nil {
				return errorsmod.Wrap(types.ErrInvalidCliCmdArg, err.Error())
			}
			assetID := strings.ToLower(args[1])
			if _, _, err := types.ValidateID(assetID, false, false); err != nil {
				return errorsmod.Wrap(types.ErrInvalidCliCmdArg, err.Error())
			}
			queryClient := delegationtype.NewQueryClient(clientCtx)
			req := &delegationtype.DelegationInfoReq{
				StakerID: strings.ToLower(stakerID),
				AssetID:  strings.ToLower(assetID),
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

// QueryUndelegations queries all undelegations for staker and asset
func QueryUndelegations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "QueryUndelegations stakerID assetID",
		Short: "Get undelegations",
		Long:  "Get undelegations",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := delegationtype.NewQueryClient(clientCtx)
			req := &delegationtype.UndelegationsReq{
				StakerID: args[0],
				AssetID:  args[1],
			}
			res, err := queryClient.QueryUndelegations(context.Background(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// QueryWaitCompleteUndelegations queries all undelegations waiting to be completed by height
func QueryWaitCompleteUndelegations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "QueryWaitCompleteUndelegations height",
		Short: "Get undelegations waiting to be completed",
		Long:  "Get undelegations waiting to be completed",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			height, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			queryClient := delegationtype.NewQueryClient(clientCtx)
			req := &delegationtype.WaitCompleteUndelegationsReq{
				BlockHeight: height,
			}
			res, err := queryClient.QueryWaitCompleteUndelegations(context.Background(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// QueryUndelegationHoldCount queries undelegation hold count for a record key.
func QueryUndelegationHoldCount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "QueryUndelegationHoldCount recordKey",
		Short: "Get undelegation hold count",
		Long:  "Get undelegation hold count",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := delegationtype.NewQueryClient(clientCtx)
			req := &delegationtype.UndelegationHoldCountReq{
				RecordKey: strings.ToLower(args[0]),
			}
			res, err := queryClient.QueryUndelegationHoldCount(context.Background(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
