package merkletree

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	poseidon "github.com/iden3/go-iden3-crypto/goldenposeidon"
)

// Key stores key of the leaf
type Key [32]byte

const (
	HashPoseidonAllZeroes = "0xc71603f33a1144ca7953db0ab48808f4c4055e3364a246c33c18a9786cb0b359"
)

// keyEthAddr is the common code for all the keys related to ethereum addresses.
func keyEthAddr(ethAddr common.Address, leafType leafType, key1Capacity [4]uint64) ([]byte, error) {
	ethAddrBI := new(big.Int).SetBytes(ethAddr.Bytes())
	ethAddrArr := scalar2fea(ethAddrBI)

	key1 := [8]uint64{
		ethAddrArr[0],
		ethAddrArr[1],
		ethAddrArr[2],
		ethAddrArr[3],
		ethAddrArr[4],
		0,
		uint64(leafType),
		0,
	}

	result, err := poseidon.Hash(key1, key1Capacity)
	if err != nil {
		return nil, err
	}

	return h4ToFilledByteSlice(result[:]), nil
}

func defaultCapIn() ([4]uint64, error) {
	capIn, err := stringToh4(HashPoseidonAllZeroes)
	if err != nil {
		return [4]uint64{}, err
	}

	return [4]uint64{capIn[0], capIn[1], capIn[2], capIn[3]}, nil
}

// KeyEthAddrBalance returns the key of balance leaf:
// hk0: H([0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0])
// key: H([ethAddr[0:4], ethAddr[4:8], ethAddr[8:12], ethAddr[12:16], ethAddr[16:20], 0, 0, 0], [hk0[0], hk0[1], hk0[2], hk0[3]])
func KeyEthAddrBalance(ethAddr common.Address) ([]byte, error) {
	capIn, err := defaultCapIn()
	if err != nil {
		return nil, err
	}

	return keyEthAddr(ethAddr, leafTypeBalance, capIn)
}

// KeyEthAddrNonce returns the key of nonce leaf:
// hk0: H([0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0])
// key: H([ethAddr[0:4], ethAddr[4:8], ethAddr[8:12], ethAddr[12:16], ethAddr[16:20], 0, 1, 0], [hk0[0], hk0[1], hk0[2], hk0[3]]
func KeyEthAddrNonce(ethAddr common.Address) ([]byte, error) {
	capIn, err := defaultCapIn()
	if err != nil {
		return nil, err
	}

	return keyEthAddr(ethAddr, leafTypeNonce, capIn)
}

// KeyContractCode returns the key of contract code leaf:
// hk0: H([0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0])
// key: H([ethAddr[0:4], ethAddr[4:8], ethAddr[8:12], ethAddr[12:16], ethAddr[16:20], 0, 2, 0], [hk0[0], hk0[1], hk0[2], hk0[3]]
func KeyContractCode(ethAddr common.Address) ([]byte, error) {
	capIn, err := defaultCapIn()
	if err != nil {
		return nil, err
	}

	return keyEthAddr(ethAddr, leafTypeCode, capIn)
}

// KeyContractStorage returns the key of contract storage position leaf:
// hk0: H([stoPos[0:4], stoPos[4:8], stoPos[8:12], stoPos[12:16], stoPos[16:20], stoPos[20:24], stoPos[24:28], stoPos[28:32], [0, 0, 0, 0])
// key: H([ethAddr[0:4], ethAddr[4:8], ethAddr[8:12], ethAddr[12:16], ethAddr[16:20], 0, 3, 0], [hk0[0], hk0[1], hk0[2], hk0[3])
func KeyContractStorage(ethAddr common.Address, storagePos []byte) ([]byte, error) {
	storageBI := new(big.Int).SetBytes(storagePos)

	storageArr := scalar2fea(storageBI)

	hk0, err := poseidon.Hash([8]uint64{
		storageArr[0],
		storageArr[1],
		storageArr[2],
		storageArr[3],
		storageArr[4],
		storageArr[5],
		storageArr[6],
		storageArr[7],
	}, [4]uint64{})
	if err != nil {
		return nil, err
	}

	return keyEthAddr(ethAddr, leafTypeStorage, hk0)
}