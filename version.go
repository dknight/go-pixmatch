package pixmatch

import (
	"fmt"
	"strings"
)

// Version is extremly simple semantic version structure.
// Me personally wouldn't like to use external (especially heavy) dependices
// to manage versions.
//
// For monstonous version management you can check this one
// https://pkg.go.dev/github.com/hashicorp/go-version
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
