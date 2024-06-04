package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// Encrypt encrypts plain text string using AES in CFB mode and returns base64 encoded string.
func Encrypt(text string, key string) (string, error) {
	// Generate a new AES cipher using the key
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// Create a new byte array the size of the text plus the block size for the IV
	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]

	// Read a random IV
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Encrypt the text using CFB mode
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(text))

	// Return the encoded base64 string
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a base64 encoded string using AES in CFB mode and returns the plain text string.
func Decrypt(encodedText string, key string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encodedText)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	// Generate a new AES cipher using the key
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// Get the IV from the ciphertext
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// Decrypt the text using CFB mode
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
