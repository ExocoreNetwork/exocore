package types

import (
	sdkmath "cosmossdk.io/math"
	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/ExocoreNetwork/exocore/x/delegation/keeper"
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ OracleKeeper = MockOracle{}
	_ AVSKeeper    = MockAVS{}
)

type AssetsKeeper interface {
	GetStakingAssetInfo(
		ctx sdk.Context, assetID string,
	) (info *assetstype.StakingAssetInfo, err error)
	GetAssetsDecimal(
		ctx sdk.Context, assets map[string]interface{},
	) (decimals map[string]uint32, err error)
	IteratorAssetsForOperator(
		ctx sdk.Context, operator string, assetsFilter map[string]interface{},
		f func(assetID string, state *assetstype.OperatorAssetInfo) error,
	) error
	AppChainInfoIsExist(ctx sdk.Context, chainID string) bool
	GetOperatorAssetInfos(
		ctx sdk.Context, operatorAddr sdk.Address, assetsFilter map[string]interface{},
	) (assetsInfo map[string]*assetstype.OperatorAssetInfo, err error)
	GetOperatorSpecifiedAssetInfo(ctx sdk.Context, operatorAddr sdk.Address, assetID string) (info *assetstype.OperatorAssetInfo, err error)
	UpdateStakerAssetState(
		ctx sdk.Context, stakerID string, assetID string,
		changeAmount assetstype.DeltaStakerSingleAsset,
	) (err error)
	UpdateOperatorAssetState(
		ctx sdk.Context, operatorAddr sdk.Address, assetID string,
		changeAmount assetstype.DeltaOperatorSingleAsset,
	) (err error)
	GetAllStakingAssetsInfo(ctx sdk.Context) (allAssets map[string]*assetstype.StakingAssetInfo, err error)
}

var _ DelegationKeeper = &keeper.Keeper{}

type DelegationKeeper interface {
	GetSingleDelegationInfo(
		ctx sdk.Context, stakerID, assetID, operatorAddr string,
	) (*delegationtype.DelegationAmounts, error)
	DelegationStateByOperatorAssets(
		ctx sdk.Context, operatorAddr string, assetsFilter map[string]interface{},
	) (map[string]map[string]delegationtype.DelegationAmounts, error)
	UpdateDelegationState(
		ctx sdk.Context, stakerID, assetID, opAddr string, deltaAmounts *delegationtype.DeltaDelegationAmounts,
	) (bool, error)
	IterateUndelegationsByOperator(
		ctx sdk.Context, operator string, heightFilter *uint64, isUpdate bool,
		opFunc func(undelegation *delegationtype.UndelegationRecord) error) error
	GetStakersByOperator(
		ctx sdk.Context, operator, assetID string,
	) (delegationtype.StakerMap, error)
	SetStakerShareToZero(
		ctx sdk.Context, operator, assetID string, stakerMap delegationtype.StakerMap,
	) error
	DeleteStakerMapForOperator(ctx sdk.Context, operator, assetID string) error
	GetStakerUndelegationRecords(
		ctx sdk.Context, stakerID, assetID string,
	) (records []*delegationtype.UndelegationRecord, err error)
	SetSingleUndelegationRecord(
		ctx sdk.Context, record *delegationtype.UndelegationRecord,
	) (recordKey []byte, err error)
	CalculateSlashShare(
		ctx sdk.Context, operator sdk.AccAddress, stakerID, assetID string, slashAmount sdkmath.Int,
	) (share sdkmath.LegacyDec, err error)
	RemoveShare(
		ctx sdk.Context, isUndelegation bool, operator sdk.AccAddress,
		stakerID, assetID string, share sdkmath.LegacyDec,
	) (removeToken sdkmath.Int, err error)
}

type PriceChange struct {
	OriginalPrice sdkmath.Int
	NewPrice      sdkmath.Int
	Decimal       uint8
}

// Price represents the expected return type from the price Oracle
// the first field is the value and the second is the decimal
// it's same as the price type of ChainLink
type Price struct {
	Value   sdkmath.Int
	Decimal uint8
}

// OracleKeeper is the oracle interface expected by operator module
// These functions need to be implemented by the oracle module
type OracleKeeper interface {
	// GetSpecifiedAssetsPrice is a function to retrieve the asset price according to the
	// assetID.
	GetSpecifiedAssetsPrice(ctx sdk.Context, assetID string) (Price, error)
	// GetMultipleAssetsPrices is a function to retrieve multiple assets prices according to the
	// assetID.
	GetMultipleAssetsPrices(ctx sdk.Context, assets map[string]interface{}) (map[string]Price, error)
	// GetPriceChangeAssets the operator module expect a function that can retrieve all
	// information about assets price change. Then it can update the USD share state according
	// to the change information. This function need to return a map, the key is assetID and the
	// value is PriceChange
	GetPriceChangeAssets(ctx sdk.Context) (map[string]*PriceChange, error)
}

type MockOracle struct{}

func (MockOracle) GetSpecifiedAssetsPrice(_ sdk.Context, _ string) (Price, error) {
	return Price{
		Value:   sdkmath.NewInt(1),
		Decimal: 0,
	}, nil
}

func (MockOracle) GetPriceChangeAssets(_ sdk.Context) (map[string]*PriceChange, error) {
	// use USDT as the mock asset
	ret := make(map[string]*PriceChange, 0)
	usdtAssetID := "0xdac17f958d2ee523a2206206994597c13d831ec7_0x65"
	ret[usdtAssetID] = &PriceChange{
		NewPrice:      sdkmath.NewInt(1),
		OriginalPrice: sdkmath.NewInt(1),
		Decimal:       0,
	}
	return nil, nil
}

func (MockOracle) GetMultipleAssetsPrices(_ sdk.Context, assets map[string]interface{}) (map[string]Price, error) {
	ret := make(map[string]Price, 0)
	for assetID := range assets {
		ret[assetID] = Price{
			Value:   sdkmath.NewInt(1),
			Decimal: 0,
		}
	}
	return ret, nil
}

type MockAVS struct {
	AssetsKeeper AssetsKeeper
}

func (a MockAVS) GetAVSSupportedAssets(ctx sdk.Context, _ string) (map[string]interface{}, error) {
	// set all registered assets as the default asset supported by mock AVS
	ret := make(map[string]interface{})
	allAssets, err := a.AssetsKeeper.GetAllStakingAssetsInfo(ctx)
	if err != nil {
		return nil, err
	}
	for assetID := range allAssets {
		ret[assetID] = nil
	}
	return ret, nil
}

func (a MockAVS) GetAVSSlashContract(_ sdk.Context, _ string) (string, error) {
	return "", nil
}

func (a MockAVS) GetAVSAddrByChainID(_ sdk.Context, chainID string) (string, error) {
	return chainID, nil
}

func (a MockAVS) GetAVSMinimumSelfDelegation(_ sdk.Context, _ string) (sdkmath.LegacyDec, error) {
	return sdkmath.LegacyNewDec(0), nil
}

func (a MockAVS) GetEpochEndAVSs(ctx sdk.Context) ([]string, error) {
	avsList := make([]string, 0)
	avsList = append(avsList, ctx.ChainID())
	return avsList, nil
}

func (a MockAVS) GetHeightForVotingPower(_ sdk.Context, _ string, height int64) (int64, error) {
	return height, nil
}

type AVSKeeper interface {
	// GetAVSSupportedAssets The ctx can be historical or current, depending on the state you
	// wish to retrieve. If the caller want to retrieve a historical assets info supported by
	// Avs, it needs to generate a historical context through calling
	// `ContextForHistoricalState` implemented in x/assets/types/general.go
	GetAVSSupportedAssets(ctx sdk.Context, avsAddr string) (map[string]interface{}, error)
	GetAVSSlashContract(ctx sdk.Context, avsAddr string) (string, error)
	// GetAVSAddrByChainID get the general Avs address for dogfood module.
	GetAVSAddrByChainID(ctx sdk.Context, chainID string) (string, error)
	// GetAVSMinimumSelfDelegation returns the USD value of minimum self delegation, which
	// is set for operator
	GetAVSMinimumSelfDelegation(ctx sdk.Context, avsAddr string) (sdkmath.LegacyDec, error)
	// GetEpochEndAVSs returns the AVS list where the current block marks the end of their epoch.
	// todo: maybe the epoch of different AVSs should be implemented in the AVS module,then
	// the other modules implement the EpochsHooks to trigger state updating.
	GetEpochEndAVSs(ctx sdk.Context) ([]string, error)
	// GetHeightForVotingPower retrieves the height of the last block in the epoch
	// where the voting power used at the current height resides
	GetHeightForVotingPower(ctx sdk.Context, avsAddr string, height int64) (int64, error)
}

// add for dogfood

type SlashKeeper interface {
	IsOperatorFrozen(ctx sdk.Context, addr sdk.AccAddress) bool
}

type OperatorHooks interface {
	// This hook is called when an operator opts in to a chain.
	AfterOperatorOptIn(
		ctx sdk.Context, addr sdk.AccAddress, chainID string,
		pubKey *tmprotocrypto.PublicKey,
	)
	// This hook is called when an operator's consensus key is replaced for
	// a chain.
	AfterOperatorKeyReplacement(
		ctx sdk.Context, addr sdk.AccAddress, oldKey *tmprotocrypto.PublicKey,
		newKey *tmprotocrypto.PublicKey, chainID string,
	)
	// This hook is called when an operator opts out of a chain.
	AfterOperatorOptOutInitiated(
		ctx sdk.Context, addr sdk.AccAddress, chainID string, key *tmprotocrypto.PublicKey,
	)
}
