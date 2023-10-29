//go:build darwin

package main

import (
	"log"

	"github.com/99designs/keyring"
)

func init() {
	var err error
	ring, err = keyring.Open(keyring.Config{
		ServiceName:     "github.com",
		AllowedBackends: []keyring.BackendType{keyring.KeychainBackend},
		KeychainName:    "login",
	})

	if err != nil {
		log.Fatal("Opening keyring failed")
	}
}
