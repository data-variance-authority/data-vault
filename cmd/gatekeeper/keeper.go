package gatekeeper

// Exec initializes the gatekeeper server and starts it
func Exec() {
	Init()
	Server()
}
