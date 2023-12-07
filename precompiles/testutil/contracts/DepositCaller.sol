// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity >=0.8.17;

import "deposit/deposit.sol" as deposit;

contract DepositCaller {

    event callDepositToResult(bool indexed success, uint256 indexed latestAssetState);
    event ErrorOccurred(string errorMessage);

    function testDepositTo(
        uint16 clientChainLzId,
        bytes memory assetsAddress,
        bytes memory stakerAddress,
        uint256 opAmount
    ) public returns (bool, uint256) {
        return
            deposit.DEPOSIT_CONTRACT.depositTo(
            clientChainLzId,
            assetsAddress,
            stakerAddress,
            opAmount
        );
    }

    function testCallDepositToAndEmitEvent(
        uint16 clientChainLzId,
        bytes memory assetsAddress,
        bytes memory stakerAddress,
        uint256 opAmount
    ) public returns (bool, uint256) {
        (bool success,uint256 latestAssetState) = deposit.DEPOSIT_CONTRACT.depositTo(
            clientChainLzId,
            assetsAddress,
            stakerAddress,
            opAmount
        );

        emit callDepositToResult(success, latestAssetState);
        return (success, latestAssetState);
    }

    function testCallDepositToWithTryCatch(
        uint16 clientChainLzId,
        bytes memory assetsAddress,
        bytes memory stakerAddress,
        uint256 opAmount
    ) public returns (bool, uint256) {
        try deposit.DEPOSIT_CONTRACT.depositTo(
            clientChainLzId,
            assetsAddress,
            stakerAddress,
            opAmount
        ) returns (bool success, uint256 latestAssetState){
            //call successfully
            emit callDepositToResult(success, latestAssetState);
            return (success, latestAssetState);
        }catch Error(string memory errorMessage){
            // An error occurred, handle it
            emit ErrorOccurred(errorMessage);
        }
        return (false,0);
    }
/*    function testDelegateCallDepositTo(
        uint16 clientChainLzId,
        bytes memory assetsAddress,
        bytes memory stakerAddress,
        uint256 opAmount
    ) public returns (bool, uint256) {
        (bool success,uint256 latestAssetState) = deposit.DEPOSIT_PRECOMPILE_ADDRESS.delegatecall(
            abi.encodeWithSignature(
                "depositTo(uint16,bytes,bytes,uint256)",
                clientChainLzId,
                assetsAddress,
                stakerAddress,
                opAmount
            )
        );
        require(success, "failed delegateCall to precompile");
        return (success, latestAssetState);
    }

    function testStaticCallDepositTo(
        uint16 clientChainLzId,
        bytes memory assetsAddress,
        bytes memory stakerAddress,
        uint256 opAmount
    ) public view {
        (bool success,) = deposit
            .DEPOSIT_PRECOMPILE_ADDRESS
            .staticcall(
            abi.encodeWithSignature(
                "depositTo(uint16,bytes,bytes,uint256)",
                clientChainLzId,
                assetsAddress,
                stakerAddress,
                opAmount
            )
        );
        require(success, "failed staticCall to precompile");
    }*/
}
