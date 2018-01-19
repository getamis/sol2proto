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
	"strings"
	"testing"
)

func TestToCamelCase(t *testing.T) {
	testCases := map[string]string{
		"a_b_c":       "ABC",
		"aa_bb_cc":    "AaBbCc",
		"aaa_bbb_ccc": "AaaBbbCcc",
	}

	for k, exp := range testCases {
		v := ToCamelCase(k)
		if strings.Compare(exp, v) != 0 {
			t.Errorf("Unexpected output, input: %v, got: %v, expected: %v\n", k, v, exp)
		}
	}
}

func TestToUnderScore(t *testing.T) {
	testCases := map[string]string{
		"_requestID": "request_id",
		"requestId":  "request_id",
		"RequestId":  "request_id",
		"requestID":  "request_id",
		"request":    "request",
		"Request":    "request",
		"ID":         "id",
		// FIXME: cannot handle below case
		//"requestIDd":       "request_i_dd",
		//"requestIDAndTxID": "request_id_and_tx_id",
	}

	for k, exp := range testCases {
		v := ToUnderScore(k)
		if strings.Compare(exp, v) != 0 {
			t.Errorf("Unexpected output, input: %v, got: %v, expected: %v\n", k, v, exp)
		}
	}
}

func TestToHyphen(t *testing.T) {
	testCases := map[string]string{
		"_requestID": "request-id",
		"requestId":  "request-id",
		"RequestId":  "request-id",
		"requestID":  "request-id",
		"request":    "request",
		"Request":    "request",
		"ID":         "id",
		// FIXME: cannot handle below case
		//"requestIDd":       "request_i_dd",
		//"requestIDAndTxID": "request_id_and_tx_id",
	}

	for k, exp := range testCases {
		v := ToHyphen(k)
		if strings.Compare(exp, v) != 0 {
			t.Errorf("Unexpected output, input: %v, got: %v, expected: %v\n", k, v, exp)
		}
	}
}

func TestSplitAndJoin(t *testing.T) {
	testCases := map[string]string{
		"_requestID": "request*id",
		"requestId":  "request*id",
		"RequestId":  "request*id",
		"requestID":  "request*id",
		"request":    "request",
		"Request":    "request",
		"ID":         "id",
		// FIXME: cannot handle below case
		//"requestIDd":       "request_i_dd",
		//"requestIDAndTxID": "request_id_and_tx_id",
	}

	for k, exp := range testCases {
		v := SplitAndJoin(k, "*")
		if strings.Compare(exp, v) != 0 {
			t.Errorf("Unexpected output, input: %v, got: %v, expected: %v\n", k, v, exp)
		}
	}
}
