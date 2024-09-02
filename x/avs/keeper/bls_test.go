package keeper_test

import (
	"encoding/hex"
	"fmt"
	"github.com/ExocoreNetwork/exocore/x/avs/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/prysmaticlabs/prysm/v4/crypto/bls/blst"
	blscommon "github.com/prysmaticlabs/prysm/v4/crypto/bls/common"
	"math/big"
)

func (suite *AVSTestSuite) TestOperator_pubkey() {
	operatorAddr := "exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkjr"

	privateKey, err := blst.RandKey()
	publicKey := privateKey.PublicKey()
	blsPub := &types.BlsPubKeyInfo{
		Operator: operatorAddr,
		PubKey:   publicKey.Marshal(),
		Name:     "",
	}
	fmt.Println("pubkey:", hex.EncodeToString(publicKey.Marshal()))
	err = suite.App.AVSManagerKeeper.SetOperatorPubKey(suite.Ctx, blsPub)
	suite.NoError(err)

	pub, err := suite.App.AVSManagerKeeper.GetOperatorPubKey(suite.Ctx, operatorAddr)
	suite.NoError(err)
	suite.Equal(publicKey.Marshal(), pub.PubKey)

	taskRes := types.TaskResponse{TaskID: 1, NumberSum: big.NewInt(100)}

	msg, _ := types.GetTaskResponseDigest(taskRes)
	msgBytes := msg[:]
	sig := privateKey.Sign(msgBytes)
	fmt.Println("sig:", hex.EncodeToString(sig.Marshal()))
	jsonData, err := types.MarshalTaskResponse(taskRes)
	hash := crypto.Keccak256Hash(jsonData)
	fmt.Println("res:", hex.EncodeToString(jsonData))
	fmt.Println("hash:", hash.String())
	sig1, _ := hex.DecodeString("af22f968871395eca62fdb91bc39c2d93569b50678ed73f00c3a6e054512bdc6cb73da7972c9553931aec25bce4973cf15227d2d596492642baaaf2ac1a1a9605b5cf1312fc1e3532aa43a22460e5ce7c081d643dce806f95f26a2df84bdfc66")
	fmt.Println(sig1)

	valid := sig.Verify(publicKey, msgBytes)
	suite.True(valid)

	valid1, _ := blst.VerifySignature(sig.Marshal(), msg, publicKey)
	suite.NoError(err)

	suite.True(valid1)

}

func (suite *AVSTestSuite) Test_hash() {
	taskres := types.TaskResponse{TaskID: 1, NumberSum: big.NewInt(100)}
	jsonData, err := types.MarshalTaskResponse(taskres)
	suite.NoError(err)
	hash := crypto.Keccak256Hash(jsonData)
	taskResponseDigest := hash.Bytes()
	suite.Equal(len(taskResponseDigest), 32)
}

// For deterministic result（msg）, aggregate signatures and verify them
func (suite *AVSTestSuite) Test_bls_agg() {
	taskres := types.TaskResponse{TaskID: 1, NumberSum: big.NewInt(100)}
	jsonData, _ := types.MarshalTaskResponse(taskres)
	msg := crypto.Keccak256Hash(jsonData)
	privateKeys := make([]blscommon.SecretKey, 4)
	publicKeys := make([]blscommon.PublicKey, 4)

	sigs := make([]blscommon.Signature, 4)
	for i := 0; i < 4; i++ {
		privateKey, _ := blst.RandKey()
		privateKeys[i] = privateKey
		publicKeys[i] = privateKey.PublicKey()
		sigs[i] = privateKey.Sign(msg[:])
	}

	aggPublicKey := blst.AggregateMultiplePubkeys(publicKeys)

	aggSignature := blst.AggregateSignatures(sigs)

	valid := aggSignature.Verify(aggPublicKey, msg.Bytes())

	suite.True(valid, "Signature verification failed")

	valid1 := aggSignature.FastAggregateVerify(publicKeys, msg)
	suite.True(valid1, "Signature verification failed")
}

// For uncertain results, i.e. multiple msgs, the test aggregation signature verification fails
func (suite *AVSTestSuite) Test_bls_agg_uncertainMsgs() {
	privateKeys := make([]blscommon.SecretKey, 4)
	publicKeys := make([]blscommon.PublicKey, 4)
	msgs := make([][]byte, 4)
	sigs := make([]blscommon.Signature, 4)
	for i := 0; i < 4; i++ {
		privateKey, _ := blst.RandKey()
		privateKeys[i] = privateKey
		publicKeys[i] = privateKey.PublicKey()
		msgs[i] = []byte{byte(i)}
		sigs[i] = privateKey.Sign([]byte{byte(i)})
	}

	aggPublicKey := blst.AggregateMultiplePubkeys(publicKeys)

	aggSignature := blst.AggregateSignatures(sigs)

	valid := aggSignature.Verify(aggPublicKey, msgs[1])

	suite.False(valid, "Signature verification failed")

	var array32 [32]byte

	// Copy data into the 32-byte array
	copy(array32[:], msgs[1])
	valid1 := aggSignature.FastAggregateVerify(publicKeys, array32)
	suite.False(valid1, "Signature verification failed")

}
