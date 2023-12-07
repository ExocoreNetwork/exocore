// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity >=0.8.17 .0;

/// @dev The WITHDRAW contract's address.
address constant WITHDRAW_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000000807;

/// @dev The WITHDRAW contract's instance.
IWithdraw constant WITHDRAW_CONTRACT = IWithdraw(
    WITHDRAW_PRECOMPILE_ADDRESS
);

/// @author Exocore Team
/// @title WITHDRAW Precompile Contract
/// @dev The interface through which solidity contracts will interact with WITHDRAW
/// @custom:address 0x0000000000000000000000000000000000000807
interface IWithdraw {
/// TRANSACTIONS
/// @dev withdraw To the staker, that will change the state in withdraw module
/// Note that this address cannot be a module account.
/// @param ClientChainLzId The lzId of client chain
/// @param AssetsAddress The client chain asset Address
/// @param WithdrawAddress The withdraw address
/// @param OpAmount The withdraw amount
    function withdrawTo(
        uint16 ClientChainLzId,
        bytes memory AssetsAddress,
        bytes memory WithdrawAddress,
        uint256 OpAmount
    ) external returns (bool success);
}

