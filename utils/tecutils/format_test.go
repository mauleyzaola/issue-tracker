package tecutils

import (
	"testing"
	"time"
)

func TestDateFormat(t *testing.T) {
	value := time.Date(1968, 9, 15, 21, 15, 0, 0, time.Local)
	expected := "15-09-1968"
	format := FormatDate(&value)
	if format != expected {
		t.Errorf("expected: %s but got: %s", expected, format)
		return
	}

	expected = "15-09-1968 21:15"
	format = FormatDateTime(&value)
	if format != expected {
		t.Errorf("expected: %s but got: %s", expected, format)
		return
	}
}
