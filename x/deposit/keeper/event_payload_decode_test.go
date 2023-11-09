package keeper_test

import (
	"bytes"
	"cosmossdk.io/math"
	"encoding/binary"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"log"
	"math/big"
	"testing"
)

func Test_DepositEventPayloadDecode(t *testing.T) {
	payload := "0x000000000000000000000000000000000000000001000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000640000000000000000000000000000000000000000000000"

	log.Println("the payload original data is:", payload)
	payloadBytes, err := hexutil.Decode(payload)
	assert.NoError(t, err)
	//assert.Equal(t, 73, len(payloadBytes))
	log.Println("the payloadBytes length is:", len(payloadBytes))

	action := payloadBytes[0]
	log.Println("the action is:", action)

	tokenAddress := hexutil.Encode(payloadBytes[1:21])
	toAddress := hexutil.Encode(payloadBytes[21:41])
	amount, ok := types.NewIntFromString(hexutil.Encode(payloadBytes[41:73]))
	assert.True(t, ok)

	log.Println("the tokenAddress toAddress and amount is:", tokenAddress, toAddress, amount)
}

func Test_DecodeThroughBigEndian(t *testing.T) {
	payload := "0x000000000000000000000000000000000000000001000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000640000000000000000000000000000000000000000000000"
	payloadBytes, err := hexutil.Decode(payload)
	assert.NoError(t, err)
	r := bytes.NewReader(payloadBytes[0:1])
	var action uint8
	err = binary.Read(r, binary.BigEndian, &action)
	assert.NoError(t, err)

	log.Println("the action is:", action)

	r = bytes.NewReader(payloadBytes[1:21])
	tokenAddress := make([]byte, 0)
	err = binary.Read(r, binary.BigEndian, tokenAddress)
	assert.NoError(t, err)
	log.Println("the tokenAddress is:", tokenAddress)

	r = bytes.NewReader(payloadBytes[21:41])
	var toAddress common.Address
	err = binary.Read(r, binary.BigEndian, &toAddress)
	assert.NoError(t, err)
	log.Println("the toAddress is:", toAddress)

	amount := math.NewIntFromBigInt(big.NewInt(0).SetBytes(payloadBytes[41:73]))
	log.Println("the amount is:", amount)
}
