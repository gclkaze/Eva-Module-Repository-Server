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

package utils

import (
	"strings"
	"testing"
)

// ============================================================================
// IsValidRepr Tests (Module Representation)
// Min: 2 chars, Max: 64 chars, Allows: letters, numbers, spaces, '-', '_'
// ============================================================================

func TestIsValidRepr_ValidInputs(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Minimum length 2 chars", "AB"},
		{"Exactly 3 chars", "ABC"},
		{"With uppercase and lowercase", "MyModule"},
		{"With numbers", "Module123"},
		{"With hyphens", "My-Module"},
		{"With underscores", "My_Module"},
		{"With spaces", "My Module"},
		{"Mixed separators", "My-Module_Name"},
		{"Maximum length 64", strings.Repeat("A", 64)},
		{"64 chars with separators", "A" + strings.Repeat("-A", 31) + "B"},
		{"Numbers and letters mix", "Module2023"},
		{"All spaces in middle", "A B C"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsValidRepr(tt.input)
			if err != nil {
				t.Errorf("IsValidRepr(%q) = %v, want nil", tt.input, err)
			}
		})
	}
}

func TestIsValidRepr_InvalidInputs(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Empty string", "", true},
		{"Single character", "A", true},
		{"Only spaces", "   ", true},
		{"Too long 65 chars", strings.Repeat("A", 65), true},
		{"Starts with hyphen", "-Module", true},
		{"Starts with underscore", "_Module", true},
		{"Starts with space", " Module", false}, // Gets trimmed
		{"Ends with hyphen", "Module-", true},
		{"Ends with underscore", "Module_", true},
		{"Ends with space", "Module ", false}, // Gets trimmed
		{"Control character tab", "Module\tName", true},
		{"Control character newline", "Module\nName", true},
		{"Invalid special chars", "Module@Name", true},
		{"Invalid special chars exclamation", "Module!Name", true},
		{"Invalid special chars hash", "Module#Name", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsValidRepr(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValidRepr(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

// ============================================================================
// IsValidDescription Tests
// Max: 2000 chars, allows most characters including newlines and tabs
// ============================================================================

func TestIsValidDescription_ValidInputs(t *testing.T) {
	tests := []struct {
		name string
		desc string
	}{
		{"Empty string", ""},
		{"Single character", "A"},
		{"Normal text", "This is a module description"},
		{"With numbers", "Version 1.0 Release"},
		{"With punctuation", "Hello, world! How are you?"},
		{"With newlines", "Line 1\nLine 2\nLine 3"},
		{"With tabs", "Column1\tColumn2\tColumn3"},
		{"Maximum length 2000", strings.Repeat("A", 2000)},
		{"Long text with newlines", strings.Repeat("Line\n", 300)},
		{"Mixed special chars", "Module: Description (v1.0) [BETA]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsValidDescription(tt.desc)
			if err != nil {
				t.Errorf("IsValidDescription(%q) = %v, want nil", tt.desc, err)
			}
		})
	}
}

func TestIsValidDescription_InvalidInputs(t *testing.T) {
	tests := []struct {
		name    string
		desc    string
		wantErr bool
	}{
		{"Too long 2001 chars", strings.Repeat("A", 2001), true},
		{"Control character (null)", "Module\x00Name", true},
		{"Control character (SOH)", "Module\x01Name", true},
		{"Control character (STX)", "Module\x02Name", true},
		{"Control character (BEL)", "Module\x07Name", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsValidDescription(tt.desc)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValidDescription(%q) error = %v, wantErr %v", tt.desc, err, tt.wantErr)
			}
		})
	}
}

// ============================================================================
// IsValidPassword Tests
// Min: 12 chars, Max: 128 chars
// ============================================================================

func TestIsValidPassword_ValidInputs(t *testing.T) {
	tests := []struct {
		name string
		pwd  string
	}{
		{"Minimum length 12", "MyPass12345!"},
		{"With mixed case", "MyPassword123"},
		{"With numbers and letters", "Password1234"},
		{"With special chars", "P@ssw0rd!123"},
		{"Maximum length 128", strings.Repeat("P", 128)},
		{"Long valid password", "MySecurePassword1234567890!@#$%^&*()"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsValidPassword(tt.pwd)
			if err != nil {
				t.Errorf("IsValidPassword(%q) = %v, want nil", tt.pwd, err)
			}
		})
	}
}

func TestIsValidPassword_InvalidInputs(t *testing.T) {
	tests := []struct {
		name    string
		pwd     string
		wantErr bool
	}{
		{"Empty string", "", true},
		{"Too short 11 chars", "MyPass1234", true},
		{"Single char", "A", true},
		{"Too long 129 chars", strings.Repeat("P", 129), true},
		{"Exactly 11 chars", strings.Repeat("P", 11), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsValidPassword(tt.pwd)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValidPassword(%q) error = %v, wantErr %v", tt.pwd, err, tt.wantErr)
			}
		})
	}
}

// ============================================================================
// IsValidModuleNameWithError Tests
// Min: 3 chars, Max: 50 chars
// Allows: letters, digits, dashes, spaces
// ============================================================================

func TestIsValidModuleNameWithError_ValidInputs(t *testing.T) {
	tests := []struct {
		name       string
		moduleName string
	}{
		{"Minimum length 3", "ABC"},
		{"Simple name", "MyModule"},
		{"With spaces", "My Module"},
		{"With numbers", "Module2023"},
		{"Maximum length 50", strings.Repeat("A", 50)},
		{"Lowercase letters", "mymodule"},
		{"Uppercase letters", "MYMODULE"},
		{"With numbers and letters", "Module2023Test"},
		{"Simple with space", "My Module"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsValidModuleNameWithError(tt.moduleName)
			if err != nil {
				t.Errorf("IsValidModuleNameWithError(%q) = %v, want nil", tt.moduleName, err)
			}
		})
	}
}

func TestIsValidModuleNameWithError_InvalidInputs(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Too short 2 chars", "AB", true},
		{"Single char", "A", true},
		{"Too long 51 chars", strings.Repeat("A", 51), true},
		{"Invalid special char @", "Module@Name", true},
		{"Invalid special char #", "Module#Name", true},
		{"Invalid special char $", "Module$Name", true},
		{"Invalid special char %", "Module%Name", true},
		{"Invalid special char &", "Module&Name", true},
		{"Invalid special char !", "Module!Name", true},
		{"Invalid underscore", "Module_Name", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsValidModuleNameWithError(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValidModuleNameWithError(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

// ============================================================================
// IsValidVersion Tests
// Min: 5 chars, Max: 64 chars
// Must be valid semver (with or without leading 'v')
// ============================================================================

func TestIsValidVersion_ValidInputs(t *testing.T) {
	tests := []struct {
		name    string
		version string
	}{
		{"Semver with v prefix", "v0.0.0"},
		{"Semver without v prefix", "0.0.0"},
		{"Major.Minor.Patch", "1.2.3"},
		{"With v and larger numbers", "v10.20.30"},
		{"With pre-release", "v1.0.0-alpha"},
		{"With pre-release and build", "v1.0.0-alpha+001"},
		{"Long valid version", "v123.456.789"},
		{"Prerelease beta", "v2.0.0-beta.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidVersion(tt.version)
			if !result {
				t.Errorf("IsValidVersion(%q) = %v, want true", tt.version, result)
			}
		})
	}
}

func TestIsValidVersion_InvalidInputs(t *testing.T) {
	tests := []struct {
		name    string
		version string
	}{
		{"Empty string", ""},
		{"Too short 4 chars", "v0.0"},
		{"Single digit", "1"},
		{"Invalid format no dots", "version"},
		{"Too long 65 chars", strings.Repeat("A", 65)},
		{"Invalid semver", "not-a-version"},
		{"Only letters", "versionstring"},
		{"Invalid characters", "v1.2.3@@@"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidVersion(tt.version)
			if result {
				t.Errorf("IsValidVersion(%q) = %v, want false", tt.version, result)
			}
		})
	}
}

// ============================================================================
// IsValidName Tests (Personal Names)
// Min: 1 char, Max: 50 chars
// Must match nameRegex: letters only, with optional apostrophes/hyphens
// ============================================================================

func TestIsValidName_ValidInputs(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Single letter", "A"},
		{"Simple name", "John"},
		{"With apostrophe", "O'Brien"},
		{"With hyphen", "Mary-Jane"},
		{"Longer name", "Alexander"},
		{"Maximum 50 chars", strings.Repeat("A", 50)},
		{"With multiple hyphens", "Jean-Pierre-Louis"},
		{"With multiple apostrophes", "O'Brien's"},
		{"Accented letters", "Francçois"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidName(tt.input)
			if !result {
				t.Errorf("IsValidName(%q) = %v, want true", tt.input, result)
			}
		})
	}
}

func TestIsValidName_InvalidInputs(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Empty string", ""},
		{"Too long 51 chars", strings.Repeat("A", 51)},
		{"With numbers", "John123"},
		{"With special chars", "John@Doe"},
		{"Starting with apostrophe", "'John"},
		{"Ending with apostrophe", "John'"},
		{"Starting with hyphen", "-John"},
		{"Ending with hyphen", "John-"},
		{"With underscore", "John_Doe"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidName(tt.input)
			if result {
				t.Errorf("IsValidName(%q) = %v, want false", tt.input, result)
			}
		})
	}
}

// ============================================================================
// IsValidHandle Tests (User Handles/Usernames)
// Min: 3 chars, Max: 30 chars
// Pattern: starts with letter, letters/numbers, with optional ._- separators
// ============================================================================

func TestIsValidHandle_ValidInputs(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Minimum 3 chars", "abc"},
		{"Simple handle", "john"},
		{"With underscore", "john_doe"},
		{"With period", "john.doe"},
		{"With hyphen", "john-doe"},
		{"With numbers", "john123"},
		{"Maximum 30 chars", strings.Repeat("a", 30)},
		{"Mixed case", "JohnDoe"},
		{"Multiple separators", "john_doe.smith"},
		{"30 chars with separators", "a" + strings.Repeat("-a", 14) + "b"},
		{"Uppercase letters", "JOHN"},
		{"Numbers after letters", "john1doe2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidHandle(tt.input)
			if !result {
				t.Errorf("IsValidHandle(%q) = %v, want true", tt.input, result)
			}
		})
	}
}

func TestIsValidHandle_InvalidInputs(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Empty string", ""},
		{"Too short 2 chars", "ab"},
		{"Too long 31 chars", strings.Repeat("a", 31)},
		{"Starting with number", "1john"},
		{"Starting with underscore", "_john"},
		{"Starting with period", ".john"},
		{"Starting with hyphen", "-john"},
		{"Ending with underscore", "john_"},
		{"Ending with period", "john."},
		{"Ending with hyphen", "john-"},
		{"With space", "john doe"},
		{"With special chars @", "john@doe"},
		{"Double underscore", "john__doe"},
		{"Double period", "john..doe"},
		{"Double hyphen", "john--doe"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidHandle(tt.input)
			if result {
				t.Errorf("IsValidHandle(%q) = %v, want false", tt.input, result)
			}
		})
	}
}

// ============================================================================
// IsValidModuleName Tests
// Min: 3 chars, Max: 50 chars
// Allows: letters, numbers, and single separators (space, dash, underscore)
// ============================================================================

func TestIsValidModuleName_ValidInputs(t *testing.T) {
	tests := []struct {
		name   string
		input  string
	}{
		{"Minimum 3 chars", "ABC"},
		{"Simple module", "MyModule"},
		{"With space", "My Module"},
		{"With numbers", "Module123"},
		{"Maximum 50 chars", strings.Repeat("A", 50)},
		{"Lowercase", "mymodule"},
		{"Uppercase", "MYMODULE"},
		{"Numbers and letters", "Module2023Test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidModuleName(tt.input)
			if !result {
				t.Errorf("IsValidModuleName(%q) = %v, want true", tt.input, result)
			}
		})
	}
}

func TestIsValidModuleName_InvalidInputs(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Too short 2 chars", "AB"},
		{"Single char", "A"},
		{"Too long 51 chars", strings.Repeat("A", 51)},
		{"Empty string", ""},
		{"Invalid special @", "Module@Name"},
		{"Invalid special #", "Module#Name"},
		{"Invalid special $", "Module$Name"},
		{"With dash", "Module-Name"},
		{"With underscore", "Module_Name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidModuleName(tt.input)
			if result {
				t.Errorf("IsValidModuleName(%q) = %v, want false", tt.input, result)
			}
		})
	}
}

// ============================================================================
// ParseModuleReleaseVersion Tests
// Format: "module" or "module@version" or "module@latest"
// ============================================================================

func TestParseModuleReleaseVersion_ValidInputs(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantModule  string
		wantVersion string
		wantErr     bool
	}{
		{"Module only", "mymodule", "mymodule", "", false},
		{"With latest", "mymodule@latest", "mymodule", "latest", false},
		{"With version v prefix", "mymodule@v1.0.0", "mymodule", "v1.0.0", false},
		{"With version no v", "mymodule@1.0.0", "mymodule", "v1.0.0", false},
		{"With prerelease", "mymodule@v1.0.0-alpha", "mymodule", "v1.0.0-alpha", false},
		{"Complex module name", "my-module@v2.1.0", "my-module", "v2.1.0", false},
		{"Numbers in module", "module2023@v1.2.3", "module2023", "v1.2.3", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			module, version, err := ParseModuleReleaseVersion(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseModuleReleaseVersion(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if module != tt.wantModule {
				t.Errorf("ParseModuleReleaseVersion(%q) module = %q, want %q", tt.input, module, tt.wantModule)
			}
			if version != tt.wantVersion {
				t.Errorf("ParseModuleReleaseVersion(%q) version = %q, want %q", tt.input, version, tt.wantVersion)
			}
		})
	}
}

func TestParseModuleReleaseVersion_InvalidInputs(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Empty string", "", true},
		{"Empty module", "@latest", true},
		{"Empty version", "mymodule@", true},
		{"Multiple @ symbols", "mymodule@v1.0.0@extra", true},
		{"Invalid version format", "mymodule@invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := ParseModuleReleaseVersion(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseModuleReleaseVersion(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

// ============================================================================
// StringToUint and UintToString Tests
// Bidirectional string<->uint conversion
// ============================================================================

func TestStringToUint_ValidInputs(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  uint
	}{
		{"Zero", "0", 0},
		{"Single digit", "5", 5},
		{"Two digits", "42", 42},
		{"Large number", "999999999", 999999999},
		{"Maximum uint32", "4294967295", 4294967295},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := StringToUint(tt.input)
			if err != nil {
				t.Errorf("StringToUint(%q) error = %v, want nil", tt.input, err)
			}
			if result != tt.want {
				t.Errorf("StringToUint(%q) = %v, want %v", tt.input, result, tt.want)
			}
		})
	}
}

func TestStringToUint_InvalidInputs(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Negative number", "-1"},
		{"Not a number", "abc"},
		{"Empty string", ""},
		{"With spaces", "123 456"},
		{"Floating point", "123.45"},
		{"With letters mixed", "123abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := StringToUint(tt.input)
			if err == nil {
				t.Errorf("StringToUint(%q) error = nil, want error", tt.input)
			}
		})
	}
}

func TestUintToString_ConversionRoundTrip(t *testing.T) {
	tests := []uint{0, 1, 42, 999, 123456789}

	for _, tt := range tests {
		t.Run("RoundTrip_"+string(rune(tt)), func(t *testing.T) {
			str := UintToString(tt)
			back, err := StringToUint(str)
			if err != nil {
				t.Errorf("StringToUint(%q) error = %v, want nil", str, err)
			}
			if back != tt {
				t.Errorf("RoundTrip failed: %v -> %q -> %v", tt, str, back)
			}
		})
	}
}

// ============================================================================
// HashPassword Tests
// ============================================================================

func TestHashPassword_ValidInputs(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{"Simple password", "MySecurePassword123!"},
		{"Long password 72 bytes", strings.Repeat("P", 72)},
		{"Special characters", "P@$$w0rd!#%&*"},
		{"With spaces", "My Secure Password 123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if err != nil {
				t.Errorf("HashPassword(%q) error = %v, want nil", tt.password, err)
			}
			if hash == "" {
				t.Errorf("HashPassword(%q) returned empty hash", tt.password)
			}
			if hash == tt.password {
				t.Errorf("HashPassword(%q) = %q, hash should be different from input", tt.password, hash)
			}
		})
	}
}

// ============================================================================
// ValidatePassword Tests
// Min: 12 chars, Max: 128 chars, checks against common passwords and email
// ============================================================================

func TestValidatePassword_ValidInputs(t *testing.T) {
	// Enable quality checks
	oldDisabled := PasswordQualityCheckDisabled
	PasswordQualityCheckDisabled = false
	defer func() { PasswordQualityCheckDisabled = oldDisabled }()

	tests := []struct {
		name     string
		password string
		email    string
	}{
		{"Valid password with email", "MySecurePassword123!", "user@example.com"},
		{"Valid password without email", "MySecurePassword123!", ""},
		{"Long valid password", strings.Repeat("P", 128), "user@example.com"},
		{"With special chars", "P@ssw0rd!#%&*123", "user@example.com"},
		{"Minimum length 12", "MyPassword123", "user@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password, tt.email)
			if err != nil {
				t.Errorf("ValidatePassword(%q, %q) = %v, want nil", tt.password, tt.email, err)
			}
		})
	}
}

func TestValidatePassword_InvalidInputs(t *testing.T) {
	// Enable quality checks
	oldDisabled := PasswordQualityCheckDisabled
	PasswordQualityCheckDisabled = false
	defer func() { PasswordQualityCheckDisabled = oldDisabled }()

	tests := []struct {
		name     string
		password string
		email    string
		wantErr  bool
	}{
		{"Too short", "Pass123", "user@example.com", true},
		{"Too long", strings.Repeat("P", 129), "user@example.com", true},
		{"Common password 'password'", "password", "user@example.com", true},
		{"Common password '123456'", "123456", "user@example.com", true},
		{"Common password 'qwerty'", "qwerty", "user@example.com", true},
		{"Common password 'admin'", "admin", "user@example.com", true},
		{"Contains username", "ValidUserPassword123", "user@example.com", true},
		{"Contains username case insensitive", "ValidUSERPassword123", "user@example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password, tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword(%q, %q) error = %v, wantErr %v", tt.password, tt.email, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePassword_QualityCheckDisabled(t *testing.T) {
	// Disable quality checks
	oldDisabled := PasswordQualityCheckDisabled
	PasswordQualityCheckDisabled = true
	defer func() { PasswordQualityCheckDisabled = oldDisabled }()

	// Even common passwords should pass when disabled
	err := ValidatePassword("password", "user@example.com")
	if err != nil {
		t.Errorf("ValidatePassword with disabled checks should pass, got error: %v", err)
	}
}

// ============================================================================
// GetRandomNumber Tests
// Returns value in range [0, max)
// ============================================================================

func TestGetRandomNumber_ReturnsBoundedValue(t *testing.T) {
	for i := 0; i < 100; i++ {
		result := GetRandomNumber(10)
		if result < 0 || result >= 10 {
			t.Errorf("GetRandomNumber(10) = %d, want value in [0, 10)", result)
		}
	}
}

func TestGetRandomNumber_LargeRange(t *testing.T) {
	for i := 0; i < 100; i++ {
		result := GetRandomNumber(1000000)
		if result < 0 || result >= 1000000 {
			t.Errorf("GetRandomNumber(1000000) = %d, want value in [0, 1000000)", result)
		}
	}
}

// ============================================================================
// GetRandomUintNumber Tests
// Returns value in range [0, max)
// ============================================================================

func TestGetRandomUintNumber_ReturnsBoundedValue(t *testing.T) {
	for i := 0; i < 100; i++ {
		result := GetRandomUintNumber(10)
		if result >= 10 {
			t.Errorf("GetRandomUintNumber(10) = %d, want value in [0, 10)", result)
		}
	}
}

func TestGetRandomUintNumber_LargeRange(t *testing.T) {
	for i := 0; i < 100; i++ {
		result := GetRandomUintNumber(1000000)
		if result >= 1000000 {
			t.Errorf("GetRandomUintNumber(1000000) = %d, want value in [0, 1000000)", result)
		}
	}
}

// ============================================================================
// GetRandomNumberInRange Tests
// Returns value in range [min, max]
// ============================================================================

func TestGetRandomNumberInRange_ValidRange(t *testing.T) {
	for i := 0; i < 100; i++ {
		result := GetRandomNumberInRange(5, 10)
		if result < 5 || result > 10 {
			t.Errorf("GetRandomNumberInRange(5, 10) = %d, want value in [5, 10]", result)
		}
	}
}

func TestGetRandomNumberInRange_SingleValue(t *testing.T) {
	result := GetRandomNumberInRange(5, 5)
	if result != 5 {
		t.Errorf("GetRandomNumberInRange(5, 5) = %d, want 5", result)
	}
}

func TestGetRandomNumberInRange_InvalidRange(t *testing.T) {
	result := GetRandomNumberInRange(10, 5)
	if result != -1 {
		t.Errorf("GetRandomNumberInRange(10, 5) = %d, want -1", result)
	}
}

func TestGetRandomNumberInRange_LargeRange(t *testing.T) {
	for i := 0; i < 100; i++ {
		result := GetRandomNumberInRange(0, 1000000)
		if result < 0 || result > 1000000 {
			t.Errorf("GetRandomNumberInRange(0, 1000000) = %d, want value in [0, 1000000]", result)
		}
	}
}

// ============================================================================
// GetRandomUintRange Tests
// Returns value in range [min, max]
// ============================================================================

func TestGetRandomUintRange_ValidRange(t *testing.T) {
	for i := 0; i < 100; i++ {
		result := GetRandomUintRange(5, 10)
		if result < 5 || result > 10 {
			t.Errorf("GetRandomUintRange(5, 10) = %d, want value in [5, 10]", result)
		}
	}
}

func TestGetRandomUintRange_SingleValue(t *testing.T) {
	result := GetRandomUintRange(5, 5)
	if result != 5 {
		t.Errorf("GetRandomUintRange(5, 5) = %d, want 5", result)
	}
}

func TestGetRandomUintRange_InvalidRange(t *testing.T) {
	result := GetRandomUintRange(10, 5)
	if result != 0 {
		t.Errorf("GetRandomUintRange(10, 5) = %d, want 0", result)
	}
}

// ============================================================================
// GetRepoName Tests
// Generates URL-friendly slug from input
// ============================================================================

func TestGetRepoName_ValidInputs(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Simple word", "MyModule"},
		{"With spaces", "My Module"},
		{"With hyphens", "My-Module"},
		{"With uppercase", "MYMODULE"},
		{"Mixed case with space", "My Super Module"},
		{"With special chars", "My@Module!"},
		{"Numbers", "Module2023"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetRepoName(tt.input)
			if result == "" {
				t.Errorf("GetRepoName(%q) returned empty string", tt.input)
			}
			// Result should be lowercase and slug-like
			if result != strings.ToLower(result) && strings.Contains(result, " ") {
				t.Errorf("GetRepoName(%q) should be lowercase slug format", tt.input)
			}
		})
	}
}

// ============================================================================
// IsValidModuleOrModuleVersion Tests
// Validates either module name alone or module@version format
// ============================================================================

func TestIsValidModuleOrModuleVersion_ValidInputs(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Module only", "mymodule"},
		{"Module with latest", "mymodule@latest"},
		{"Module with version", "mymodule@v1.0.0"},
		{"Module with version no v", "mymodule@1.0.0"},
		{"Complex module name", "my-module@v2.1.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := IsValidModuleOrModuleVersion(tt.input)
			if err != nil {
				t.Errorf("IsValidModuleOrModuleVersion(%q) error = %v, want nil", tt.input, err)
			}
			if !result {
				t.Errorf("IsValidModuleOrModuleVersion(%q) = %v, want true", tt.input, result)
			}
		})
	}
}

func TestIsValidModuleOrModuleVersion_InvalidInputs(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Empty string", ""},
		{"Invalid module name", "a"},
		{"Invalid version", "mymodule@1.0"},
		{"Invalid semver", "mymodule@invalid"},
		{"Empty module", "@v1.0.0"},
		{"Empty version", "mymodule@"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := IsValidModuleOrModuleVersion(tt.input)
			if err == nil {
				t.Errorf("IsValidModuleOrModuleVersion(%q) error = nil, want error", tt.input)
			}
			if result {
				t.Errorf("IsValidModuleOrModuleVersion(%q) = %v, want false", tt.input, result)
			}
		})
	}
}

// ============================================================================
// ValidateTags Tests
// Max: 10 tags, each tag: 1-32 chars, alphanumeric with - and _
// ============================================================================

func TestValidateTags_ValidInputs(t *testing.T) {
	tests := []struct {
		name string
		tags []string
	}{
		{"Empty slice", []string{}},
		{"Single tag", []string{"tag"}},
		{"Multiple tags", []string{"tag1", "tag2", "tag3"}},
		{"Tags with hyphen", []string{"my-tag", "another-tag"}},
		{"Tags with underscore", []string{"my_tag", "another_tag"}},
		{"Tags with numbers", []string{"tag1", "tag2"}},
		{"Maximum 10 tags", []string{"t1", "t2", "t3", "t4", "t5", "t6", "t7", "t8", "t9", "t10"}},
		{"Minimum length 1 char", []string{"a"}},
		{"Maximum length 32 chars", []string{strings.Repeat("a", 32)}},
		{"Mixed valid tags", []string{"python", "golang", "webapp", "v1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTags(tt.tags)
			if err != nil {
				t.Errorf("ValidateTags(%v) = %v, want nil", tt.tags, err)
			}
		})
	}
}

func TestValidateTags_InvalidInputs(t *testing.T) {
	tests := []struct {
		name    string
		tags    []string
		wantErr bool
	}{
		{"Too many tags 11", []string{"t1", "t2", "t3", "t4", "t5", "t6", "t7", "t8", "t9", "t10", "t11"}, true},
		{"Tag too short empty", []string{"", "tag"}, true},
		{"Tag too long 33 chars", []string{strings.Repeat("a", 33)}, true},
		{"Tag with space", []string{"my tag"}, true},
		{"Tag with special char @", []string{"my@tag"}, true},
		{"Tag with period", []string{"my.tag"}, true},
		{"Tag with control char", []string{"my\ttag"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTags(tt.tags)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTags(%v) error = %v, wantErr %v", tt.tags, err, tt.wantErr)
			}
		})
	}
}

// ============================================================================
// ValidateNames Tests
// Max: 5 names, each name: 1-64 chars, letters/numbers/spaces/hyphens/underscores
// ============================================================================

func TestValidateNames_ValidInputs(t *testing.T) {
	tests := []struct {
		name  string
		names []string
	}{
		{"Empty slice", []string{}},
		{"Single name", []string{"John"}},
		{"Multiple names", []string{"John", "Jane", "Bob"}},
		{"Names with space", []string{"John Doe", "Jane Smith"}},
		{"Names with hyphen", []string{"Mary-Jane"}},
		{"Names with underscore", []string{"john_doe"}},
		{"Names with numbers", []string{"Person1", "Person2"}},
		{"Maximum 5 names", []string{"Name1", "Name2", "Name3", "Name4", "Name5"}},
		{"Minimum 1 char", []string{"A"}},
		{"Maximum 64 chars", []string{strings.Repeat("A", 64)}},
		{"Mixed valid names", []string{"John Doe", "Mary-Jane", "bob_smith"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNames(tt.names)
			if err != nil {
				t.Errorf("ValidateNames(%v) = %v, want nil", tt.names, err)
			}
		})
	}
}

func TestValidateNames_InvalidInputs(t *testing.T) {
	tests := []struct {
		name    string
		names   []string
		wantErr bool
	}{
		{"Too many names 6", []string{"N1", "N2", "N3", "N4", "N5", "N6"}, true},
		{"Name too short empty", []string{""}, true},
		{"Name too long 65 chars", []string{strings.Repeat("A", 65)}, true},
		{"Name with special char @", []string{"Name@Doe"}, true},
		{"Name with control char", []string{"Name\tDoe"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNames(tt.names)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNames(%v) error = %v, wantErr %v", tt.names, err, tt.wantErr)
			}
		})
	}
}

// ============================================================================
// ValidateDescriptions Tests
// Max: 5 descriptions, each: 1-128 chars, letters/numbers/spaces/hyphens/underscores
// ============================================================================

func TestValidateDescriptions_ValidInputs(t *testing.T) {
	tests := []struct {
		name         string
		descriptions []string
	}{
		{"Empty slice", []string{}},
		{"Single description", []string{"This is a description"}},
		{"Multiple descriptions", []string{"Desc1", "Desc2", "Desc3"}},
		{"Description with numbers", []string{"Version 1.0"}},
		{"Description with hyphens", []string{"My-Description"}},
		{"Description with underscores", []string{"my_description"}},
		{"Description with spaces", []string{"This is a long description"}},
		{"Maximum 5 descriptions", []string{"D1", "D2", "D3", "D4", "D5"}},
		{"Minimum 1 char", []string{"A"}},
		{"Maximum 128 chars", []string{strings.Repeat("A", 128)}},
		{"Mixed valid descriptions", []string{"Main description", "Sub-description", "detail_info"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDescriptions(tt.descriptions)
			if err != nil {
				t.Errorf("ValidateDescriptions(%v) = %v, want nil", tt.descriptions, err)
			}
		})
	}
}

func TestValidateDescriptions_InvalidInputs(t *testing.T) {
	tests := []struct {
		name        string
		descriptions []string
		wantErr     bool
	}{
		{"Too many descriptions 6", []string{"D1", "D2", "D3", "D4", "D5", "D6"}, true},
		{"Description too short empty", []string{""}, true},
		{"Description too long 129 chars", []string{strings.Repeat("A", 129)}, true},
		{"Description with special char @", []string{"Desc@ription"}, true},
		{"Description with control char", []string{"Desc\tription"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDescriptions(tt.descriptions)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDescriptions(%v) error = %v, wantErr %v", tt.descriptions, err, tt.wantErr)
			}
		})
	}
}
