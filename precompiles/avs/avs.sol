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
        string memory avsName,
        string memory avsAddress,
        string memory operatorAddress,
        uint64 action,
        string memory avsOwnerAddress,
        string memory assetID
    ) external returns (bool success);
}