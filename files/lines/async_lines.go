package lines

func ProcessLineWorker(id int, jobs <-chan string, results chan<- string, indexCols []int, csvDelimiter string) {
	for line := range jobs {
		results <- ProcessLine(line, indexCols, csvDelimiter)
	}
}
