package utils

import (
	"strconv"

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
