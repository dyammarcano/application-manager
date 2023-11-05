package crypto

import (
	"github.com/dyammarcano/application-manager/internal/mock"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestEncryptString(t *testing.T) {
	encrypted, err := AutoEncryptString(mock.Message5kChars)
	assert.Nil(t, err)

	decrypted, err := AutoDecryptString(encrypted)
	assert.Nil(t, err)

	assert.Equal(t, decrypted, mock.Message5kChars)
}

func TestEncryptBytes(t *testing.T) {
	encrypted, err := AutoEncryptBytes([]byte(mock.Message5kChars))
	assert.Nil(t, err)

	decrypted, err := AutoDecryptBytes(encrypted)
	assert.Nil(t, err)

	assert.Equal(t, decrypted, []byte(mock.Message5kChars))
}
