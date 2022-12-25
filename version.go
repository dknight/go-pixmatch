package pixmatch

import (
	"fmt"
	"strings"
)

// Version is extremly simple semantic version structure.
// Me personally wouldn't like to use external (especially heavy) dependices
// to manage versions.
type Version struct {
	Major uint
	Minor uint
	Patch uint
	Pre   string
}

func (v Version) String() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Pre != "" {
		builder.WriteString("-")
		builder.WriteString(v.Pre)
	}
	return builder.String()
}
