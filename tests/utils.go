// Copyright (c) 2025 Michail Dorgiakis - gclkaze@gmail.com - https://github.com/gclkaze
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
