package main

import (
	"fmt"
	"os"

	"github.com/gclkaze/evamodulerepositoryserver/server"
)

var TheServer *server.EvaModuleRepositoryServer

func main() {
	TheServer = server.NewEvaModuleRepositoryServer()
	error := TheServer.Initialize()
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	} else {
		TheServer.Run()
		os.Exit(0)
	}
}
