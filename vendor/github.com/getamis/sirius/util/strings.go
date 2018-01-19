// Copyright 2017 AMIS Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"regexp"
	"strings"
)

// ToCamelCase converts an under-score string to a camel-case string
func ToCamelCase(inputUnderScoreStr string) (camelCase string) {
	isToUpper := false

	for k, v := range inputUnderScoreStr {
		if k == 0 {
			camelCase = strings.ToUpper(string(inputUnderScoreStr[0]))
		} else {
			if isToUpper {
				camelCase += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					camelCase += string(v)
				}
			}
		}
	}
	return
}

// ToUnderScore converts a camel-case string to a under-score string
func ToUnderScore(s string) string {
	return SplitAndJoin(s, "_")
}

// ToHyphen converts a camel-case string to a hyphen-style string
func ToHyphen(s string) string {
	return SplitAndJoin(s, "-")
}

// SplitAndJoin converts a camel-case string to a string joined by the provided symbol
func SplitAndJoin(s string, symbol string) string {
	var camel = regexp.MustCompile("(^[^A-Z0-9]*)?([A-Z0-9]{2,}|[A-Z0-9][^A-Z]+|$)")
	var a []string
	for _, sub := range camel.FindAllStringSubmatch(s, -1) {
		if sub[1] != "" {
			a = append(a, sub[1])
		}
		if sub[2] != "" {
			a = append(a, sub[2])
		}
	}

	result := strings.ToLower(strings.Join(a, symbol))
	result = strings.TrimPrefix(result, "_")
	result = strings.TrimSuffix(result, "_")
	return result
}
