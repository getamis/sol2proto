package util

import "testing"
import "strings"

func TestToUnderScore(t *testing.T) {
	testCases := map[string]string{
		"_requestID": "request_id",
		"requestId":  "request_id",
		"RequestId":  "request_id",
		"requestID":  "request_id",
		"ID":         "id",
		// FIXME: cannot handle below case
		//"requestIDd": "request_i_dd",
		//"requestIDAndTxID": "request_id_and_tx_id",
	}

	for k, exp := range testCases {
		v := ToUnderScore(k)
		if strings.Compare(exp, v) != 0 {
			t.Errorf("Unexpected output, input: %v, got: %v, expected: %v\n", k, v, exp)
		}
	}
}
