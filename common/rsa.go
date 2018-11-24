package common

/*
A simple helper class that deals all the RSA encryption.
*/

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"log"
)

// Create RSA key pair in PEM format - used by seller only
func GenerateRSA() (*rsa.PrivateKey, rsa.PublicKey) {
	reader := rand.Reader
	bitSize := 2048
	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		log.Fatalf("Error generating key: %v", err)
		// TODO: handle error
	}
	return key, key.PublicKey
}

// Marshal rsa public key
func MarshalKeyToPem(key rsa.PublicKey) []byte {
	asn1Bytes, err := asn1.Marshal(key)
	if err != nil {
		log.Fatalf("Error on encoding to pem: %v")
		// TODO: handle error
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	})
}

// Unmarshal rsa public key
func UnmarshalPemToKey(rawKey []byte) rsa.PublicKey {
	block, _ := pem.Decode(rawKey)
	if block == nil || block.Type != "PUBLIC KEY" {
		log.Fatalf("Error decoding the key")
		// TODO: handle error
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatalf("Error decoding: %v", err)
		// TODO: handle error
	}
	pk, ok := pub.(rsa.PublicKey)
	if !ok {
		log.Fatalf("type assertion failed")
		// TODO: handle error
	}
	return pk
}

func EncryptID(ipAddress, price string, publicKey *rsa.PublicKey) ([]byte, error) {
	// ID will be bidder's IP address + price encrypted in seller's public key
	// EX: "127.0.0.1:9091 300" -> encrypted with public key
	// We use OAEP encryption standrad and NOT PKCK1
	plainByte := []byte(ipAddress + " " + price)
	rng := rand.Reader
	return rsa.EncryptOAEP(sha256.New(), rng, publicKey, plainByte, []byte(""))
}

func DecryptID(rawMsg []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	rng := rand.Reader
	return rsa.DecryptOAEP(sha256.New(), rng, privateKey, rawMsg, []byte(""))
}
