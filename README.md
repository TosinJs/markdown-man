# markdown-man

## Description

markdown-man is a cli tool used to convert markdown files (.md) to html files (.html) and vice versa.

## Functionality
markdown-man converts markdown files to HTML and vice versa. <br />
It has a preview toggle which opens the converted file in your browser. <br />
All generated files are saved in the directory of the source file.

## Flags
```bash
$ -file: sets the file to be worked On. Either HTML (.html) file or a Markdown (.md) file.
$ -s: a boolean to determine if the preview functionality should be skipped. The default value is false.
```

## Usage
```bash
# Clone the repository
$ git clone https://github.com/TosinJs/markdown-man.git

# Install dependencies
$ cd markdown-man
$ go get ./...

# Start the application
In the root folder run:
$ go run ./src -file "path to file"
```
## Test

```bash
# unit tests
$ go test
```
