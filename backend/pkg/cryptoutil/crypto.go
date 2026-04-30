package cryptoutil

import (
	"crypto/sha256"
	"encoding/hex"
)

func SHA256Hex(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
