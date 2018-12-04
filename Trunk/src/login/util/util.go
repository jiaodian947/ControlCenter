package util

import (
	"encoding/hex"

	"github.com/google/uuid"
)

const primeRK = 16777619

// hashStr returns the hash and the appropriate multiplicative
// factor for use in Rabin-Karp algorithm.
func Hash(sep string) uint32 {
	hash := uint32(0)
	for i := 0; i < len(sep); i++ {
		//hash = 0 * 16777619 + sep[i]
		hash = hash*primeRK + uint32(sep[i])
	}
	return hash
}

func UUID() string {
	uid := uuid.New()
	var buf [32]byte
	hex.Encode(buf[:], uid[:])
	return string(buf[:])
}
