package internal

import (
	"fmt"
	"strconv"
	"time"
)

var (
	commitVersion = "v0.2"       //nolint:gochecknoglobals // May be updated using build flags
	commitDate    = "1731409630" //nolint:gochecknoglobals // commitDate in Epoch seconds (may be overridden using build flags)
)

// GetVersion - get version and also commitHash and commitDate if inserted via Makefile.
func GetVersion() string {
	seconds, _ := strconv.Atoi(commitDate)
	if commitDate != "" {
		t := time.Unix(int64(seconds), 0)
		return fmt.Sprintf("%s, date: %s", commitVersion, t.Format("2006-01-02"))
	}
	return commitVersion
}
