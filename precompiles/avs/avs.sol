pragma solidity >=0.8.17;

/// @dev The avs-manager contract's address.
address constant AVSMANAGER_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000000901;

/// @dev The avs-manager contract's instance.
IAVSManager constant AVSMANAGER_CONTRACT = IAVSManager(
    AVSMANAGER_PRECOMPILE_ADDRESS
);

/// @author Exocore Team
/// @title AVS-Manager Precompile Contract
/// @dev The interface through which solidity contracts will interact with AVS-Manager
/// @custom:address 0x0000000000000000000000000000000000000901
interface IAVSManager {
    function RegisterAVS(
        string[] memory avsOwnerAddress,
        string memory avsName,
        string memory rewardContractAddr,
        string memory slashContractAddr,
        string[] memory assetID,
        uint64 minSelfDelegation,
        uint64 unbondingPeriod,
        string memory epochIdentifier
    ) external returns (bool success);

    function UpdateAVS(
        string[] memory avsOwnerAddress,
        string memory avsName,
        string memory rewardContractAddr,
        string memory slashContractAddr,
        string[] memory assetID,
        uint64 minSelfDelegation,
        uint64 unbondingPeriod,
        string memory epochIdentifier
    ) external returns (bool success);


    function DeregisterAVS(
        string memory avsName
    ) external returns (bool success);

    function RegisterOperatorToAVS(
    ) external returns (bool success);

    function DeregisterOperatorFromAVS(
    ) external returns (bool success);


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
    function getRegisteredPubkey(string memory operator) external returns (bytes memory pubkey);

    // EVENTS
    /// @notice Emitted when `operator` registers with the public keys `pubKey`.
    event NewPubkeyRegistration(string indexed operator, bytes pubKey);

}


