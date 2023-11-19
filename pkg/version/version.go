package version

import "fmt"

// EDIT VERSION HERE
const (
	MAJOR int = 0
	MINOR int = 1
	PATCH int = 0
)

type Version struct {
	Major int
	Minor int
	Patch int
	Tag   string
}

func (v *Version) String() string {
	if v.Tag != "" {
		return fmt.Sprintf("v%d.%d.%d-%s", v.Major, v.Minor, v.Patch, v.Tag)
	}
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func GetVersion() Version {
	return Version{
		Major: MAJOR,
		Minor: MINOR,
		Patch: PATCH,
		Tag:   "alpha",
	}
}

func Compare(v1, v2 Version) int {
	if v1.Major < v2.Major {
		return -1
	} else if v1.Major > v2.Major {
		return 1
	} else if v1.Minor < v2.Minor {
		return -1
	} else if v1.Minor > v2.Minor {
		return 1
	} else if v1.Patch < v2.Patch {
		return -1
	} else if v1.Patch > v2.Patch {
		return 1
	}
	return 0
}
