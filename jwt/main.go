package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/simba-fs/go-oauth/types"
)

func keyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

func getKeyFromFile(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

// getKey returns private key. If not exists, it will generate a new one and save. If exists, it will load it for you.
func getKey(path string) (*rsa.PrivateKey, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// if key file not exists
		key, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return nil, err
		}

		// save to file
		keyBytes := keyToBytes(key)
		err = os.WriteFile(path, keyBytes, 0644)
		if err != nil {
			return nil, err
		}
		return key, nil

	}
	// if key file is exists
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	key, err := getKeyFromFile(string(keyBytes))
	if err != nil {
		return nil, err
	}

	return key, nil
}

// Sign issue a jwt token
func Sign(token *jwt.Token, config *types.Config) ([]byte, error) {
	(*token).Set(jwt.IssuerKey, config.ClientID)
	(*token).Set(jwt.IssuedAtKey, time.Now().Unix())
	(*token).Set(jwt.ExpirationKey, time.Now().Add(time.Second*time.Duration(config.TokenExp)).Unix())

	key, err := getKey("../private.pem")
	if err != nil {
		return make([]byte, 0), nil
	}

	return jwt.Sign(*token, jwa.RS256, key)
}
