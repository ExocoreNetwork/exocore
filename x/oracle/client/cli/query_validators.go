package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	"github.com/spf13/cast"
)

func CmdListValidators() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-validators",
		Short: "list all validators",
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

			params := &types.QueryAllValidatorsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.ValidatorsAll(cmd.Context(), params)
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

func CmdShowValidators() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-validators [block]",
		Short: "shows a validators",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argBlock, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			params := &types.QueryGetValidatorsRequest{
				Block: argBlock,
			}

			res, err := queryClient.Validators(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
