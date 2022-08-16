package interactions

import (
	"crypto/ed25519"
	"encoding/hex"
)

func parsePublicKey(key string) (ed25519.PublicKey, error) {
	return hex.DecodeString(key)
}

func verifyMessage(data []byte, signature string, key ed25519.PublicKey) bool {
	sig, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}
	return ed25519.Verify(key, data, sig)
}
