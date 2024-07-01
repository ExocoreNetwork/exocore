//go:build skiptest

// It's used to facilitate the precompile test without depending on the client chain and Exocore
// gateway contract. This test depends on a running local node, so set the ignore flag to skip the
//  tests. You can execute this test by clicking in an IDE or by using the `go test` command
// These tests can be referred to when implementing automated integration testing tools, which can reduce
// the test workload for the basic precompile functions when the code has some significant changed.

package testutil_test

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ExocoreNetwork/exocore/precompiles/assets"
	"github.com/ExocoreNetwork/exocore/precompiles/delegation"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"math/big"
	"strings"
	"testing"
	"time"
)

type BasicInfoToSendTx struct {
	ctx      context.Context
	sk       *ecdsa.PrivateKey
	caller   common2.Address
	signer   types.Signer
	ethC     *ethclient.Client
	nonce    uint64
	GasPrice *big.Int
	GasLimit uint64
}

var ExocorePrecompileGaLimit = uint64(500000)

func PaddingAddressTo32(address common2.Address) []byte {
	paddingLen := 32 - len(address)
	ret := make([]byte, len(address))
	copy(ret, address[:])
	for i := 0; i < paddingLen; i++ {
		ret = append(ret, 0)
	}
	fmt.Println("the ret is:", hexutil.Encode(ret))
	return ret
}

func SendTxPreparation(t *testing.T) BasicInfoToSendTx {
	// this is the private key of local funded address
	privateKey := "D196DCA836F8AC2FFF45B3C9F0113825CCBB33FA1B39737B948503B263ED75AE"
	sk, err := crypto.HexToECDSA(privateKey)
	assert.NoError(t, err)
	callAddr := crypto.PubkeyToAddress(sk.PublicKey)

	ctx := context.Background()
	exocoreNode := "http://127.0.0.1:8545"
	c, err := rpc.DialContext(ctx, exocoreNode)
	assert.NoError(t, err)
	ethC := ethclient.NewClient(c)
	chainID, err := ethC.ChainID(ctx)
	assert.NoError(t, err)
	signer := types.LatestSignerForChainID(chainID)

	gasLimit := ExocorePrecompileGaLimit
	nonce, err := ethC.NonceAt(ctx, callAddr, nil)
	assert.NoError(t, err)
	gasPrice, err := ethC.SuggestGasPrice(ctx)
	assert.NoError(t, err)
	return BasicInfoToSendTx{
		ctx:      ctx,
		sk:       sk,
		caller:   callAddr,
		ethC:     ethC,
		signer:   signer,
		nonce:    nonce,
		GasLimit: gasLimit,
		GasPrice: gasPrice,
	}
}
func SignAndSendTx(basicInfo *BasicInfoToSendTx, contractAddr *common2.Address, data []byte, t *testing.T) {
	retTx := types.NewTx(&types.LegacyTx{
		Nonce:    basicInfo.nonce,
		To:       contractAddr,
		Value:    big.NewInt(0),
		Gas:      basicInfo.GasLimit,
		GasPrice: basicInfo.GasPrice,
		Data:     data,
	})

	signTx, err := types.SignTx(retTx, basicInfo.signer, basicInfo.sk)
	assert.NoError(t, err)
	fmt.Println("the txID is:", signTx.Hash().String())
	msg := ethereum.CallMsg{
		From: basicInfo.caller,
		To:   retTx.To(),
		Data: retTx.Data(),
	}
	_, err = basicInfo.ethC.CallContract(context.Background(), msg, nil)
	assert.NoError(t, err)

	err = basicInfo.ethC.SendTransaction(basicInfo.ctx, signTx)
	assert.NoError(t, err)

	time.Sleep(20 * time.Second)

	receipt, err := basicInfo.ethC.TransactionReceipt(basicInfo.ctx, signTx.Hash())
	assert.NoError(t, err)
	assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
	fmt.Println("the block height is:", receipt.BlockNumber)
}

func Test_ExocoreClient(t *testing.T) {
	basicInfo := SendTxPreparation(t)

	blockNumber, err := basicInfo.ethC.BlockNumber(basicInfo.ctx)
	assert.NoError(t, err)
	fmt.Println("the blockNumber is:", blockNumber)

	balance, err := basicInfo.ethC.BalanceAt(basicInfo.ctx, basicInfo.caller, nil)
	assert.NoError(t, err)
	fmt.Println("the balance is:", balance)

	//assetsAddr := common2.HexToAddress("0x0000000000000000000000000000000000000800")
	contractBytes, err := basicInfo.ethC.CodeAt(basicInfo.ctx, basicInfo.caller, nil)
	assert.NoError(t, err)
	fmt.Println("the contractBytes is:", hexutil.Encode(contractBytes))
}

func Test_CheckTxStatus(t *testing.T) {
	basicInfo := SendTxPreparation(t)
	txID := common2.HexToHash("0x19d44adc35607e0187300bfedca35a3da11e7684a553f7837294054c3aa3d147")
	tx, isPending, err := basicInfo.ethC.TransactionByHash(basicInfo.ctx, txID)
	assert.NoError(t, err)
	fmt.Println("isPending", isPending)
	fmt.Println("the tx is:", tx.Hash())

	receipt, err := basicInfo.ethC.TransactionReceipt(basicInfo.ctx, txID)
	assert.NoError(t, err)
	assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
	fmt.Println("the block height is:", receipt.BlockNumber)
}

func Test_RegisterOrUpdateClientChain(t *testing.T) {
	basicInfo := SendTxPreparation(t)
	contractAddr := common2.HexToAddress("0x0000000000000000000000000000000000000804")
	assetsAbi, err := abi.JSON(strings.NewReader(assets.AssetsABI))
	assert.NoError(t, err)

	clientChainID := uint32(111)
	addressLength := uint8(20)
	name := "testClientChain"
	metaInfo := "it's a test client chain"
	signatureType := "ecdsa"
	data, err := assetsAbi.Pack("registerOrUpdateClientChain", clientChainID, addressLength, name, metaInfo, signatureType)
	assert.NoError(t, err)
	SignAndSendTx(&basicInfo, &contractAddr, data, t)
}

func Test_RegisterOrUpdateTokens(t *testing.T) {
	basicInfo := SendTxPreparation(t)
	contractAddr := common2.HexToAddress("0x0000000000000000000000000000000000000804")
	assetsAbi, err := abi.JSON(strings.NewReader(assets.AssetsABI))
	assert.NoError(t, err)

	clientChainID := uint32(101)
	testTokenAddr0 := common2.HexToAddress("0xb82381a3fbd3fafa77b3a7be693342618240067b")
	decimal := uint8(8)
	tvlLimit := big.NewInt(3000000000000000000)
	name := "WSTETH"
	metaInfo := "Wrapped STETH"

	data, err := assetsAbi.Pack("registerOrUpdateTokens", clientChainID, PaddingAddressTo32(testTokenAddr0), decimal, tvlLimit, name, metaInfo)

	assert.NoError(t, err)
	SignAndSendTx(&basicInfo, &contractAddr, data, t)
}

func Test_Deposit(t *testing.T) {
	basicInfo := SendTxPreparation(t)

	stakerAddr := common2.HexToAddress("0x217F1887cCE09BFFc4194cca5d561Bc447298d24")
	contractAddr := common2.HexToAddress("0x0000000000000000000000000000000000000804")
	assetsAbi, err := abi.JSON(strings.NewReader(assets.AssetsABI))
	assert.NoError(t, err)
	assetAddr := common2.HexToAddress("0xb82381a3fbd3fafa77b3a7be693342618240067b")
	opAmount := big.NewInt(0).Exp(big.NewInt(100), big.NewInt(8), nil)
	fmt.Println("the opAmount is:", opAmount)
	clientChainID := uint32(101)

	data, err := assetsAbi.Pack("depositTo", clientChainID, PaddingAddressTo32(assetAddr), PaddingAddressTo32(stakerAddr), opAmount)
	SignAndSendTx(&basicInfo, &contractAddr, data, t)
}

func Test_Delegate(t *testing.T) {
	basicInfo := SendTxPreparation(t)

	stakerAddr := common2.HexToAddress("0x217F1887cCE09BFFc4194cca5d561Bc447298d24")
	contractAddr := common2.HexToAddress("0x0000000000000000000000000000000000000805")
	delegationAbi, err := abi.JSON(strings.NewReader(delegation.DelegationABI))
	assetAddr := common2.HexToAddress("0xb82381a3fbd3fafa77b3a7be693342618240067b")
	opAmount := big.NewInt(0).Exp(big.NewInt(100), big.NewInt(8), nil)
	clientChainID := uint32(101)
	// todo: need to ensure the operator address has been registered
	operatorAddr := "exo18cggcpvwspnd5c6ny8wrqxpffj5zmhklprtnph"

	// using the nonce of caller as the layer zero nonce
	data, err := delegationAbi.Pack("delegateToThroughClientChain", clientChainID, basicInfo.nonce, PaddingAddressTo32(assetAddr), PaddingAddressTo32(stakerAddr), []byte(operatorAddr), opAmount)
	assert.NoError(t, err)
	SignAndSendTx(&basicInfo, &contractAddr, data, t)
}

func Test_AssociateOperatorWithStaker(t *testing.T) {
	basicInfo := SendTxPreparation(t)
	stakerAddr := common2.HexToAddress("0x217F1887cCE09BFFc4194cca5d561Bc447298d24")
	clientChainID := uint32(101)
	contractAddr := common2.HexToAddress("0x0000000000000000000000000000000000000805")
	delegationAbi, err := abi.JSON(strings.NewReader(delegation.DelegationABI))
	// todo: need to ensure the operator address has been registered
	operatorAddr := "exo18cggcpvwspnd5c6ny8wrqxpffj5zmhklprtnph"

	// using the nonce of caller as the layer zero nonce
	data, err := delegationAbi.Pack("associateOperatorWithStaker", clientChainID, PaddingAddressTo32(stakerAddr), []byte(operatorAddr))
	assert.NoError(t, err)
	SignAndSendTx(&basicInfo, &contractAddr, data, t)
}
