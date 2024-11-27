package configs

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// Configuration is the configuration for the application
type Configuration struct {
	Command        string // Command is the command to run
	ConfigFilePath string // ConfigFilePath is the filepath of the configuration, the configuration file is different for each command
	ConfigFileData []byte // RefFile is the file to use for the configuration
}

var Instance Configuration

// Init initializes the configuration
func (c *Configuration) Init() {
	// Initialize the application
	flag.StringVar(&c.Command, "cmd", "", "Command to run (required, options: 'keeper', 'vault')")
	flag.StringVar(&c.ConfigFilePath, "ref", "", "Reference file path to use for configuration (required)")
	help := flag.Bool("help", false, "Show help")
	flag.Parse()

	if *help {
		CliHelper()
		os.Exit(0)
	}

	if c.Command == "" || (c.Command != "gatekeeper" && c.Command != "vault") {
		log.Println("Valid Command is required")
		CliHelper()
		os.Exit(1)
	}

	if c.ConfigFilePath == "" {
		log.Println("Valid Path for Config File is required")
		CliHelper()
		os.Exit(1)
	}

	// Read the file
	tbytes, err := os.ReadFile(c.ConfigFilePath)
	if err != nil {
		log.Println("Error reading configuration file")
		os.Exit(1)
	}
	c.ConfigFileData = tbytes
}

// CliHelper prints the help message for the CLI
func CliHelper() {
	fmt.Printf("Data Vault by D.V.A.\n")
	//detect os
	if os.Getenv("OS") == "Windows_NT" {
		fmt.Println("Usage: datavault.exe -cmd <command> -ref <reference file path>")
	} else {
		fmt.Println("Usage: ./datavault -cmd <command> -ref <reference file path>")
	}

	flag.PrintDefaults()
}
