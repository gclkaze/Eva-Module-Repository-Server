// Package utils contains
package utils

import (
	"fmt"
	"strings"
)

func BuildWhereConditionStringForUniqueAttrsContaining(attr string, tags []string) string {
	unique := make(map[string]struct{})
	for _, s := range tags {
		unique[s] = struct{}{}
	}

	var conditions []string
	for s := range unique {
		conditions = append(conditions, fmt.Sprintf("%s LIKE '%%%s%%'", attr, s))
	}

	whereClause := strings.Join(conditions, " OR ")
	return whereClause
}
