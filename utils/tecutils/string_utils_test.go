package tecutils

import (
	"testing"
)

func TestSplit(t *testing.T) {
	t.Log("try to split with empty values")
	values := "1,2,,,4,5,4,2,1"
	res := Split(values, ",")
	if len(res) != 7 {
		t.Error()
	}

	t.Log("try to split unique values")
	res = SplitUnique(values, ",")
	if len(res) != 4 {
		t.Error()
	}
}
