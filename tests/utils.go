// Package tests contains the tests!
package tests

import (
	"fmt"
	"os"
	"path"

	"github.com/gclkaze/evamodulerepositoryserver/cmd/server"
	"github.com/gin-gonic/gin"
)

var TheTestServer *server.EvaModuleRepositoryServer

func StartServer() *gin.Engine {
	TheTestServer = server.NewEvaModuleRepositoryServer()
	currentPath, _ := os.Getwd()
	p := path.Join(currentPath, "application_test.properties")
	error := TheTestServer.InitializeWithPropertiesPath(p)
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	} else {
		return TheTestServer.GetRouter()

	}
	return nil
}

func ClearModuleFolders() {
	TheTestServer.ClearModuleFolders()
}

func TeardownServer() {
	if TheTestServer == nil {
		fmt.Print("test Server hasn't been initialized")
		os.Exit(1)
	}

	TheTestServer.CleanDB()
}

func GetDefaultUserPassword() string {
	return TheTestServer.GetBackend().GetProperties().GetString("default_password", "")
}
