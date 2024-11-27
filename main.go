package main

import (
	"datavault/cmd/gatekeeper"
	"datavault/cmd/vault"
	"datavault/configs"
	"log"
)

// Main functions
func main() {
	configs.Instance.Init()

	// Run the command
	switch configs.Instance.Command {
	case "gatekeeper":
		gatekeeper.Exec()
	case "vault":
		vault.Exec()
	default:
		log.Fatalln("Command not recognized")
	}

	log.Printf("Application %s shut down successfully\n", configs.Instance.Command)
}
