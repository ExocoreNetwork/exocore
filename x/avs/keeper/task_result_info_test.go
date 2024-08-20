package keeper

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/prysmaticlabs/prysm/v4/crypto/bls/blst"
	blscommon "github.com/prysmaticlabs/prysm/v4/crypto/bls/common"
	"golang.org/x/crypto/sha3"
	"math/big"
	"testing"
)

type TaskResponse struct {
	TaskId    uint64
	NumberSum *big.Int
}

func Test_bls_sig(t *testing.T) {

	privateKeys := make([]blscommon.SecretKey, 3)
	for i := 0; i < 3; i++ {
		privateKeys[i], _ = blst.RandKey()
	}

	publicKeys := make([]blscommon.PublicKey, 3)
	for i := 0; i < 3; i++ {
		publicKeys[i] = privateKeys[i].PublicKey()
	}

	taskres := TaskResponse{TaskId: 1, NumberSum: big.NewInt(100)}

	msg, _ := GetTaskResponseDigest(taskres)
	msgBytes := msg[:]

	signatures := make([]blscommon.Signature, 3)
	for i := 0; i < 3; i++ {
		signatures[i] = privateKeys[i].Sign(msgBytes)
	}

	aggsignature := blst.AggregateSignatures(signatures)

	valid2 := aggsignature.FastAggregateVerify(publicKeys, msg)
	valid3 := aggsignature.Eth2FastAggregateVerify(publicKeys, msg)
	fmt.Println("Aggregate signature2 is valid for all messages:", valid2)
	fmt.Println("Aggregate signature3 is valid for all messages:", valid3)

	sigN := privateKeys[1].Sign(msgBytes)
	a := sigN.Marshal()
	fmt.Println(a)
	valid := sigN.Verify(publicKeys[1], msgBytes)
	fmt.Println(" sigN is valid for all messages:", valid)

	b, _ := blst.VerifySignature(a, msg, publicKeys[1])
	fmt.Println(" b is valid for all messages:", b)

}

// GetTaskResponseDigest returns the hash of the TaskResponse, which is what operators sign over
func GetTaskResponseDigest(h TaskResponse) ([32]byte, error) {

	jsonData, err := json.Marshal(h)
	if err != nil {
		fmt.Println("Error marshalling struct to JSON:", err)
		return [32]byte{}, err
	}

	fmt.Println(jsonData)
	fmt.Println(string(jsonData))

	var newPerson TaskResponse
	err = json.Unmarshal(jsonData, &newPerson)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return [32]byte{}, err
	}
	fmt.Println(newPerson)

	var taskResponseDigest [32]byte
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(jsonData)
	copy(taskResponseDigest[:], hasher.Sum(nil)[:32])

	return taskResponseDigest, nil
}

func Test_hash(t *testing.T) {
	taskres := TaskResponse{TaskId: 1, NumberSum: big.NewInt(100)}
	jsonData, _ := json.Marshal(taskres)

	var taskResponseDigest [32]byte
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(jsonData)
	copy(taskResponseDigest[:], hasher.Sum(nil)[:32])
	fmt.Println(hex.EncodeToString(taskResponseDigest[:]))
}
