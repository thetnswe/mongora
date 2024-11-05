package goNest

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net"
	"strings"
)

func ComparePassword(password, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		// Passwords do not match
		return false, nil
	}
	// Passwords match
	return true, nil
}

// Encrypt function with PKCS7 padding
func Encrypt(textToEncrypt, key, iv string) (string, error) {
	// Convert inputs to byte slices
	plaintext := []byte(textToEncrypt)

	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		return "", fmt.Errorf("failed to decode key: %v", err)
	}

	ivBytes, err := hex.DecodeString(iv)
	if err != nil {
		return "", fmt.Errorf("failed to decode iv: %v", err)
	}

	// Initialize AES cipher with CBC mode
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	// Add PKCS7 padding to the plaintext
	plaintext = pkcs7Pad(plaintext, aes.BlockSize)

	// Encrypt the padded plaintext using CBC
	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, ivBytes)
	mode.CryptBlocks(ciphertext, plaintext)

	// Return hex-encoded encrypted content
	return hex.EncodeToString(ciphertext), nil
}

// Decrypt encrypted content using AES-256-CBC and PKCS7 unpadding
func Decrypt(encryptedContent, key, iv string) (string, error) {
	// Decode the hex-encoded strings

	encryptedBytes, err := hex.DecodeString(strings.TrimSpace(encryptedContent))
	if err != nil {
		return "", fmt.Errorf("failed to decode encrypted content %s : %v", encryptedContent, err)
	}

	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		return "", fmt.Errorf("failed to decode key: %v", err)
	}

	ivBytes, err := hex.DecodeString(iv)
	if err != nil {
		return "", fmt.Errorf("failed to decode iv: %v", err)
	}

	// Initialize AES cipher with CBC mode
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	// Decrypt the content using CBC
	mode := cipher.NewCBCDecrypter(block, ivBytes)
	decrypted := make([]byte, len(encryptedBytes))
	mode.CryptBlocks(decrypted, encryptedBytes)

	// Remove padding (PKCS7)
	decrypted, err = pkcs7Unpad(decrypted, aes.BlockSize)
	if err != nil {
		return "", fmt.Errorf("failed to unpad decrypted content: %v", err)
	}

	return string(decrypted), nil
}

// PKCS7 padding function
func pkcs7Pad(data []byte, blockSize int) []byte {
	paddingLength := blockSize - len(data)%blockSize
	padding := make([]byte, paddingLength)
	for i := range padding {
		padding[i] = byte(paddingLength)
	}
	return append(data, padding...)
}

// PKCS7 unpadding function
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 || length%blockSize != 0 {
		return nil, fmt.Errorf("invalid padding size")
	}

	paddingLength := int(data[length-1])
	if paddingLength > blockSize || paddingLength == 0 {
		return nil, fmt.Errorf("invalid padding length")
	}

	for i := 0; i < paddingLength; i++ {
		if data[length-1-i] != byte(paddingLength) {
			return nil, fmt.Errorf("invalid padding byte")
		}
	}

	return data[:length-paddingLength], nil
}

// GetMacAddresses returns a list of MAC addresses from network interfaces
func GetMacAddresses() []string {
	var macAddresses []string
	interfaces, _ := net.Interfaces()
	for _, iFace := range interfaces {
		if iFace.Flags&net.FlagUp != 0 && (iFace.Name == "eth0" || iFace.Name == "eth1" || iFace.Name == "eth2" || iFace.Name == "Local Area Connection" || iFace.Name == "Ethernet 2" || iFace.Name == "Wi-Fi" || iFace.Name == "Wireless Network Connection") {
			mac := iFace.HardwareAddr.String()
			if mac != "" {
				macAddresses = append(macAddresses, mac)
			}
		}
	}
	return macAddresses
}
