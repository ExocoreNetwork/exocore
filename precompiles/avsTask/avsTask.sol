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
    struct Task {
        uint256 numberToBeSquared;
        uint32 taskCreatedBlock;
        // task submitter decides on the criteria for a task to be completed
        // note that this does not mean the task was "correctly" answered (i.e. the number was squared correctly)
        //      this is for the challenge logic to verify
        // task is completed (and contract will accept its TaskResponse) when each quorumNumbers specified here
        // are signed by at least quorumThresholdPercentage of the operators
        // note that we set the quorumThresholdPercentage to be the same for all quorumNumbers, but this could be changed
        bytes quorumNumbers;
        uint32 quorumThresholdPercentage;
    }

    struct TaskResponse {
        // Can be obtained by the operator from the event NewTaskCreated.
        uint32 referenceTaskIndex;
        // This is just the response that the operator has to compute by itself.
        uint256 numberSquared;
    }

// Extra information related to taskResponse, which is filled inside the contract.
// It thus cannot be signed by operators, so we keep it in a separate struct than TaskResponse
// This metadata is needed by the challenger, so we emit it in the TaskResponded event
    struct TaskResponseMetadata {
        uint32 taskResponsedBlock;
        bytes32 hashOfNonSigners;
    }
interface IAVSTask {
/// TRANSACTIONS
/// @dev IAVSTask the oprator, that will change the state in AVSTask module
/// @param numberToBeSquared The Numbers that need to be squared
/// @param quorumThresholdPercentage The Quorum threshold
/// @param quorumNumbers The Quorum numbers
    function createNewTask(
        uint256 numberToBeSquared,
        uint32 quorumThresholdPercentage,
        bytes calldata quorumNumbers
    ) external returns (bool success);

    /// TRANSACTIONS
/// @dev this function responds to existing tasks.
/// @param Task The task of avs already created
/// @param TaskResponse The Task response parameters
    function respondToTask(
        Task calldata task,
        TaskResponse calldata taskResponse
    ) external returns (bool success);

    /// TRANSACTIONS
/// @dev Get the count of the current task
    function taskNumber() external view returns (uint32);

    /// TRANSACTIONS
/// @dev Get the task window block for the current response
    function getTaskResponseWindowBlock() external view returns (uint32);

    /// @dev This event is emitted when a task created.
    event NewTaskCreated(uint32 indexed taskIndex, Task task);

}
