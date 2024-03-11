package operator

import (
	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

const (
	// MethodRegisterOperator defines the ABI method name for the operator
	//  transaction.
	MethodRegisterOperator = "RegisterOperator"
)

// RegisterOperator
func (p Precompile) RegisterOperator(ctx sdk.Context, origin common.Address, contract *vm.Contract, stateDB vm.StateDB, method *abi.Method, args []interface{}) ([]byte, error) {
	// check the invalidation of caller contract
	//Registration operator qualification review
	operatorParams, err := p.GetOperatorParamsFromInputs(ctx, args)
	if err != nil {
		return nil, err
	}
	registerReq := &operatortypes.RegisterOperatorReq{
		FromAddress: contract.CallerAddress.String(),
		Info: &operatortypes.OperatorInfo{
			EarningsAddr:     string(operatorParams.EarningsAddr),
			ApproveAddr:      string(operatorParams.ApproveAddr),
			OperatorMetaInfo: operatorParams.operatorMetaInfo,
			ClientChainEarningsAddr: &operatortypes.ClientChainEarningAddrList{
				EarningInfoList: []*operatortypes.ClientChainEarningAddrInfo{
					{operatorParams.ClientChainLzId, operatorParams.ClientChainEarningAddr},
				},
			},
		},
	}
	_, err = p.operatorKeeper.RegisterOperator(ctx, registerReq)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}
