pragma solidity >=0.8.17;

/// @dev The CLIENT_CHAINS contract's address.
address constant CLIENT_CHAINS_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000000801;

/// @dev The CLIENT_CHAINS contract's instance.
IClientChains constant CLIENT_CHAINS_CONTRACT = IClientChains(
    CLIENT_CHAINS_PRECOMPILE_ADDRESS
);

/// @author Exocore Team
/// @title Client Chains Precompile Contract
/// @dev The interface through which solidity contracts will interact with ClientChains
/// @custom:address 0x0000000000000000000000000000000000000801
interface IClientChains {
    /// @dev Returns the chain indices of the client chains.
    function getClientChains() external view returns (bool, uint16[] memory);
}

