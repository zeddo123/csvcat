package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/zeddo123/csvcat/v2/files"
	//"github.com/pkg/profile"
)


func main() {
	//defer profile.Start(profile.MemProfile).Stop()
	// Define flags
	directory := flag.String("directory", ".", "Directory containing the files to be compilled")
	delim := flag.String("delimiter", ",", "Csv delimiter of files")
	output_file := flag.String("output", "output.csv", "Output filename")
	csvColumns := flag.String("columns", "", "Columns to be selected")
	batchSize := flag.Int("batch", 30, "Batch size")
	concurrently := flag.Bool("concurrency", true, "Set flag to disable concurrency")
	verbose := flag.Bool("v", false, "Set to true to have verbose output")
	checkExtension := flag.Bool("c", true, "Set to false to ignore checking extension")
	flag.Parse()

	// Disabling the logger
	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}

	start := time.Now()
	if *concurrently {
		files.AsyncCompileFiles(*directory, *delim, *output_file, *csvColumns, *batchSize, *checkExtension)
	} else {
		files.CompileFiles(*directory, *delim, *output_file, *csvColumns, *checkExtension)
	}
	fmt.Println("============ Total", time.Now().Sub(start), "===================")
}

