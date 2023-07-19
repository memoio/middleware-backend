package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
)

var aesKey = []byte("")

func encryptKey(key *Key) ([]byte, error) {
	keyJSON, err := json.Marshal(key)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	_, err = rand.Read(iv)
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, iv)

	encrypted := make([]byte, aes.BlockSize+len(keyJSON))
	copy(encrypted[:aes.BlockSize], iv)

	stream.XORKeyStream(encrypted[aes.BlockSize:], keyJSON)

	return encrypted, nil
}

func decryptKey(encrypted []byte) (*Key, error) {
	iv := encrypted[:aes.BlockSize]

	ciphertext := encrypted[aes.BlockSize:]

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, iv)

	decrypted := make([]byte, len(ciphertext))

	stream.XORKeyStream(decrypted, ciphertext)

	var key Key
	err = json.Unmarshal(decrypted, &key)
	if err != nil {
		return nil, err
	}

	return &key, nil
}
