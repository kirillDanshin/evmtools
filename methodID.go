package evmtools

import "github.com/ethereum/go-ethereum/crypto"

func MethodID(signature string) []byte {
	hash := crypto.Keccak256([]byte(signature))
	return hash[:4]
}
