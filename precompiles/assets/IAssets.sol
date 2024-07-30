// SPDX-License-Identifier: MIT
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
    function depositTo(uint32 clientChainID, bytes memory assetsAddress, bytes memory stakerAddress, uint256 opAmount)
        external
        returns (bool success, uint256 latestAssetState);

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

    /// @dev registers or updates a client chain to allow deposits / staking, etc.
    /// from that chain.
    /// @param clientChainID is the layerZero chainID if it is supported.
    //  It might be allocated by Exocore when the client chain isn't supported
    //  by layerZero
    function registerOrUpdateClientChain(
        uint32 clientChainID,
        uint8 addressLength,
        string calldata name,
        string calldata metaInfo,
        string calldata signatureType
    ) external returns (bool success, bool updated);

    /// @dev register or update token addresses to exocore
    /// @dev note that there is no way to delete a token. If a token is to be removed,
    /// the TVL limit should be set to 0.
    /// @param clientChainID is the identifier of the token's home chain (LZ or otherwise)
    /// @param token is the address of the token on the home chain
    /// @param decimals is the number of decimals of the token
    /// @param tvlLimit is the number of tokens that can be deposited in the system. Set to
    /// maxSupply if there is no limit
    /// @param name is the name of the token
    /// @param metaData is the arbitrary metadata of the token
    /// @return success if the token registration is successful
    /// @return updated whether the token was added or updated
    function registerOrUpdateTokens(
        uint32 clientChainID,
        bytes calldata token,
        uint8 decimals,
        uint256 tvlLimit,
        string calldata name,
        string calldata metaData
    ) external returns (bool success, bool updated);

    /// QUERIES
    /// @dev Returns the chain indices of the client chains.
    function getClientChains() external view returns (bool, uint32[] memory);

}
