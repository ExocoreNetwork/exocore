package reward

import (
	"fmt"
	"math/big"

	exocmn "github.com/ExocoreNetwork/exocore/precompiles/common"

	sdkmath "cosmossdk.io/math"

	"github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/ExocoreNetwork/exocore/x/reward/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
)

func (p Precompile) GetRewardParamsFromInputs(ctx sdk.Context, args []interface{}) (*keeper.RewardParams, error) {
	if len(args) != 4 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 4, len(args))
	}
	rewardParams := &keeper.RewardParams{}
	clientChainLzID, ok := args[0].(uint32)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 0, "uint32", clientChainLzID)
	}
	rewardParams.ClientChainLzID = uint64(clientChainLzID)

	info, err := p.assetsKeeper.GetClientChainInfoByIndex(ctx, rewardParams.ClientChainLzID)
	if err != nil {
		return nil, err
	}
	clientChainAddrLength := info.AddressLength

	// the length of client chain address inputted by caller is 32, so we need to check the length and remove the padding according to the actual length.
	assetAddr, ok := args[1].([]byte)
	if !ok || assetAddr == nil {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 1, "[]byte", assetAddr)
	}
	if len(assetAddr) != types.GeneralClientChainAddrLength {
		return nil, fmt.Errorf(exocmn.ErrInvalidAddrLength, len(assetAddr), types.GeneralClientChainAddrLength)
	}
	rewardParams.AssetsAddress = assetAddr[:clientChainAddrLength]

	stakerAddr, ok := args[2].([]byte)
	if !ok || stakerAddr == nil {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 2, "[]byte", stakerAddr)
	}
	if len(assetAddr) != types.GeneralClientChainAddrLength {
		return nil, fmt.Errorf(exocmn.ErrInvalidAddrLength, len(assetAddr), types.GeneralClientChainAddrLength)
	}
	rewardParams.WithdrawRewardAddress = stakerAddr[:clientChainAddrLength]

	opAmount, ok := args[3].(*big.Int)
	if !ok || opAmount == nil || opAmount.Cmp(big.NewInt(0)) == 0 {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 3, "*big.Int", opAmount)
	}

	rewardParams.OpAmount = sdkmath.NewIntFromBigInt(opAmount)
	return rewardParams, nil
}
