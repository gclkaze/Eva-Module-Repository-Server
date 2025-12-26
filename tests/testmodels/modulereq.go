// Package testmodels contains models that are used specifically for testing
package testmodels

type ModuleRequest struct {
	Title       string
	Repr        string
	Tags        string
	Description string
	TheFile     string
	FilePath    string
}

type ModuleRequestMultipleFiles struct {
	Title       string
	Repr        string
	Tags        string
	Description string
	TheFiles    []string
	FilePath    []string
}
