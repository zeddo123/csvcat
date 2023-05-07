package files

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	Line "github.com/zeddo123/csvcat/v2/files/lines"
)


func processFileAsync(fileName string, ch chan<- FileContent, csvColumns []string, csvDelimiter string) int {
	bytes, err := readFile(fileName)
	if err != nil {
		return  0
	}

	start := time.Now()
	content := filterColumnsAsync(bytes, csvColumns, csvDelimiter)
	log.Println("[PROCESS] Finished processing", filepath.Base(fileName), "in", time.Now().Sub(start))
	ch <- FileContent{fileName, content}

	return 1
}

func filterColumnsAsync(bytes []byte, targetCols []string, csvDelimiter string) string {
	content := string(bytes)
	if len(content) == 0 {
		return content
	}

	lines := strings.Split(content, "\n")
	firstLine := lines[0]
	lines = lines[1:]
	
	indexCols := getColumnIndices(firstLine, targetCols, csvDelimiter)
	
	jobchan := make(chan string, len(lines))
	resultschan := make(chan string, len(lines))
	for id := 0; id < int(float64(len(lines)) * 0.01); id++ {
		go Line.ProcessLineWorker(id, jobchan, resultschan, indexCols, csvDelimiter)
	}


	for _, line := range lines {
		jobchan <- line
	}
	close(jobchan)

	var ls = []string{}
	for i := 0; i < len(lines); i++ {
		ls = append(ls, <-resultschan)
	}
	return strings.Join(ls, "\n")
}

func processFileSync(fileName string, csvColumns []string, csvDelimiter string) FileContent {
	bytes, err := readFile(fileName)
	if err != nil {
		return FileContent{fileName, ""}
	}
	
	return FileContent{fileName, filterColumnsSync(bytes, csvColumns, csvDelimiter)}
}

func filterColumnsSync(bytes []byte, targetCols []string, csvDelimiter string) string {
	content := string(bytes)
	if len(content) == 0 {
		return content
	}

	lines := strings.Split(content, "\n")
	firstLine := lines[0]
	lines = lines[1:]
	
	indexCols := getColumnIndices(firstLine, targetCols, csvDelimiter)
	
	var ls = make([]string, 0, len(lines))
	for _, line := range lines {
		ls = append(ls, Line.ProcessLine(line, indexCols, csvDelimiter))
	}

	return strings.Join(ls, "\n")
}

func readFile(fileName string) ([]byte, error) {
	var bytes []byte
	fi, err := os.Open(fileName)
	if err != nil {
		log.Println("Couldn't open ", fileName)
		return bytes, errors.New("Couldn't open file " + fileName)
	}

	bytes, err = io.ReadAll(fi)
	if err != nil {
		log.Println("Couldn't read ", fileName)
		return bytes, errors.New("Couldn't read file " + fileName)
	}
	return bytes, nil
}

func getColumnIndices(firstLine string, targetCols []string, csvDelimiter string) []int {
	columns := strings.Split(firstLine, csvDelimiter)
	indexCols := make([]int, 0, len(columns))
	for _, col := range targetCols {
		flag := false
		for i, c := range columns {
			if c == col {
				indexCols = append(indexCols, i)
				flag = true
				break
			}
		}
		// if the column was not found we need to still have
		// place for it in the final file.
		// e.g A^B^C -> A^^B^C
		if !flag {
			indexCols = append(indexCols, len(columns))
		}
	}
	return indexCols
}
