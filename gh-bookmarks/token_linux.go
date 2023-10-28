//go:build linux

package main

import (
	"os/exec"
)

// storeToken securely stores the GitHub access token using keyctl
func storeToken(token string) error {
	addSecret := exec.Command("keyctl", "add", "user", "github-token", token, "@s")
	return addSecret.Run()
}

// getToken retrieves the GitHub access token from keyctl
func getToken() ([]byte, error) {
	getSecret := exec.Command("keyctl", "pipe", "%user:github-token")
	return getSecret.Output()
}

// deleteToken removes the GitHub access token using keyctl
func deleteToken() error {
	deleteSecret := exec.Command("keyctl", "purge", "-s", "user", "github-token")
	return deleteSecret.Run()
}
