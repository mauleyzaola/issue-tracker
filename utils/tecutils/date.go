package tecutils

import (
	"time"
)

func truncMs(d *time.Time, loc *time.Location) {
	if d == nil {
		return
	}
	*d = time.Date(d.Year(), d.Month(), d.Day(), d.Hour(), d.Minute(), d.Second(), 0, loc)
}

func TruncMs(d *time.Time, loc *time.Location) {
	truncMs(d, loc)
}

func TruncMsUTC(d *time.Time) {
	truncMs(d, time.UTC)
}
