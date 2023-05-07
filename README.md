# Csvcat
`csvcat` is a very fast csv files compiler (with filtering) written in Go. Using concurrency, `csvcat` can concat and filter 
a huge number of files without loosing to much in terms of memeroy and processing time.

## Usage of `csvcat`
```
Usage of ./csvcat:
  -batch int
    	Batch size (default 30)
  -columns string
    	Columns to be selected
  -concurrency
    	Set flag to disable concurrency (default true)
  -delimiter string
    	Csv delimiter of files (default ",")
  -directory string
    	Directory containing the files to be compilled (default ".")
  -output string
    	Output filename (default "output.csv")
  -v	Set to true to have verbose output
```
Here's an example of how you might run `csvcat` with its flags:
```
./csvcat --batch 20 --columns "A,B,C" --delimiter "," --directory files
```

`csvcat` expects every csv file to have a header in it's first line where all the columns are labled so that
it can filter the correct columns. If the csv file is not correctly formated (some lines have more/less columns),
it will try to add an empty column in the correct location.


## Building `csvcat`
To build `csvcat` you need to run:
```sh
go build .
// or
go install .
```
