package vault

// Exec initializes and starts the vault server
func Exec() {
	Init()
	Server()
}
