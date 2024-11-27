package vault

import (
	"log"
	"net/http"
	"os"
)

// Server starts the vault server
func Server() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", HandlerPing) // Ping the vault server

	mux.HandleFunc("GET /groups", HandlerGroups)        // Get all groups
	mux.HandleFunc("GET /group", HandlerGroup)          // Get all records in a group
	mux.HandleFunc("PUT /group", HandlerGroupUpload)    // Upload files into a group
	mux.HandleFunc("DELETE /group", HandlerGroupDelete) // Delete a group

	mux.HandleFunc("GET /group/element", HandlerElementGet)      // Get an element
	mux.HandleFunc("DELETE /group/element", HandleElementDelete) // Delete an element

	// setup server
	server := &http.Server{
		Addr:     ":" + VaultConfig.Port,
		Handler:  mux,
		ErrorLog: log.New(os.Stderr, "http: ", log.LstdFlags),
	}

	//Start the vault server
	log.Println("Starting vault server...")
	log.Println("Listening on port", VaultConfig.Port)
	log.Fatal(server.ListenAndServe())
}
