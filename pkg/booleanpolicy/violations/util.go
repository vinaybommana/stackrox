package violations

import (
	"fmt"
	"sort"
	"strings"
)

func stringSliceToSortedSentence(s []string) string {
	var sb strings.Builder
	sort.Strings(s)
	switch sLen := len(s); {
	case sLen == 1:
		fmt.Fprintf(&sb, "%s", s[0])
	case sLen == 2:
		fmt.Fprintf(&sb, "%s and %s", s[0], s[1])
	default:
		for idx, elem := range s {
			if idx < sLen-1 {
				fmt.Fprintf(&sb, "%s, ", elem)
			} else {
				fmt.Fprintf(&sb, "and %s", elem)
			}
		}
	}
	return sb.String()
}
