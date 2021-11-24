// wallet.go
package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"

	"golang.org/x/crypto/ripemd160"
)

const (
	checksumLength = 4
	// hexadecimal for 0
	version = byte(0x00)
)

type Wallet struct {
	// elliptical curve digital signature algorithm
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pub
}

// to get PublicKey, private key -> sha256 -> RipeMD160
func PublicKeyHash(publicKey []byte) []byte {
	hashedPublicKey := sha256.Sum256(publicKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(hashedPublicKey[:])
	if err != nil {
		log.Panic(err)
	}
	publicRipeMd := hasher.Sum(nil)

	return publicRipeMd
}

func Checksum(ripMdHash []byte) []byte {
	firstHash := sha256.Sum256(ripeMdHash)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checksumLength]
}

func (w *Wallet) Address() []byte {
	// steps 1 and 2 get public key
	pubHash := PublicKeyHash(w.PublicKey)
	// step 3 append Version variable
	versionedHash := append([]byte{version}, pubHash...)
	// step 4 get checksum
	checksum := Checksum(versionedHash)
	// step 5
	finalHash := append(versionedHash, checksum...)
	// step 6
	address := base58Encode(finalHash)
	return address
}

func MakeWallet() *Wallet {
	privateKey, publicKey := NewKeyPair()
	wallet := Wallet{privateKey, publicKey}
	return &wallet
}
