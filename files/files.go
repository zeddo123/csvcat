package files

import (
	"os"
	"time"
	"path/filepath"
	"log"
	"strings"
	"fmt"
)

func AsyncCompileFiles(directory string, csvDelimiter string, output string, csvColumns string, batchSize int, checkFileExtension bool) {
	var files = getfilenames(directory, checkFileExtension)

	targetCols := strings.Split(csvColumns, csvDelimiter)

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
				processFileAsync(file, results, targetCols, csvDelimiter)
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

func CompileFiles(directory string, csvDelimiter string, output string, csvColumns string, checkFileExtension bool) {
	// get files in the directory
	var files = getfilenames(directory, checkFileExtension)

	targetCols := strings.Split(csvColumns, csvDelimiter)

	fmt.Println("Number of files found:", len(files))

	of, err := os.Create(output)
	defer of.Close()
	if err != nil {
		log.Fatal(err)
	}
	of.WriteString(csvColumns + "\n")

	for _, file := range files {
		fname := filepath.Base(file)
		log.Println("[JOB] Starting job", fname)
		start := time.Now()
		processedFile := processFileSync(file, targetCols, csvDelimiter)
		log.Println("[JOB] Finished job", fname, "in ------------------", time.Now().Sub(start))
		start = time.Now()
		of.WriteString(processedFile.Content)
		log.Println("[IO] Finished Writing", fname, "✓ ( +", time.Now().Sub(start), ")")
	}
}

func getfilenames(directory string, checkFileExtension bool) []string {
	var files []string

	entries, err := os.ReadDir(directory)
	path, _ := filepath.Abs(directory)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		// check if entry is a file and has correct file type
		if e.Type().IsRegular() {
			if !checkFileExtension || strings.HasSuffix(e.Name(), ".csv") {
				files = append(files, filepath.Join(path, e.Name()))
			}
		}
	}

	return files
}
