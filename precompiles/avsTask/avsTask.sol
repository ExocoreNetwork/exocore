pragma solidity >=0.8.17 .0;

/// @dev The AVSTask contract's address.
address constant AVSTASK_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000000901;

/// @dev The AVSTask contract's instance.
IAVSTask constant AVSTASK_CONTRACT = IAVSTask(
    AVSTASK_PRECOMPILE_ADDRESS
);

/// @author Exocore Team
/// @title AVSTask Precompile Contract
/// @dev The interface through which solidity contracts will interact with AVSTask
/// @custom:address 0x0000000000000000000000000000000000000901

interface IAVSTask {
    /// @dev IAVSTask the oprator, that will change the state in AVSTask module
    /// @param TaskContractAddress avstask Contract Address
    /// @param Name avstask name
    /// @param MetaInfo avstask desc
    function registerAVSTask(
        string memory TaskContractAddress,
        string memory Name,
        string memory MetaInfo
    ) external returns (bool success);

    /// @dev Called by the avs manager service register an operator as the owner of a BLS public key.
    /// @param operator is the operator for whom the key is being registered
    /// @param pubKey  the public keys of the operator
    function registerBLSPublicKey(
        string memory operator,
        bytes calldata pubKey
    ) external returns (bool success);

    /// @dev Returns the pubkey and pubkey hash of an operator,Reverts if the operator has not registered a valid pubkey
    /// @param operator is the operator for whom the key is being registered
    function getRegisteredPubkey(address operator) external returns (bytes32);

    // EVENTS
    /// @notice Emitted when `operator` registers with the public keys `pubkeyG1`.
    event NewPubkeyRegistration(address indexed operator, bytes pubkeyG1);
}
