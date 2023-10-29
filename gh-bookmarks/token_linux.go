//go:build linux

package main

import (
	"log"

	"github.com/99designs/keyring"
)

func init() {
	var err error
	ring, err = keyring.Open(keyring.Config{
		ServiceName:     "github-token",
		AllowedBackends: []keyring.BackendType{keyring.KeyCtlBackend},
		KeyCtlScope:     "session",
		// KeyCtlPerm:      0x3f3f0000, // "alswrvalswrv------------"
	})

	if err != nil {
		log.Fatal("Opening keyring failed")
	}
}
