// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity >=0.8.17 .0;

/// @dev The DEPOSIT contract's address.
address constant DEPOSIT_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000000804;

/// @dev The DEPOSIT contract's instance.
DepositI constant DEPOSIT_CONTRACT = DepositI(
    DEPOSIT_PRECOMPILE_ADDRESS
);

/// @author Exocore Team
/// @title Deposit Precompile Contract
/// @dev The interface through which solidity contracts will interact with Deposit
/// @custom:address 0x0000000000000000000000000000000000000804
interface DepositI {
/// TRANSACTIONS
/// @dev deposit the client chain assets to the staker, that will change the state in deposit module
/// Note that this address cannot be a module account.
/// @param ClientChainLzId The lzId of client chain
/// @param AssetsAddress The client chain asset Address
/// @param StakerAddress The staker address
/// @param OpAmount The deposit amount
    function DepositTo(
        uint64 ClientChainLzId,
        bytes memory AssetsAddress,
        bytes memory StakerAddress,
        uint256 OpAmount
    ) external returns (bool success);
}

