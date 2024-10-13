package cli

import (
	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdQueryStakerInfos() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-staker-infos [assetID]",
		Short: "shows all staker infos including stakerAddr, validators of that staker, latest balance...",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			assetID := args[0]

			if _, _, err := assetstypes.ValidateID(assetID, true, false); err != nil {
				return err
			}

			request := &types.QueryStakerInfosRequest{
				AssetId: assetID,
			}

			res, err := queryClient.StakerInfos(cmd.Context(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryStakerInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-staker-info [assetID] [stakerAddr]",
		Short: "shows staker info of the specified staker including stakerAddr, validators of that staker, latest balance...",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			assetID := args[0]
			stakerAddr := args[1]

			if _, _, err := assetstypes.ValidateID(assetID, true, false); err != nil {
				return err
			}

			request := &types.QueryStakerInfoRequest{
				AssetId:    assetID,
				StakerAddr: stakerAddr,
			}

			res, err := queryClient.StakerInfo(cmd.Context(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryStakerList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-staker-list",
		Short: "shows staker list including all staker addresses",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.StakerList(cmd.Context(), &types.QueryStakerListRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
