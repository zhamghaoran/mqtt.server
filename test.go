package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"math/big"
)

func test() {
	publickeyString := "-----BEGIN PUBLIC KEY-----\nMHYwEAYHKoZIzj0CAQYFK4EEACIDYgAE1Hf//ofxfUAQbi+a/FMvXc8/sS5n21hp\nY2A0dM/+OOtTruT0qLtgMnEtjR2gKqnsnVEgKZwen7vEyiZhy3iSIadrknjz51P1\nC5Is4LI43fFmb4vjyogR+V66efO6QmIC\n-----END PUBLIC KEY-----"
	//p,_ := pem.Decode([]byte(publickeyString))
	keyh, _ := loadPublicKey(publickeyString)
	// The public key shared by Amazon. Replace this with the actual key.
	//publicKeyBytes := []byte(publickeyString)

	// The JSON string. Replace this with the actual string.
	jsonString := `{"vendorId":4933,"productId":40961,"uniqueDeviceId":"MjIwOTIzNDc0MDIyMjY1MjAwMDFmMDAwMDAwMDAwMTc=","rotatingIdAlgorithm":"MATTER_V0","discriminator":23,"passcode":63660476}`

	// Convert JSON string to byte array
	jsonBytes := []byte(jsonString)

	// Create a new elliptic curve
	curve := elliptic.P384()

	// Get x, y coordinates from the public key
	//x, y := Unmarshal(curve, publicKeyBytes)

	//pubkey, _ := crypto.DecompressPubkey(p.Bytes)
	// Create a new ECDSA public key
	publicKey := &ecdsa.PublicKey{Curve: curve, X: keyh.X, Y: keyh.Y}

	// Convert the ECDSA public key to ECIES public key
	eciesPublicKey := ecies.ImportECDSAPublic(publicKey)

	// Encrypt the byte array with the ECIES public key
	encryptedBytes, err := ecies.Encrypt(rand.Reader, eciesPublicKey, jsonBytes, nil, nil)
	if err != nil {
		panic(err)
	}

	// Base64 encode the encrypted byte array
	encodedString := base64.StdEncoding.EncodeToString(encryptedBytes)

	// Print the encoded string
	fmt.Println(encodedString)
}
func Unmarshal(curve elliptic.Curve, data []byte) (x, y *big.Int) {
	byteLen := (curve.Params().BitSize + 7) >> 3
	if len(data) != 1+2*byteLen {
		return
	}

	if data[0] != 4 { // uncompressed form
		return
	}
	p := curve.Params().P

	x = new(big.Int).SetBytes(data[1 : 1+byteLen])
	y = new(big.Int).SetBytes(data[1+byteLen:])
	if x.Cmp(p) >= 0 || y.Cmp(p) >= 0 {
		return nil, nil
	}
	if !curve.IsOnCurve(x, y) {
		return nil, nil
	}
	return
}
func loadPublicKey(publicKeyStr string) (*ecdsa.PublicKey, error) {
	keyData := []byte(publicKeyStr)
	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not of type ECDSA")
	}

	return publicKey, nil
}
