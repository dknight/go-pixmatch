package pixmatch

import (
	"fmt"
	"strings"
)

// currentVersion of the software, do not forget to change it with release.
var currentVersion = Version{1, 0, 0, ""}

// Version is extremely simple semantic version structure. Personally, I
// wouldn't like to use external (especially heavy) dependencies to manage
// versions.
//
// Check out https://pkg.go.dev/github.com/hashicorp/go-version for monstrous
// version management.
type Version struct {
	Major int
	Minor int
	Patch int
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

// GetVersion gets current version of pixmatch.
func GetVersion() string {
	return currentVersion.String()
}
