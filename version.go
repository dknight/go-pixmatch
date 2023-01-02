package pixmatch

import (
	"fmt"
	"strings"
)

// Version is extremly simple semantic version structure. Personally, I
// wouldn't like to use external (especially heavy) dependencies to manage
// versions.
//
// Check out https://pkg.go.dev/github.com/hashicorp/go-version for monstrous
// version management.
type Version struct {
	Major uint
	Minor uint
	Patch uint
	Meta  string
}

func (v Version) String() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Meta != "" {
		builder.WriteString("-")
		builder.WriteString(v.Meta)
	}
	return builder.String()
}

// GetVersion gets current version of pixmatch.
func GetVersion() string {
	return Version{0, 0, 2, "alpha-3"}.String()
}
