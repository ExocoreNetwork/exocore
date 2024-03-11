package operator

import (
	"fmt"
	"github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
	"reflect"
)

type operatorParams struct {
	ClientChainLzId        uint64
	EarningsAddr           []byte
	ApproveAddr            []byte
	ClientChainEarningAddr string
	operatorMetaInfo       string
}

func (p Precompile) GetOperatorParamsFromInputs(ctx sdk.Context, args []interface{}) (*operatorParams, error) {
	if len(args) != 8 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 4, len(args))
	}
	registerReq := &operatorParams{}

	clientChainLzID, ok := args[0].(uint16)
	if !ok {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 0, reflect.TypeOf(args[0]), clientChainLzID)
	}
	registerReq.ClientChainLzId = uint64(clientChainLzID)

	// the length of client chain address inputted by caller is 32, so we need to check the length and remove the padding according to the actual length.
	assetAddr, ok := args[1].([]byte)
	if !ok || assetAddr == nil {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 1, reflect.TypeOf(args[0]), assetAddr)
	}
	if len(assetAddr) != types.GeneralClientChainAddrLength {
		return nil, fmt.Errorf(ErrInputClientChainAddrLength, len(assetAddr), types.GeneralClientChainAddrLength)
	}
	registerReq.EarningsAddr = assetAddr[:]

	stakerAddr, ok := args[2].([]byte)
	if !ok || stakerAddr == nil {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 2, reflect.TypeOf(args[0]), stakerAddr)
	}
	if len(assetAddr) != types.GeneralClientChainAddrLength {
		return nil, fmt.Errorf(ErrInputClientChainAddrLength, len(assetAddr), types.GeneralClientChainAddrLength)
	}
	registerReq.ApproveAddr = stakerAddr[:]

	registerReq.ClientChainEarningAddr = args[3].(string)

	registerReq.operatorMetaInfo = args[4].(string)

	return registerReq, nil
}
