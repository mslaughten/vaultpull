// Package main is the entry point for the vaultpull application.
// vaultpull is a CLI tool for pulling secrets from HashiCorp Vault
// and making them available as environment variables or files.
package main

import "github.com/vaultpull/vaultpull/cmd"

func main() {
	cmd.Execute()
}
