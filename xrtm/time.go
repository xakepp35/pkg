package xrtm

import (
	"time"
)

// timestampFunc for UTC time
func TimeUTC() time.Time {
	return time.Now().UTC()
}
