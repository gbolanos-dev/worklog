package update

import (
	"strconv"
	"strings"
)

func newer(current string, latest string) bool {
	parse := func(v string) (int, int, int) {
		p := strings.SplitN(strings.TrimPrefix(v, "v"), ".", 3)
		major, _ := strconv.Atoi(p[0])
		minor, _ := strconv.Atoi(p[1])
		patch, _ := strconv.Atoi(p[2])

		return major, minor, patch
	}

	cMajor, cMinor, cPatch := parse(current)
	lMajor, lMinor, lPatch := parse(latest)

	if lMajor != cMajor {
		return lMajor > cMajor
	}

	if lMinor != cMinor {
		return lMinor > cMinor
	}

	return lPatch > cPatch
}
