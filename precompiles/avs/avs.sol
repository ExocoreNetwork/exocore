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
    /// @dev Register AVS contract to EXO.
    /// @param avsName The name of AVS.
    /// @param minStakeAmount The minimum amount of funds staked by each operator.
    /// @param taskAddr The task address of AVS.
    /// @param slashAddr The slash address of AVS.
    /// @param rewardAddr The reward address of AVS.
    /// @param avsOwnerAddress The owners who have permission for AVS.
    /// @param assetIds The basic asset information of AVS.
    /// @param avsUnbondingPeriod The unbonding duration of AVS.
    /// @param minSelfDelegation The minimum delegation amount for an operator.
    /// @param epochIdentifier The AVS epoch identifier.
    /// @param params 1.miniOptInOperators The minimum number of opt-in operators.
    ///2.minTotalStakeAmount The minimum total amount of stake by all operators.
    ///3.avsReward The proportion of reward for AVS.
    ///4.avsSlash The proportion of slash for AVS.
    function registerAVS(
        string memory avsName,
        uint64 minStakeAmount,
        address taskAddr,
        address slashAddr,
        address rewardAddr,
        string[] memory avsOwnerAddress,
        string[] memory assetIds,
        uint64 avsUnbondingPeriod,
        uint64 minSelfDelegation,
        string memory epochIdentifier,
        uint64[] memory params
    ) external returns (bool success);

    /// @dev Update AVS info to EXO.
    /// @param avsName The name of AVS.
    /// @param minStakeAmount The minimum amount of funds staked by each operator.
    /// @param taskAddr The task address of AVS.
    /// @param slashAddr The slash address of AVS.
    /// @param rewardAddr The reward address of AVS.
    /// @param avsOwnerAddress The owners who have permission for AVS.
    /// @param assetIds The basic asset information of AVS.
    /// @param avsUnbondingPeriod The unbonding duration of AVS.
    /// @param minSelfDelegation The minimum delegation amount for an operator.
    /// @param epochIdentifier The AVS epoch identifier.
    /// @param params 1.miniOptInOperators The minimum number of opt-in operators.
    ///2.minTotalStakeAmount The minimum total amount of stake by all operators.
    ///3.avsReward The proportion of reward for AVS.
    ///4.avsSlash The proportion of slash for AVS.
    function updateAVS(
        string memory avsName,
        uint64 minStakeAmount,
        address taskAddr,
        address slashAddr,
        address rewardAddr,
        string[] memory avsOwnerAddress,
        string[] memory assetIds,
        uint64 avsUnbondingPeriod,
        uint64 minSelfDelegation,
        string memory epochIdentifier,
        uint64[] memory params
    ) external returns (bool success);

    /// @dev Deregister avs from exo
    /// @param avsName The name of AVS.
    function deregisterAVS(
        string memory avsName
    ) external returns (bool success);

    /// @dev RegisterOperatorToAVS operator opt in current avs
    function registerOperatorToAVS(
    ) external returns (bool success);

    /// @dev DeregisterOperatorFromAVS operator opt out current avs
    function deregisterOperatorFromAVS(
    ) external returns (bool success);


    /// @dev CreateTask , avs owner create a new task
    /// @param name The name of the task.
    /// @param data The data supplied by the contract, usually ABI-encoded.
    /// @param taskId The task ID of the task.
    /// @param taskResponsePeriod The deadline for task response.
    /// @param taskChallengePeriod The challenge period for the task.
    /// @param thresholdPercentage The signature threshold percentage.
    function createTask(
        string memory name,
        bytes calldata data,
        string memory taskId,
        uint64 taskResponsePeriod,
        uint64 taskChallengePeriod,
        uint64 thresholdPercentage
    ) external returns (bool success);

    /// @dev SubmitProof ,After processing the task contract, aggregate the signature and submit the processed proof
    /// @param taskId The task ID of the task.
    /// @param taskContractAddress The contract address of AVSTask.
    /// @param aggregator The aggregator address.
    /// @param avsAddress The address of AVS.
    /// @param operatorStatus The status and proof of operators.
    function submitProof(
        string memory taskId,
        string memory taskContractAddress,
        string memory aggregator,
        string memory avsAddress,
        bytes calldata operatorStatus
    ) external returns (bool success);


    /// @dev Called by the avs manager service register an operator as the owner of a BLS public key.
    /// @param operator is the operator for whom the key is being registered
    /// @param name the name of public keys
    /// @param pubKey the public keys of the operator
    /// @param pubkeyRegistrationSignature the public keys of the operator
    /// @param pubkeyRegistrationMessageHash the public keys of the operator
    function registerBLSPublicKey(
        string memory operator,
        string calldata name,
        bytes calldata pubKey,
        bytes calldata pubkeyRegistrationSignature,
        bytes calldata pubkeyRegistrationMessageHash
    ) external returns (bool success);



    /// @dev Returns the pubkey and pubkey hash of an operator
    /// @param operator is the operator for whom the key is being registered
    function getRegisteredPubkey(string memory operator) external returns (bytes calldata pubkey);

    /// @dev Returns the operators of all opt-in in the current avs
    /// @param avsAddress avs address
    function getOptInOperators(address avsAddress) external returns (string[] calldata operators);

    /// @dev RegisterBLSPublicKey Emitted when `operator` registers with the public keys `pubKey`.
    /// @param operator the address of the delegator
    /// @param pubKey the address of the validator
    event RegisterBLSPublicKey(
        string indexed operator,
        string name,
        bytes pubKey,
        bytes pubkeyRegistrationSignature,
        bytes pubkeyRegistrationMessageHash
    );

    /// @dev RegisterAVS Emitted when `avs` register to exocore.
    /// @dev Register AVS contract to EXO.
    /// @param avsAddress The address of AVS.
    /// @param avsName The name of AVS.
    /// @param minStakeAmount The minimum amount of funds staked by each operator.
    /// @param taskAddr The task address of AVS.
    /// @param slashAddr The slash address of AVS.
    /// @param rewardAddr The reward address of AVS.
    /// @param avsOwnerAddress The owners who have permission for AVS.
    /// @param assetIds The basic asset information of AVS.
    /// @param avsUnbondingPeriod The unbonding duration of AVS.
    /// @param minSelfDelegation The minimum delegation amount for an operator.
    /// @param epochIdentifier The AVS epoch identifier.
    /// @param params 1.miniOptInOperators The minimum number of opt-in operators.
    ///2.minTotalStakeAmount The minimum total amount of stake by all operators.
    ///3.avsReward The proportion of reward for AVS.
    ///4.avsSlash The proportion of slash for AVS.
    event RegisterAVS(
        string indexed avsAddress,
        string avsName,
        uint64 minStakeAmount,
        address taskAddr,
        address slashAddr,
        address rewardAddr,
        address[] avsOwnerAddress,
        string[] assetIds,
        uint64 avsUnbondingPeriod,
        uint64 minSelfDelegation,
        string epochIdentifier,
        uint64[] params
    );

    /// @dev UpdateAVS Emitted when `avs` update to exocore.
    /// @param avsAddress The address of AVS.
    /// @param avsName The name of AVS.
    /// @param minStakeAmount The minimum amount of funds staked by each operator.
    /// @param taskAddr The task address of AVS.
    /// @param slashAddr The slash address of AVS.
    /// @param rewardAddr The reward address of AVS.
    /// @param avsOwnerAddress The owners who have permission for AVS.
    /// @param assetIds The basic asset information of AVS.
    /// @param avsUnbondingPeriod The unbonding duration of AVS.
    /// @param minSelfDelegation The minimum delegation amount for an operator.
    /// @param epochIdentifier The AVS epoch identifier.
    /// @param params 1.miniOptInOperators The minimum number of opt-in operators.
    ///2.minTotalStakeAmount The minimum total amount of stake by all operators.
    ///3.avsReward The proportion of reward for AVS.
    ///4.avsSlash The proportion of slash for AVS.
    event UpdateAVS(
        string indexed avsAddress,
        string avsName,
        uint64 minStakeAmount,
        address taskAddr,
        address slashAddr,
        address rewardAddr,
        address[] avsOwnerAddress,
        string[] assetIds,
        uint64 avsUnbondingPeriod,
        uint64 minSelfDelegation,
        string epochIdentifier,
        uint64[]  params
    );

    /// @dev DeregisterAVS Emitted when `avs` Deregister to exocore.
    /// @param avsAddress The address of AVS.
    /// @param avsName The name of AVS.
    event DeregisterAVS(
        string indexed avsAddress,
        string  avsName
    );

    /// @dev RegisterOperatorToAVS Emitted when `operator` opt-in to avs.
    /// @param operator address.
    /// @param avsAddress The address of AVS.
    event RegisterOperatorToAVS(
        string indexed operator,
        string indexed avsAddress
    );

    /// @dev DeregisterOperatorFromAVS Emitted when `operator` opt-out to avs.
    /// @param operator address.
    /// @param avsAddress The address of AVS.
    event DeregisterOperatorFromAVS(
        string indexed operator,
        string indexed avsAddress
    );

    /// @dev CreateTask Emitted when `avs` CreateTask.
    /// @param taskContractAddress The contract address of AVSTask.
    /// @param taskId The task ID of the task.
    /// @param name The name of the task.
    /// @param data The data supplied by the contract, usually ABI-encoded.
    /// @param taskResponsePeriod The deadline for task response.
    /// @param taskChallengePeriod The challenge period for the task.
    /// @param thresholdPercentage The signature threshold percentage.
    event CreateTask(
        string indexed taskContractAddress,
        string indexed taskId,
        string name,
        bytes data,
        uint64 taskResponsePeriod,
        uint64 taskChallengePeriod,
        uint64 thresholdPercentage
    );
    /// @dev SubmitProof Emitted when task contract submit proof.
    /// @param taskContractAddress The contract address of AVSTask.
    /// @param taskId The task ID of the task.
    /// @param aggregator The aggregator address.
    /// @param avsAddress The address of AVS.
    /// @param operatorStatuses The status and proof of operators.
    event SubmitProof(
        string indexed taskContractAddress,
        string indexed taskId,
        string aggregator,
        string avsAddress,
        bytes operatorStatuses
    );
}
