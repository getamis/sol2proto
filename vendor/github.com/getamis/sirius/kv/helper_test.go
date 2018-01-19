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

package kv

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKeyMaker(t *testing.T) {
	seperator := ":"
	vals := []struct {
		expectResult string
		maker        keyMaker
		identity     string
		args         []string
	}{
		{
			// full information
			"pkg:tag1:tag1:idendtity1:arg1:arg1",
			NewKey(seperator, "pkg", "tag1", "tag1"),
			"idendtity1",
			[]string{"arg1", "arg1"},
		},
		{
			// missing tags
			"pkg:idendtity2",
			NewKey(seperator, "pkg"),
			"idendtity2",
			nil,
		},
		{
			// missing args
			"pkg:tag3:tag3:idendtity3",
			NewKey(seperator, "pkg", "tag3", "tag3"),
			"idendtity3",
			nil,
		},
		{
			// empty a args
			"pkg:tag4:tag4:idendtity4",
			NewKey(seperator, "pkg", "tag4", "tag4"),
			"idendtity4",
			[]string{},
		},
		{
			// empty identity
			"pkg:",
			NewKey(seperator, "pkg"),
			"",
			[]string{},
		},
		{
			// empty identity with a arg
			"pkg::arg6",
			NewKey(seperator, "pkg"),
			"",
			[]string{"arg6"},
		},
	}
	for _, v := range vals {
		result := v.maker(v.identity, v.args...)
		assert.Equal(t, v.expectResult, result, "should be equal")
	}
}

func ExampleNewKeyMaker() {
	keyMaker := NewKey("/", "package", "tag1", "tag2")
	fmt.Println(keyMaker("golang"))
	// Output: package/tag1/tag2/golang
}
