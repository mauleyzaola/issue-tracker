package tecutils

import (
	"testing"
	"time"
)

func TestTruncMs(t *testing.T) {
	var now *time.Time
	now = &time.Time{}
	*now = time.Now()
	t.Log(now.String())
	TruncMs(now, time.Local)
	t.Log(now.String())
	TruncMs(now, time.UTC)
	t.Log(now.String())
}
