package concurrency

import (
	"strings"
)

func ProcessLineWorker(id int, jobs <-chan string, results chan<- string, indexCols []int, csvDelimiter string) {
	for line := range jobs {
		if len(line) == 0 {
			results <- ""
			continue
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
		results <- strings.Join(nl, csvDelimiter)
	}
}
