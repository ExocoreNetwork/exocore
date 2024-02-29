pragma solidity >=0.8.17 .0;

/// @dev The WITHDRAW contract's address.
address constant WITHDRAW_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000000808;

/// @dev The WITHDRAW contract's instance.
IWithdraw constant WITHDRAW_CONTRACT = IWithdraw(
    WITHDRAW_PRECOMPILE_ADDRESS
);

/// @author Exocore Team
/// @title WITHDRAW Precompile Contract
/// @dev The interface through which solidity contracts will interact with WITHDRAW
/// @custom:address 0x0000000000000000000000000000000000000808
interface IWithdraw {
/// TRANSACTIONS
/// @dev withdraw To the staker, that will change the state in withdraw module
/// Note that this address cannot be a module account.
/// @param clientChainLzID The LzID of client chain
/// @param assetsAddress The client chain asset Address
/// @param withdrawAddress The withdraw address
/// @param opAmount The withdraw amount
    function withdrawPrinciple(
        uint16 clientChainLzID,
        bytes memory assetsAddress,
        bytes memory withdrawAddress,
        uint256 opAmount
    ) external returns (bool success,uint256 latestAssetState);
}