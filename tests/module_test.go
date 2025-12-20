package tests

import (
	"fmt"
	"os"
	"testing"
)

func setup() {
	StartServer()
}

func TestMain(m *testing.M) {
	// 🔧 setup (runs once)
	setup()
	// ▶ run all tests
	code := m.Run()

	// 🧹 teardown (runs once)
	teardown()
	os.Exit(code)
}
func teardown() {
	TeardownServer()
}

func TestModuleCreation(t *testing.T) {
	fmt.Printf("lets see")
}
