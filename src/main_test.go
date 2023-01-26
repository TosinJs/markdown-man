package main

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

const (
	inputMD   = "../testdata/test.data.md"
	inputHTML = "../testdata/test.data.html"
)

func cleanString(s string) []byte {
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	return []byte(s)
}

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

	got = cleanString(string(got))
	expected = cleanString(string(expected))

	if !bytes.Equal(got, expected) {
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

	got = cleanString(string(got))
	expected = cleanString(string(expected))

	if !bytes.Equal(got, expected) {
		t.Logf("golden:\n%s\n", expected)
		t.Logf("result:\n%s", got)
		t.Error("Result content does not match golden file")
	}
}
