# Csvcat
`csvcat` is a very fast csv files compiler (with filtering) written in Go. Using concurrency, `csvcat` can concat and filter 
a huge number of files without loosing to much in terms of memeroy and processing time.

Using a dummy dataset generated from `generate_set.py` with `100` files that have each `100000` lines (around 1.4G) in concurrent and non-concurrent
modes to filter 5 columns out of 10, cvscat take around (on an i7-7820HQ (8)):
```sh 
$ ./csvcat --columns "B,A,E,C,F" --delimiter "," --directory "csvset" --concurrency=true
Number of files found: 100
============ Total 3.933102857s ===================
$ ./csvcat --columns "B,A,E,C,F" --delimiter "," --directory "csvset" --concurrency=false
Number of files found: 100
============ Total 10.588019261s ===================
```

## Usage of `csvcat`
```
Usage of ./csvcat:
  -batch int
    	Batch size (default 30)
  -c	Set to false to ignore checking extension (default true)
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

`csvcat` expects every csv file to have a header in its first line where all the columns are labled so that
it can filter the correct columns. If the csv file is not correctly formated (some lines have more/less columns),
it will try to add an empty column in the correct location.


## Building `csvcat`
To build `csvcat` you need to run:
```sh
go build .
// or
go install .
```
