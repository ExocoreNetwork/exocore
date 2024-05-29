pragma solidity >=0.8.17;

/// @dev The avs-manager contract's address.
address constant AVSMANAGER_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000000902;

/// @dev The avs-manager contract's instance.
IAVSManager constant AVSMANAGER_CONTRACT = IAVSManager(
    AVSMANAGER_PRECOMPILE_ADDRESS
);

/// @author Exocore Team
/// @title AVS-Manager Precompile Contract
/// @dev The interface through which solidity contracts will interact with AVS-Manager
/// @custom:address 0x0000000000000000000000000000000000000902
interface IAVSManager {
    function AVSAction(
        string[] memory avsOwnerAddress,
        string memory avsName,
        string memory slashContractAddr,
        string[] memory assetID,
        uint64 action,
        uint64 minimumDelegation,
        uint64 unbondingEpochs
    ) external returns (bool success);


    function OperatorOptAction(
        uint64 action
    ) external returns (bool success);

}


