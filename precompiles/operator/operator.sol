pragma solidity >=0.8.17 .0;

/// @dev The OPERATOR contract's address.
address constant OPERATOR_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000000809;

/// @dev The OPERATOR contract's instance.
IOperator constant OPERATOR_CONTRACT = IOperator(
    OPERATOR_PRECOMPILE_ADDRESS
);

/// @author Exocore Team
/// @title Operator Precompile Contract
/// @dev The interface through which solidity contracts will interact with Operator
/// @custom:address 0x0000000000000000000000000000000000000809
interface IOperator {
/// TRANSACTIONS
/// @dev  the oprator register, that will change the state in Operator module
/// Note that this address cannot be a module account.
/// @param clientChainLzId The lzId of client chain
/// @param earningsAddress the earningsAddress of operator
/// @param approveAddress the approveAddress of operator
/// @param clientChainEarningAddr the clientChainEarningAddr of operator
/// @param operatorMetaInfo the  operatorMetaInfo of operator

    function RegisterOperator(
        uint16 clientChainLzId,
        bytes memory earningsAddress,
        bytes memory approveAddress,
        bytes memory clientChainEarningAddr,
        string memory operatorMetaInfo
    ) external returns (bool success);
}
