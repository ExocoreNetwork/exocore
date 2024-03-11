pragma solidity >=0.8.17;

/// @dev The delegation contract's address.
address constant DELEGATION_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000000805;

/// @dev The delegation contract's instance.
IDelegation constant DELEGATION_CONTRACT = IDelegation(
    DELEGATION_PRECOMPILE_ADDRESS
);

/// @author Exocore Team
/// @title delegation Precompile Contract
/// @dev The interface through which solidity contracts will interact with delegation
/// @custom:address 0x0000000000000000000000000000000000000805
interface IDelegation {
/// TRANSACTIONS
/// @dev delegate the client chain assets to the operator through client chain, that will change the states in delegation and assets module
/// Note that this address cannot be a module account.
/// @param clientChainLzID The LzID of client chain
/// @param lzNonce The cross chain tx layerZero nonce
/// @param assetsAddress The client chain asset Address
/// @param stakerAddress The staker address
/// @param operatorAddr  The operator address that wants to be delegated to
/// @param opAmount The delegation amount
    function delegateToThroughClientChain(
        uint16 clientChainLzID,
        uint64 lzNonce,
        bytes memory assetsAddress,
        bytes memory stakerAddress,
        bytes memory operatorAddr,
        uint256 opAmount
    ) external returns (bool success);

/// TRANSACTIONS
/// @dev undelegate the client chain assets from the operator through client chain, that will change the states in delegation and assets module
/// Note that this address cannot be a module account.
/// @param clientChainLzID The LzID of client chain
/// @param lzNonce The cross chain tx layerZero nonce
/// @param assetsAddress The client chain asset Address
/// @param stakerAddress The staker address
/// @param operatorAddr  The operator address that wants to unDelegate from
/// @param opAmount The Undelegation amount
    function undelegateFromThroughClientChain(
        uint16 clientChainLzID,
        uint64 lzNonce,
        bytes memory assetsAddress,
        bytes memory stakerAddress,
        bytes memory operatorAddr,
        uint256 opAmount
    ) external returns (bool success);
}

