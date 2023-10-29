package main

import (
	"github.com/99designs/keyring"
)

var ring keyring.Keyring

// storeToken securely stores the GitHub access token using keychain
func storeToken(token string) error {
	return ring.Set(keyring.Item{
		Key:  "token",
		Data: []byte(token),
	})
}

// getToken retrieves the GitHub access token from keychain
func getToken() ([]byte, error) {
	i, err := ring.Get("token")
	if err != nil {
		return []byte{}, err
	}
	return i.Data, nil
}

// deleteToken removes the GitHub access token using keychain
func deleteToken() error {
	return ring.Remove("token")
}
