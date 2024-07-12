
pragma solidity >=0.8.17;

import {ASSETS_CONTRACT} from "precompiles/assets/IAssets.sol";

contract DepositCaller {

    event callDepositToResult(bool indexed success, uint256 indexed latestAssetState);
    event ErrorOccurred();

    function testDepositTo(
        uint32 clientChainLzID,
        bytes memory assetsAddress,
        bytes memory stakerAddress,
        uint256 opAmount
    ) public returns (bool, uint256) {
        return
            ASSETS_CONTRACT.depositTo(
            clientChainLzID,
            assetsAddress,
            stakerAddress,
            opAmount
        );
    }

    function testCallDepositToAndEmitEvent(
        uint32 clientChainLzID,
        bytes memory assetsAddress,
        bytes memory stakerAddress,
        uint256 opAmount
    ) public returns (bool, uint256) {
        (bool success,uint256 latestAssetState) = ASSETS_CONTRACT.depositTo(
            clientChainLzID,
            assetsAddress,
            stakerAddress,
            opAmount
        );

        emit callDepositToResult(success, latestAssetState);
        return (success, latestAssetState);
    }

    function testCallDepositToWithTryCatch(
        uint32 clientChainLzID,
        bytes memory assetsAddress,
        bytes memory stakerAddress,
        uint256 opAmount
    ) public returns (bool, uint256) {
        try ASSETS_CONTRACT.depositTo(
            clientChainLzID,
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

    function testWithdrawPrincipal(
        uint32 clientChainLzID,
        bytes memory assetsAddress,
        bytes memory stakerAddress,
        uint256 opAmount
    ) public returns (bool, uint256) {
        return
            ASSETS_CONTRACT.withdrawPrincipal(
            clientChainLzID,
            assetsAddress,
            stakerAddress,
            opAmount
        );
    }

    function testCallWithdrawPrincipalWithTryCatch(
        uint32 clientChainLzID,
        bytes memory assetsAddress,
        bytes memory stakerAddress,
        uint256 opAmount
    ) public returns (bool, uint256) {
        try ASSETS_CONTRACT.withdrawPrincipal(
            clientChainLzID,
            assetsAddress,
            stakerAddress,
            opAmount
        ) returns (bool success, uint256 latestAssetState){
            //call successfully
            emit callDepositToResult(success, latestAssetState);
            return (success, latestAssetState);
        } catch {
            // An error occurred, handle it
            emit ErrorOccurred();
        }
        return (false,0);
    }
}
