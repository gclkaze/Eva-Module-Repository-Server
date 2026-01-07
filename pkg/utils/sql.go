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
		conditions = append(conditions, fmt.Sprintf("%s LIKE '%s%%'", attr, s))
		conditions = append(conditions, fmt.Sprintf("%s LIKE '%%%s'", attr, s))
	}

	whereClause := strings.Join(conditions, " OR ")
	return whereClause
}

func BuildWhereConditionStringForUniqueAttrContaining(attr string, tag string) string {

	var conditions []string
	conditions = append(conditions, fmt.Sprintf("%s LIKE '%%%s%%'", attr, tag))
	conditions = append(conditions, fmt.Sprintf("%s LIKE '%s%%'", attr, tag))
	conditions = append(conditions, fmt.Sprintf("%s LIKE '%%%s'", attr, tag))

	whereClause := strings.Join(conditions, " OR ")
	return whereClause
}
