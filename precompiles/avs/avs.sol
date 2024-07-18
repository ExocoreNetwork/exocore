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
    function registerAVS(
        string[] memory avsOwnerAddress,
        string memory avsName,
        string memory taskAddr,
        string memory rewardContractAddr,
        string memory slashContractAddr,
        string[] memory assetID,
        uint64 minSelfDelegation,
        uint64 unbondingPeriod,
        string memory epochIdentifier
    ) external returns (bool success);

    function updateAVS(
        string[] memory avsOwnerAddress,
        string memory avsName,
        string memory taskAddr,
        string memory rewardContractAddr,
        string memory slashContractAddr,
        string[] memory assetID,
        uint64 minSelfDelegation,
        uint64 unbondingPeriod,
        string memory epochIdentifier
    ) external returns (bool success);


    function deregisterAVS(
        string memory avsName
    ) external returns (bool success);

    function registerOperatorToAVS(
    ) external returns (bool success);

    function deregisterOperatorFromAVS(
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
        bytes calldata pubKey,
        bytes calldata pubkeyRegistrationSignature，
        bytes calldata pubkeyRegistrationMessageHash
    ) external returns (bool success);

    /// @dev Returns the pubkey and pubkey hash of an operator,Reverts if the operator has not registered a valid pubkey
    /// @param operator is the operator for whom the key is being registered
    function getRegisteredPubkey(string memory operator) external returns (bytes memory pubkey);


    /// @dev RegisterBLSPublicKey Emitted when `operator` registers with the public keys `pubKey`.
    /// @param operator the address of the delegator
    /// @param pubKey the address of the validator
    event RegisterBLSPublicKey(
        string indexed operator,
        bytes calldata pubKey
    );

    /// @dev RegisterAVS Emitted when `avs` register to exocore.
    event RegisterAVS(
        string indexed avsAddress,
        string[]  avsOwnerAddress,
        string  avsName,
        string   taskAddr,
        string  rewardContractAddr,
        string  slashContractAddr,
        string[]  assetID,
        uint64 minSelfDelegation,
        uint64 unbondingPeriod,
        string  epochIdentifier
    );

    /// @dev UpdateAVS Emitted when `avs` update to exocore.
    event UpdateAVS(
        string indexed avsAddress,
        string[]  avsOwnerAddress,
        string  avsName,
        string   taskAddr,
        string  rewardContractAddr,
        string  slashContractAddr,
        string[]  assetID,
        uint64 minSelfDelegation,
        uint64 unbondingPeriod,
        string  epochIdentifier
    );

    /// @dev DeregisterAVS Emitted when `avs` Deregister to exocore.
    event DeregisterAVS(
        string indexed avsAddress,
        string  avsName
    );

    /// @dev RegisterOperatorToAVS Emitted when `operator` opt-in to avs.
    event RegisterOperatorToAVS(
        string indexed operator,
        string  avsAddress
    );

    /// @dev DeregisterOperatorFromAVS Emitted when `operator` opt-out to avs.
    event DeregisterOperatorFromAVS(
        string indexed operator,
        string  avsAddress
    );


    /// @dev RegisterAVSTask the oprator, that will change the state in AVSTask module
    /// @param TaskContractAddress avstask Contract Address
    /// @param Name avstask name
    /// @param MetaInfo avstask desc
    event RegisterAVSTask(
        string indexed TaskContractAddress,
        string  Name,
        string  MetaInfo
    );
}


