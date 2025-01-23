package defs

import "fmt"

var Version = ""
var Revision = ""

func VersionString() string {
	return fmt.Sprintf("tskks v%s+%s", Version, Revision)
}
