# Csvcat
**csvcat** is a very fast csv files compiler (with filtering) written in Go. Using concurrency, **csvcat** can concat and filter 
a huge number of files without loosing to much in terms of memeroy and processing time.

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
