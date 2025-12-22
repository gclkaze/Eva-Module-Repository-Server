package utils

import (
	"math/rand/v2"
	"strconv"

	"github.com/gosimple/slug"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/mod/semver"
)

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
