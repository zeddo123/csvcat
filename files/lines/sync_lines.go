package lines

import (
	"strings"
)

func ProcessLine(line string, indexCols []int, csvDelimiter string) string {
	if len(line) == 0 {
		return ""
	}

	nl := []string{}
	cols := strings.Split(line, csvDelimiter)
	numbercol := len(cols)
	for _, index := range indexCols {
		if numbercol > index{
			nl = append(nl, cols[index])
		} else {
			nl = append(nl, "")
		}
	}
	return strings.Join(nl, csvDelimiter)
}
