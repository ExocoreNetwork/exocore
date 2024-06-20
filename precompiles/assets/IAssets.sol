pragma solidity >=0.8.17;

/// @dev The Assets contract's address.
address constant ASSETS_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000000804;

/// @dev The Assets contract's instance.
IAssets constant ASSETS_CONTRACT = IAssets(ASSETS_PRECOMPILE_ADDRESS);

/// @author Exocore Team
/// @title Assets Precompile Contract
/// @dev The interface through which solidity contracts will interact with assets module
/// @custom:address 0x0000000000000000000000000000000000000804
interface IAssets {

    /// TRANSACTIONS
    /// @dev deposit the client chain assets for the staker,
    /// that will change the state in deposit module
    /// Note that this address cannot be a module account.
    /// @param clientChainID is the layerZero chainID if it is supported.
    //  It might be allocated by Exocore when the client chain isn't supported
    //  by layerZero
    /// @param assetsAddress The client chain asset address
    /// @param stakerAddress The staker address
    /// @param opAmount The amount to deposit
    function depositTo(
        uint32 clientChainID,
        bytes memory assetsAddress,
        bytes memory stakerAddress,
        uint256 opAmount) external
    returns (bool success, uint256 latestAssetState);

    /// TRANSACTIONS
    /// @dev withdraw To the staker, that will change the state in withdraw module
    /// Note that this address cannot be a module account.
    /// @param clientChainID is the layerZero chainID if it is supported.
    //  It might be allocated by Exocore when the client chain isn't supported
    //  by layerZero
    /// @param assetsAddress The client chain asset Address
    /// @param withdrawAddress The withdraw address
    /// @param opAmount The withdraw amount
    function withdrawPrincipal(
        uint32 clientChainID,
        bytes memory assetsAddress,
        bytes memory withdrawAddress,
        uint256 opAmount
    ) external returns (bool success, uint256 latestAssetState);

    /// QUERIES
    /// @dev Returns the chain indices of the client chains.
    function getClientChains() external view returns (bool, uint32[] memory);

    /// TRANSACTIONS
    /// @dev register some client chain to allow token registration from that chain, staking
    /// from that chain, and other operations from that chain.
    /// @param clientChainID is the layerZero chainID if it is supported.
    //  It might be allocated by Exocore when the client chain isn't supported
    //  by layerZero
    function registerClientChain(
        uint32 clientChainID,
        uint32 addressLength,
        string memory name,
        string memory metaInfo,
        string memory signatureType
    ) external returns (bool success);

    /// TRANSACTIONS
    /// @dev register unwhitelisted token addresses to exocore
    /// @param clientChainID is the layerZero chainID if it is supported.
    //  It might be allocated by Exocore when the client chain isn't supported
    //  by layerZero
    /// @param tokens The token addresses that would be registered to exocore
    function registerTokens(
        uint32 clientChainID,
        bytes[] memory tokens,
        uint8[] memory decimals,
        uint256[] memory tvlLimit
    ) external returns (bool success);
}