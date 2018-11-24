package common

/*
A simple helper class that deals all the RSA encryption.
*/

import (
    "os"
    "log"
    "crypto/rand"
    "crypto/rsa"
    //"crypto/x509"
    "encoding/pem"
    "encoding/asn1"
)

// Create RSA key pair in PEM format
func GenerateRSA() (*rsa.PrivateKey, rsa.PublicKey) {
    reader := rand.Reader
    bitSize := 2048
    key, err := rsa.GenerateKey(reader, bitSize)
    if err != nil {
        log.Fatalf("Error generating key: %v", err)
        os.Exit(1)
    }
    return key, key.PublicKey
}

// Marshal/Unmarshal RSA keys
func MarshalKeyToPem(key rsa.PublicKey) []byte {
    asn1Bytes, err := asn1.Marshal(key)
    if err != nil {
        log.Fatalf("Error on encoding to pem: %v")
    }
    return pem.EncodeToMemory(&pem.Block {
        Type: "PUBLIC KEY",
        Bytes: asn1Bytes,
    })
}

