package app

import (
	"github.com/gochat/internal/server"
)

// This function is boostrap
// Start following service
// Log // use exe or out in directory with Log folder in it
// Server
func Run() {
	server.Log()
	server.Main()
}
