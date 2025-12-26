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
var OldUploadLimit int64

func StartServer() *gin.Engine {
	TheTestServer = server.NewEvaModuleRepositoryServer()
	currentPath, _ := os.Getwd()
	p := path.Join(currentPath, "application_test.properties")
	error := TheTestServer.InitializeWithPropertiesPath(p)
	//TheTestServer.ClearModuleFolders()
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	} else {
		return TheTestServer.GetRouter()

	}
	return nil
}

func ResetUploadLimit(limit int64) int64 {
	currentPath, _ := os.Getwd()
	p := path.Join(currentPath, "application_test.properties")
	old, error := TheTestServer.ResetRouterWithUploadLimit(p, limit)

	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}

	return old
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
	TheTestServer.ClearModuleFolders()
}

func GetDefaultUserPassword() string {
	return TheTestServer.GetBackend().GetProperties().GetString("default_password", "")
}
