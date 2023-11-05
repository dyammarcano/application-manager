package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"github.com/dyammarcano/application-manager/internal/algorithm/compression"
	"github.com/dyammarcano/base58"
	"io"
	"os"
	"path/filepath"
)

const (
	NonceSize    = 12
	GenKeySize   = 12
	KeysFileName = "keys.dat"
)

var keys map[int][]byte

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	keysPath := filepath.Join(home, KeysFileName)
	if _, err := os.Stat(keysPath); err != nil {
		if os.IsNotExist(err) {
			if err := GenerateKeys(keysPath); err != nil {
				panic(err)
			}
		}
	}

	if keys == nil {
		if err := GetKeys(keysPath); err != nil {
			panic(err)
		}
	}
}

// generateKeys generates a N bytes master key
func generateKeys(size int) ([]byte, error) {
	masterKey := make([]byte, size)
	if _, err := rand.Read(masterKey); err != nil {
		return nil, err
	}
	return masterKey, nil
}

func xorBytes(key []byte) ([]byte, error) {
	versionInt := int(binary.BigEndian.Uint16(key[:2]))
	sl := versionInt % len(keys)
	mt := keys[sl]

	matrixKey := make([]byte, len(mt))
	for i := range mt {
		sk := i % len(key)
		matrixKey[i] = mt[i] ^ key[sk]
	}

	return matrixKey, nil
}

// extractKeys extracts the iv, nonce and secret from the master key
func extractKeys(key []byte) ([]byte, []byte, error) {
	xorKey, err := xorBytes(key)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, NonceSize)
	copy(nonce, xorKey[:NonceSize])

	secret := make([]byte, 32)
	copy(secret, xorKey[NonceSize:])

	return secret, nonce, nil
}

// splitResult splits the result into iv, key and cypherText
func splitResult(result []byte) ([]byte, []byte, []byte, error) {
	secret, nonce, err := extractKeys(result[:GenKeySize])
	if err != nil {
		return nil, nil, nil, err
	}
	return secret, nonce, result[GenKeySize:], nil
}

// encrypt encrypts a message using AES-256-GCM
func encrypt(message []byte, raw bool) ([]byte, error) {
	masterKey, err := generateKeys(GenKeySize)
	if err != nil {
		return nil, err
	}

	response := append([]byte{}, masterKey...)
	key, nonce, err := extractKeys(response)
	if err != nil {
		return nil, err
	}

	cc, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(cc)
	if err != nil {
		return nil, err
	}

	response = append(response, gcm.Seal(nil, nonce, message, nil)...)

	if raw {
		return response, err
	}
	return []byte(base58.StdEncoding.Encode(response)), err
}

// decrypt decrypts a message using AES-256-GCM
func decrypt(message []byte, raw bool) ([]byte, error) {
	if !raw {
		decoded, err := base58.StdEncoding.Decode(string(message))
		if err != nil {
			return nil, err
		}
		message = decoded
	}

	key, nonce, cypherText, err := splitResult(message)
	if err != nil {
		return nil, err
	}

	cc, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(cc)
	if err != nil {
		return nil, err
	}

	decrypted, err := gcm.Open(nil, nonce, cypherText, nil)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

// AutoEncryptString decrypts a message using AES-256-GCM
func AutoEncryptString(message string) (string, error) {
	encrypted, err := encrypt([]byte(message), false)
	if err != nil {
		return "", err
	}
	return string(encrypted), nil
}

// AutoEncryptBytes decrypts a message using AES-256-GCM
func AutoEncryptBytes(message []byte) ([]byte, error) {
	encrypted, err := encrypt(message, true)
	if err != nil {
		return nil, err
	}
	return []byte(encrypted), nil
}

// AutoDecryptString decrypts a message using AES-256-GCM
func AutoDecryptString(message string) (string, error) {
	decrypted, err := decrypt([]byte(message), false)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

// AutoDecryptBytes decrypts a message using AES-256-GCM
func AutoDecryptBytes(message []byte) ([]byte, error) {
	decrypted, err := decrypt(message, true)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}

// EncryptPassword encrypts a message using AES-256-GCM and a password
func EncryptPassword(message, password []byte) ([]byte, error) {
	md5Hash := md5.Sum(password)
	password = []byte(hex.EncodeToString(md5Hash[:]))
	cc, err := aes.NewCipher(password)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(cc)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, _ = io.ReadFull(rand.Reader, nonce)

	cypherText := gcm.Seal(nil, nonce, message, nil)

	return append(nonce, cypherText...), nil
}

// DecryptPassword decrypts a message using AES-256-GCM and a password
func DecryptPassword(message, password []byte) ([]byte, error) {
	md5Hash := md5.Sum(password)
	password = []byte(hex.EncodeToString(md5Hash[:]))
	cc, err := aes.NewCipher(password)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(cc)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce, cipheredText := message[:nonceSize], message[nonceSize:]

	decrypted, err := gcm.Open(nil, nonce, cipheredText, nil)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

func GenerateKeys(keysPath string) error {
	keys = make(map[int][]byte)

	for i := 0; i < 1024; i++ {
		data, err := generateKeys(44)
		if err != nil {
			return err
		}
		keys[i] = data
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(keys); err != nil {
		return err
	}

	// compress data
	comp, err := compression.CompressData(buf.Bytes())
	if err != nil {
		return err
	}

	// save data
	if err = SaveKeys(keysPath, comp); err != nil {
		return err
	}

	return nil
}

func SaveKeys(keysPath string, data []byte) error {
	file, err := os.Create(keysPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func GetKeys(keysPath string) error {
	data, err := os.ReadFile(keysPath)
	if err != nil {
		return err
	}

	dec, err := compression.DecompressData(data)
	if err != nil {
		return err
	}

	if err := gob.NewDecoder(bytes.NewReader(dec)).Decode(&keys); err != nil {
		return err
	}

	return nil
}
