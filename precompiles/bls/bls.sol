pragma solidity >=0.8.17 .0;

/// @dev The BLS contract's address.
address constant BLS_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000000902;

/// @dev The BLS contract's instance.
IBLS constant BLS_CONTRACT = IBLS(
    BLS_PRECOMPILE_ADDRESS
);

/// @author Exocore Team
/// @title BLS Precompile Contract
/// @dev The interface through which solidity contracts will interact with BLS
/// @custom:address 0x0000000000000000000000000000000000000902

interface IBLS {
/// TRANSACTIONS
/// @dev Called by the avs manager service register an operator as the owner of a BLS public key.
/// @param operator is the operator for whom the key is being registered
/// @param pubKeyRegistrationParams contains the G1 & G2 public keys of the operator, and a signature proving their ownership
/// @param pubKeyRegistrationMessageHash is a hash that the operator must sign to prove key ownership
    function registerBLSPublicKey(
        address operator,
        bytes calldata pubKeyRegistrationParams,
        bytes calldata pubKeyRegistrationMessageHash
    ) external returns (bool success,bytes32 operatorId);

    /// TRANSACTIONS
/// @dev Returns the pubkey and pubkey hash of an operator,Reverts if the operator has not registered a valid pubkey
/// @param operator is the operator for whom the key is being registered
    function getRegisteredPubkey(address operator) external returns (bytes32,bytes32);

    /// TRANSACTIONS
/// @dev Get the count of the current task
    function checkSignatures(
        bytes32 msgHash,
        bytes calldata quorumNumbers,
        uint32 referenceBlockNumber,
        bytes memory params
    ) external view returns (bytes32,bytes32);

    /// TRANSACTIONS
/// @dev Get the task window block for the current response
    function trySignatureAndApkVerification(
        bytes32 msgHash,
        bytes calldata apk,
        bytes calldata apkG2,
        bytes calldata sigma
    ) external view returns (bool,bool);


    // EVENTS
    /// @notice Emitted when `operator` registers with the public keys `pubkeyG1` and `pubkeyG2`.
    event NewPubkeyRegistration(address indexed operator, bytes pubkeyG1, bytes pubkeyG2);
}
