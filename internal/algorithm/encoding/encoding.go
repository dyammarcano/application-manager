package encoding

import (
	"github.com/dyammarcano/application-manager/internal/algorithm/compression"
	"github.com/dyammarcano/application-manager/internal/algorithm/crypto"
	"github.com/dyammarcano/base58"
)

func Serialize(message string) (string, error) {
	comp, err := compression.CompressData([]byte(message))
	if err != nil {
		return "", err
	}

	enc, err := crypto.AutoEncryptBytes(comp)
	if err != nil {
		return "", err
	}

	return base58.StdEncoding.EncodeToString(enc), nil
}

func Deserialize(message string) (string, error) {
	dec, err := base58.StdEncoding.DecodeString(message)
	if err != nil {
		return "", err
	}

	dec, err = crypto.AutoDecryptBytes(dec)
	if err != nil {
		return "", err
	}

	dec, err = compression.DecompressData(dec)
	if err != nil {
		return "", err
	}

	return string(dec), nil
}
