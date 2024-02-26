package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	"github.com/spf13/cast"
)

func CmdListRoundData() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-round-data",
		Short: "list all round-data",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllRoundDataRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.RoundDataAll(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowRoundData() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-round-data [token-id]",
		Short: "shows a round-data",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argTokenId, err := cast.ToInt32E(args[0])
			if err != nil {
				return err
			}

			params := &types.QueryGetRoundDataRequest{
				TokenId: argTokenId,
			}

			res, err := queryClient.RoundData(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
