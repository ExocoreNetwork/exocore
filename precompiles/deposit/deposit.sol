pragma solidity >=0.8.17;

/// @dev The DEPOSIT contract's address.
address constant DEPOSIT_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000000804;

/// @dev The DEPOSIT contract's instance.
IDeposit constant DEPOSIT_CONTRACT = IDeposit(
    DEPOSIT_PRECOMPILE_ADDRESS
);

/// @author Exocore Team
/// @title Deposit Precompile Contract
/// @dev The interface through which solidity contracts will interact with Deposit
/// @custom:address 0x0000000000000000000000000000000000000804
interface IDeposit {
/// TRANSACTIONS
/// @dev deposit the client chain assets to the staker, that will change the state in deposit module
/// Note that this address cannot be a module account.
/// @param clientChainLzID The LzID of client chain
/// @param assetsAddress The client chain asset address
/// @param stakerAddress The staker address
/// @param opAmount The deposit amount
    function depositTo(
        uint16 clientChainLzID,
        bytes memory assetsAddress,
        bytes memory stakerAddress,
        uint256 opAmount
    ) external returns (bool success,uint256 latestAssetState);
}

