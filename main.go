package main

import (
	"fmt"

	"github.com/gclkaze/evamodulerepositoryserver/server"
)

var TheServer *server.EvaModuleRepositoryServer

func main() {
	TheServer = server.NewEvaModuleRepositoryServer()
	error := TheServer.Initialize()
	if error != nil {
		fmt.Println(error)
	} else {
		TheServer.Run()
	}
}
