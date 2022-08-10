package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
)

func main() {
	fileName := flag.String("file", "", "File to Be Worked On. Either HTML file or Markdown")
	flag.Parse()

	if *fileName == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*fileName, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filename string, w io.Writer) error {
	ext := filepath.Ext(filename)
	inputFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	var data []byte
	var tempFileStr string
	switch ext {
	case ".html":
		data, err = parseToMd(inputFile)
		if err != nil {
			return err
		}
		tempFile, err := ioutil.TempFile("", "mdp*.md")
		if err != nil {
			return err
		}
		tempFileStr = tempFile.Name()
		if err != nil {
			return err
		}
	case ".md":
		data, err = parseToHtml(inputFile)
		if err != nil {
			return err
		}
		tempFile, err := ioutil.TempFile("", "mdp*.html")
		if err != nil {
			return err
		}
		tempFileStr = tempFile.Name()
	default:
		return fmt.Errorf("invalid input file")
	}
	fmt.Fprintln(w, tempFileStr)
	saveFile(tempFileStr, data)
	return preview(tempFileStr)
}

func parseToHtml(input []byte) ([]byte, error) {
	var output bytes.Buffer
	err := goldmark.Convert(input, &output)
	if err != nil {
		return []byte{}, err
	}
	body := bluemonday.UGCPolicy().SanitizeBytes(output.Bytes())
	var buffer bytes.Buffer
	buffer.Write(body)
	return buffer.Bytes(), nil
}

func parseToMd(input []byte) ([]byte, error) {
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertBytes(input)
	if err != nil {
		return []byte{}, err
	}
	return markdown, nil
}

func saveFile(outname string, data []byte) error {
	return ioutil.WriteFile("name", data, 0446)
}

func preview(fname string) error {
	cName := ""
	cParams := []string{}

	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("OS not Supported")
	}

	cParams = append(cParams, fname)
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}
	err = exec.Command(cPath, cParams...).Run()
	time.Sleep(2 * time.Second)
	return err
}
