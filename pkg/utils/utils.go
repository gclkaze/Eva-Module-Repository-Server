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
	"math/rand/v2"
	"regexp"
	"strconv"

	"github.com/gosimple/slug"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/mod/semver"
)

var nameRegex = regexp.MustCompile(`^[\p{L}]+([\p{L}\s'-]*[\p{L}]+)?$`)
var handleRegex = regexp.MustCompile(
	`^[a-zA-Z][a-zA-Z0-9]*(?:[._-][a-zA-Z0-9]+)*$`,
)
var reprRegex = regexp.MustCompile(`^[A-Za-z0-9]+( [A-Za-z0-9]+)*$`)

var ModuleReprMin = 3
var ModuleReprMax = 50

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
