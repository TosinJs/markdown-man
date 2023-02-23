package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/browser"
	"github.com/yuin/goldmark"
)

const (
	header = `
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta http-equiv="X-UA-Compatible" content="IE=edge">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Document</title>
		</head>
		<body>`

	footer = `
		</body>
		</html>`
)

func main() {
	fileName := flag.String("file", "", "File to Be Worked On. Either HTML file or Markdown")
	skipPreview := flag.Bool("s", false, "Skip File Preview")
	flag.Parse()

	if *fileName == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*fileName, *skipPreview, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filename string, skipPreview bool, w io.Writer) error {
	ext := filepath.Ext(filename)
	inputFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	str, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	var data []byte
	var outputFileName = strings.TrimSuffix(str, ext)
	rand.Seed(time.Now().UnixNano())
	switch ext {
	case ".html":
		data, err = parseToMd(inputFile)
		if err != nil {
			return err
		}
		outputFileName = fmt.Sprintf("%s%v.mm.md", outputFileName, rand.Intn(100))
	case ".md":
		data, err = parseToHtml(inputFile)
		if err != nil {
			return err
		}
		outputFileName = fmt.Sprintf("%s%v.mm.html", outputFileName, rand.Intn(100))
	default:
		return fmt.Errorf("invalid input file")
	}

	err = saveFile(outputFileName, data)
	if err != nil {
		return err
	}
	if skipPreview {
		return nil
	}
	return browser.OpenFile(outputFileName)
}

func parseToHtml(input []byte) ([]byte, error) {
	var output bytes.Buffer
	err := goldmark.Convert(input, &output)
	if err != nil {
		return []byte{}, err
	}
	body := bluemonday.UGCPolicy().SanitizeBytes(output.Bytes())
	var buffer bytes.Buffer
	buffer.WriteString(header)
	buffer.Write(body)
	buffer.WriteString(footer)
	return buffer.Bytes(), nil
}

func parseToMd(input []byte) ([]byte, error) {
	converter := md.NewConverter("", true, nil).Remove("title")
	markdown, err := converter.ConvertBytes(input)
	if err != nil {
		return []byte{}, err
	}
	return markdown, nil
}

func saveFile(outname string, data []byte) error {
	return ioutil.WriteFile(outname, data, 0446)
}
