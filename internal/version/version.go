package version

import "time"

var (
	Development = "devel"
	Version     = Development
	Commit      = "none"
	Date        = time.Now().String()
	BuiltBy     = "unknown"
)
