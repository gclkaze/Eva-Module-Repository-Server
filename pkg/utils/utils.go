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
	"errors"
	"fmt"
	"math/rand/v2"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gosimple/slug"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/mod/semver"
)

var nameRegex = regexp.MustCompile(`^[\p{L}]+([\p{L}\s'-]*[\p{L}]+)?$`)
var handleRegex = regexp.MustCompile(
	`^[a-zA-Z][a-zA-Z0-9]*(?:[._-][a-zA-Z0-9]+)*$`,
)
var reprRegex = regexp.MustCompile(`^[A-Za-z0-9]+( [A-Za-z0-9]+)*$`)
var ErrInvalidModuleReleaseVersion = errors.New("invalid module-release version")
var PasswordLengthMin = 12
var PasswordLengthMax = 128
var ModuleReprMin = 3
var ModuleReprMax = 50
var PasswordQualityCheckDisabled = false

var repoRegex = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
var MinVersionLength = 5
var MaxVersionLength = 64

const (
	MaxDescriptionLength = 2000
	ModuleReprMinLength  = 2
	ModuleReprMaxLength  = 64
)

var moduleReprRe = regexp.MustCompile(
	`^[A-Za-z0-9]+([ _-][A-Za-z0-9]+)*$`,
)

func IsValidRepr(input string) error {
	name := strings.TrimSpace(input)

	if name == "" {
		return fmt.Errorf("module name cannot be empty")
	}

	if len(name) < ModuleReprMinLength {
		return fmt.Errorf(
			"module name too short (min %d characters)",
			ModuleReprMinLength,
		)
	}

	if len(name) > ModuleReprMaxLength {
		return fmt.Errorf(
			"module name too long (max %d characters)",
			ModuleReprMaxLength,
		)
	}

	// Reject control characters
	for _, r := range name {
		if r < 32 {
			return fmt.Errorf(
				"module name contains invalid control characters",
			)
		}
	}

	if !moduleReprRe.MatchString(name) {
		return fmt.Errorf(
			"module name may contain only letters, numbers, spaces, '-' and '_'",
		)
	}

	return nil
}
func IsValidDescription(desc string) error {
	// Allow empty description
	if desc == "" {
		return nil
	}

	// Trim outer whitespace for validation only
	s := strings.TrimSpace(desc)

	// Length check
	if len(s) > MaxDescriptionLength {
		return fmt.Errorf(
			"description too long (%d characters, max %d)",
			len(s),
			MaxDescriptionLength,
		)
	}

	// Control character check
	for i, r := range s {
		if r < 32 && r != '\n' && r != '\t' {
			return fmt.Errorf(
				"description contains invalid control character at position %d",
				i,
			)
		}
	}

	return nil
}
func IsValidModuleOrModuleVersion(s string) (bool, error) {
	module, version, err := ParseModuleReleaseVersion(s)
	if err != nil {
		return false, err
	}

	err = IsValidModuleNameWithError(module)
	if err != nil {
		return false, err
	}

	if version == "" || version == "latest" {
		return true, nil
	}
	if !IsValidVersion(version) {
		return false, fmt.Errorf("incorrect version: '%s'. Additionally its length needs to be between %d and %d", version, MinVersionLength, MaxVersionLength)
	}

	return true, nil
}

func IsValidPassword(pwd string) error {
	length := utf8.RuneCountInString(pwd)
	if length < PasswordLengthMin {
		return fmt.Errorf("password must be at least %d characters long", PasswordLengthMin)
	}
	if length > PasswordLengthMax {
		return errors.New("password is too long")
	}
	return nil
}

func IsValidModuleNameWithError(moduleName string) error {
	attr := "module name"
	if len(moduleName) < ModuleReprMin {
		return fmt.Errorf("%s must be at least %d characters long, it is %d", attr, ModuleReprMin, len(moduleName))
	}
	if len(moduleName) > ModuleReprMax {
		return fmt.Errorf("%s must be at least %d characters long, it is %d", attr, ModuleReprMax, len(moduleName))
	}

	res := reprRegex.MatchString(moduleName)
	if !res {
		res = repoRegex.MatchString(moduleName)
		if !res {
			return fmt.Errorf("only letters,digits and dashes are allowed in %s", attr)
		}
	}

	return nil
}

func StringToUint(str string) (uint, error) {
	// Convert string to uint64 first
	num, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, err
	}

	// Cast to uint
	u := uint(num)
	return u, err
}

func UintToString(i uint) string {
	return strconv.FormatUint(uint64(i), 10)
}

func HashPassword(pw string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(bytes), err
}

func IsValidVersion(v string) bool {
	if len(v) > MaxVersionLength {
		return false
	}
	if len(v) < MinVersionLength {
		return false
	}
	v = strings.TrimSpace(v)
	if v[0] == 'v' {
		return semver.IsValid(v)
	}
	return semver.IsValid("v" + v)
}

func GetRepoName(input string) string {
	return slug.Make(input)
}

func GetRandomNumber(max int) int {
	return rand.IntN(max) // 0 ≤ n < max
}

func GetRandomUintNumber(max uint) uint {
	return rand.UintN(max) // 0 ≤ n < max
}

func GetRandomNumberInRange(min, max int) int {
	if min > max {
		return -1
	}
	return min + GetRandomNumber(max-min+1)
}

func GetRandomUintRange(min, max uint) uint {
	if min > max {
		return 0
	}
	return min + GetRandomUintNumber(max-min+1)
}

func IsValidName(name string) bool {
	if len(name) < 1 || len(name) > 50 {
		return false
	}
	return nameRegex.MatchString(name)
}

func IsValidHandle(handle string) bool {
	if len(handle) < 3 || len(handle) > 30 {
		return false
	}
	return handleRegex.MatchString(handle)
}

func IsValidModuleName(moduleName string) bool {
	if len(moduleName) < ModuleReprMin || len(moduleName) > ModuleReprMax {
		return false
	}
	return reprRegex.MatchString(moduleName)
}

func ParseModuleReleaseVersion(s string) (string, string, error) {
	parts := strings.Split(s, "@")

	switch len(parts) {
	case 1:
		// "X"
		if parts[0] == "" {
			return "", "", ErrInvalidModuleReleaseVersion
		}
		return parts[0], "", nil

	case 2:
		// "X@Y"
		module := parts[0]
		release := parts[1]

		if module == "" || release == "" {
			return "", "", ErrInvalidModuleReleaseVersion
		}

		if release == "latest" {
			return module, "latest", nil
		}

		// Normalize semver: must start with 'v'
		if !strings.HasPrefix(release, "v") {
			release = "v" + release
		}

		if !semver.IsValid(release) {
			return "", "", ErrInvalidModuleReleaseVersion
		}

		return module, release, nil

	default:
		// "X@Y@Z..."
		return "", "", ErrInvalidModuleReleaseVersion
	}
}

var commonPasswords = map[string]struct{}{
	"password":  {},
	"123456":    {},
	"123456789": {},
	"qwerty":    {},
	"letmein":   {},
	"admin":     {},
}

func ValidatePassword(password, email string) error {
	if PasswordQualityCheckDisabled {
		return nil
	}
	pw := strings.TrimSpace(password)

	length := utf8.RuneCountInString(pw)
	if length < PasswordLengthMin {
		return fmt.Errorf("password must be at least %d characters long", PasswordLengthMin)
	}
	if length > PasswordLengthMax {
		return errors.New("password is too long")
	}

	lowerPw := strings.ToLower(pw)
	if _, ok := commonPasswords[lowerPw]; ok {
		return errors.New("password is too common")
	}

	if email != "" {
		user := strings.ToLower(strings.Split(email, "@")[0])
		if user != "" && strings.Contains(lowerPw, user) {
			return errors.New("password must not contain your email or username")
		}
	}

	return nil
}

const (
	// Max number of terms per field
	MaxTags         = 10
	MaxNames        = 5
	MaxDescriptions = 5

	// Length limits
	TagMinLen = 1
	TagMaxLen = 32

	NameMinLen = 1
	NameMaxLen = 64

	DescMinLen = 1
	DescMaxLen = 128
)

var (
	// tags are identifiers
	tagRe = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

	// name / description search terms (human text)
	textRe = regexp.MustCompile(`^[\p{L}\p{N} ._-]+$`)
)

func validateList(
	label string,
	values []string,
	maxItems int,
	minLen int,
	maxLen int,
	re *regexp.Regexp,
) error {
	// Empty slice is valid
	if len(values) == 0 {
		return nil
	}

	if len(values) > maxItems {
		return fmt.Errorf(
			"%s has too many values (max %d)",
			label,
			maxItems,
		)
	}

	for i, v := range values {
		if len(v) < minLen || len(v) > maxLen {
			return fmt.Errorf(
				"%s[%d] length must be between %d and %d",
				label,
				i,
				minLen,
				maxLen,
			)
		}

		// Reject control characters
		for _, r := range v {
			if r < 32 {
				return fmt.Errorf(
					"%s[%d] contains invalid control characters",
					label,
					i,
				)
			}
		}

		if !re.MatchString(v) {
			return fmt.Errorf(
				"%s[%d] contains invalid characters",
				label,
				i,
			)
		}
	}

	return nil
}

func ValidateTags(tags []string) error {
	return validateList(
		"tags",
		tags,
		MaxTags,
		TagMinLen,
		TagMaxLen,
		tagRe,
	)
}

func ValidateNames(names []string) error {
	return validateList(
		"name",
		names,
		MaxNames,
		NameMinLen,
		NameMaxLen,
		textRe,
	)
}

func ValidateDescriptions(desc []string) error {
	return validateList(
		"description",
		desc,
		MaxDescriptions,
		DescMinLen,
		DescMaxLen,
		textRe,
	)
}
func UniqueLowercaseTokens(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))

	for _, s := range in {
		v := strings.ToLower(s)
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}

	return out
}
