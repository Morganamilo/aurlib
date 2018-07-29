package slice

import (
	"github.com/Morganamilo/aurlib/utils/math"
)

func Longest(strs []string) int {
	biggest := 0

	for _, str := range strs {
		biggest = math.Max(len(str), biggest)
	}

	return biggest
}
