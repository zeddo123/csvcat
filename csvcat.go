package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zeddo123/csvcat/v2/concurrency"
	//"github.com/pkg/profile"
)


func main() {
	//defer profile.Start(profile.MemProfile).Stop()
	entrypoint();
}


func entrypoint() {
	// Define flags
	directory := flag.String("directory", ".", "Directory containing the files to be compilled")
	delim := flag.String("delimiter", ",", "Csv delimiter of files")
	output_file := flag.String("output", "output.csv", "Output filename")
	csvColumns := flag.String("columns", "", "Columns to be selected")
	batchSize := flag.Int("batch", 30, "Batch size")
	flag.Parse()

	start := time.Now()
	compileFiles(*directory, *delim, *output_file, *csvColumns, *batchSize)
	fmt.Println("============ Total", time.Now().Sub(start), "===================")
}



func process(bytes []byte, targetCols []string, csvDelimiter string) string {
	content := string(bytes)
	if len(content) == 0 {
		return content
	}

	lines := strings.Split(content, "\n")
	firstLine := lines[0]
	lines = lines[1:]
	
	columns := strings.Split(firstLine, "^")
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

	jobchan := make(chan string, len(lines))
	resultschan := make(chan string, len(lines))
	for id := 0; id < int(float64(len(lines)) * 0.01); id++ {
		go concurrency.ProcessLineWorker(id, jobchan, resultschan, indexCols, csvDelimiter)
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

func readFile(fileName string, ch chan<- FileContent, csvColumns []string, csvDelimiter string) int {
	fi, err := os.Open(fileName)
	if err != nil {
		log.Println("Couldn't open ", fileName)
		return 0
	}

	bytes, err := io.ReadAll(fi)
	if err != nil {
		log.Println("Couldn't read ", fileName)
		return  0
	}

	start := time.Now()
	content := process(bytes, csvColumns, csvDelimiter)
	log.Println("[PROCESS] Finished processing", filepath.Base(fileName), "in", time.Now().Sub(start))
	ch <- FileContent{fileName, content}

	return 1
}

func compileFiles(directory string, csvDelimiter string, output string, csvColumns string, batchSize int) {
	var files []string

	entries, err := os.ReadDir(directory)
	path, _ := filepath.Abs(directory)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		// check if entry is a file and has correct file type
		if e.Type().IsRegular() && strings.HasSuffix(e.Name(), ".csv") {
			files = append(files, filepath.Join(path, e.Name()))
		}
	}

	targetCols := strings.Split(csvColumns, "^")

	jobschan := make(chan string, len(files))
	ch := make(chan FileContent, len(files))
	fmt.Println("Number of files found:", len(files))

	// Start workers
	for i := 0; i < min(batchSize, len(files)); i++ {
		go func (id int, jobs <-chan string, results chan<- FileContent, targetCols []string){
			done := 0
			for file := range jobs {
				fname := filepath.Base(file)
				log.Println("[JOB] Starting job", fname)
				start := time.Now()
				readFile(file, results, targetCols, csvDelimiter)
				log.Println("[JOB] Finished job", fname, "in ------------------", time.Now().Sub(start))
				done++
			}
			log.Println("[Worker", id, "] Finished", done, "job(s)")
		}(i, jobschan, ch, targetCols)
	}

	for _, file := range files {
		jobschan <- file
	}
	close(jobschan)

	of, err := os.Create(output)
	defer of.Close()
	if err != nil {
		log.Fatal(err)
	}
	of.WriteString(csvColumns + "\n")
	for i:=0; i < len(files); i++ {
		start := time.Now()
		filecontent := <- ch
		of.WriteString(filecontent.Content)
		log.Println("[IO] Finished Writing", filepath.Base(filecontent.Name), "✓ ( +", time.Now().Sub(start), ")")
		//of.Sync()
	}
}
