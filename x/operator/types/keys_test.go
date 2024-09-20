package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type KeysTestSuite struct {
	suite.Suite
}

func (suite *KeysTestSuite) SetupTest() {
}

func TestKeysTestSuite(t *testing.T) {
	suite.Run(t, new(KeysTestSuite))
}

func (suite *KeysTestSuite) TestParseKeyForOperatorAndChainIDToConsKey() {
	operator := "exo1rtg0cgw94ep744epyvanc0wdd5kedwql73vlmr"
	operatorAddr, err := sdk.AccAddressFromBech32(operator)
	suite.NoError(err)
	chainIDWithoutRevision := "exocoretestnet_233"
	key := KeyForOperatorAndChainIDToConsKey(operatorAddr, chainIDWithoutRevision)

	parsedAddr, parsedChainID, err := ParseKeyForOperatorAndChainIDToConsKey(key[1:])
	suite.NoError(err)
	suite.Equal(operatorAddr, parsedAddr)
	suite.Equal(chainIDWithoutRevision, parsedChainID)
}

func (suite *KeysTestSuite) TestParsePrevConsKey() {
	operator := "exo1rtg0cgw94ep744epyvanc0wdd5kedwql73vlmr"
	operatorAddr, err := sdk.AccAddressFromBech32(operator)
	suite.NoError(err)
	chainIDWithoutRevision := "exocoretestnet_233"
	key := KeyForChainIDAndOperatorToPrevConsKey(chainIDWithoutRevision, operatorAddr)

	parsedChainID, parsedAddr, err := ParsePrevConsKey(key[1:])
	suite.NoError(err)
	suite.Equal(operatorAddr, parsedAddr)
	suite.Equal(chainIDWithoutRevision, parsedChainID)
}

func (suite *KeysTestSuite) TestParseKeyForOperatorKeyRemoval() {
	operator := "exo1rtg0cgw94ep744epyvanc0wdd5kedwql73vlmr"
	operatorAddr, err := sdk.AccAddressFromBech32(operator)
	suite.NoError(err)
	chainIDWithoutRevision := "exocoretestnet_233"
	key := KeyForOperatorKeyRemovalForChainID(operatorAddr, chainIDWithoutRevision)

	parsedAddr, parsedChainID, err := ParseKeyForOperatorKeyRemoval(key[1:])
	suite.NoError(err)
	suite.Equal(operatorAddr, parsedAddr)
	suite.Equal(chainIDWithoutRevision, parsedChainID)
}
