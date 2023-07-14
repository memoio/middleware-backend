package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"io"
)

func EncryptFile(input io.Reader, key []byte) (io.Reader, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	stream := cipher.NewCTR(block, iv)

	pr, pw := io.Pipe()

	reader := &cipher.StreamReader{S: stream, R: input}

	go func() {
		_, _ = io.Copy(pw, reader)
		pw.Close()
	}()

	return pr, nil
}

func DecryptFile(input *bytes.Buffer, output io.Writer, key []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)
	stream := cipher.NewCTR(block, iv)

	writer := &cipher.StreamWriter{S: stream, W: output}

	if _, err := io.Copy(writer, input); err != nil {
		return err
	}

	return nil
}
