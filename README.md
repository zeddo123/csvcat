# Csvcat
`csvcat` is a very fast csv files compiler (with filtering) written in Go. Using concurrency, `csvcat` can concat and filter 
a huge number of files without loosing to much in terms of memeroy and processing time.

## Usage of `csvcat`
```
Usage of csvcat:
  -batch int
    	Batch size (default 30)
  -columns string
    	Columns to be selected
  -delimiter string
    	Csv delimiter of files (default ",")
  -directory string
    	Directory containing the files to be compilled (default ".")
  -output string
    	Output filename (default "output.csv")
```
Here's an example of how you might run `csvcat` with its flags.
```
./csvcat --batch 20 --columns "A,B,C" --delimiter "," --directory files
```

## Building `csvcat`
To build `csvcat` you need to run:
```sh
go build .
// or
go install .
```
