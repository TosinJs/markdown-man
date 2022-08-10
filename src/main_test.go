package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

const (
	inputMD   = "../testdata/test.data.md"
	inputHTML = "../testdata/test.data.html"
)

func TestParseToMd(t *testing.T) {
	input, err := ioutil.ReadFile(inputHTML)
	if err != nil {
		t.Fatal(err)
	}
	got, err := parseToMd(input)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := ioutil.ReadFile(inputMD)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Logf("golden:\n%s\n", expected)
		t.Logf("result:\n%s", got)
		t.Error("Result content does not match golden file")
	}
}

func TestParseToHtml(t *testing.T) {
	input, err := ioutil.ReadFile(inputMD)
	if err != nil {
		t.Fatal(err)
	}
	got, err := parseToHtml(input)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := ioutil.ReadFile(inputHTML)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Logf("golden:\n%s\n", expected)
		t.Logf("result:\n%s", got)
		t.Error("Result content does not match golden file")
	}
}

func TestRun(t *testing.T) {
	tt := []struct {
		name       string
		inputFile  string
		outputFile string
	}{
		{
			name:       "Convert From Markdown to HTML",
			inputFile:  inputMD,
			outputFile: inputHTML,
		},
		{
			name:       "Convert From HTML to Markdown",
			inputFile:  inputHTML,
			outputFile: inputMD,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var mockStdOut bytes.Buffer
			if err := run(tc.inputFile, &mockStdOut); err != nil {
				t.Fatal(err)
			}
			resultFile := strings.TrimSpace(mockStdOut.String())
			result, err := ioutil.ReadFile(resultFile)
			if err != nil {
				t.Fatal(err)
			}
			expected, err := ioutil.ReadFile(tc.outputFile)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(expected, result) {
				t.Logf("golden:\n%s\n", expected)
				t.Logf("result:\n%s\n", result)
				t.Error("Result content does not match golden file")
			}
			os.Remove(resultFile)
		})
	}
}
