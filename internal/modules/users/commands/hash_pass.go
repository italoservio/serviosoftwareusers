package commands

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

func HashPass(secret, pass string) (string, error) {
	if secret == "" {
		return "", errors.New("Nao ha chave de criptografia de senha")
	}

	hasher := sha256.New()
	hasher.Write([]byte(pass))

	return hex.EncodeToString(hasher.Sum([]byte(secret))), nil
}
