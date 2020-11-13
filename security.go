package interactions

import (
	"crypto/ed25519"
	"encoding/hex"
)

func parsePublicKey(key string) (ed25519.PublicKey, error) {
	return hex.DecodeString(key)
}

func (a *App) verifyMessage(data []byte, signature string) bool {
	sig, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}
	return ed25519.Verify(a.pubKey, data, sig)
}
