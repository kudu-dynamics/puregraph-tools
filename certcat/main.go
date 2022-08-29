package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
)

func getItems(bundlepath string) ([]string, error) {
	var items []string

	file, err := os.Open(bundlepath)
	if err != nil {
		fmt.Println("Failed to extract from file", err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buf := bytes.Buffer{}
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			buf.WriteString(line + "\n")
		}

		if strings.Contains(line, "----END ") {
			items = append(items, strings.Trim(buf.String(), "\n"))
			buf.Reset()
		}
	}
	return items, nil
}

func getKeyPair(items []string) (int, int, error) {
	// Returns the index of the certificate, key, and nil or an error.
	for i, item1 := range items {
		for j, item2 := range items {
			if i == j {
				continue
			}
			_, err := tls.X509KeyPair([]byte(item1), []byte(item2))
			if err == nil {
				return i, j, nil
			}
		}
	}
	return -1, -1, errors.New("No KeyPair found")
}

func doExtract(bundlepath, outdir string) int {
	err := os.MkdirAll(outdir, 0755)
	if err != nil {
		fmt.Println("Failed to create outdir", err)
		return 1
	}

	items, err := getItems(bundlepath)
	if err != nil {
		fmt.Println("Failed to extract", err)
		return 1
	}
	ci, ki, err := getKeyPair(items)
	if err != nil {
		fmt.Println("Failed to extract", err)
		return 1
	}
	// If the cert, key pair has been identified, write them to files.
	cert, key := items[ci], items[ki]
	ioutil.WriteFile(outdir+"/service.crt", []byte(cert+"\n"), 0644)
	ioutil.WriteFile(outdir+"/service.key", []byte(key+"\n"), 0600)
	// Gather the remaining entries and emit them as a CA file.
	buf := bytes.Buffer{}
	for i, item := range items {
		if i == ci || i == ki {
			continue
		}
		buf.WriteString(item + "\n")
	}
	ioutil.WriteFile(outdir+"/CA.crt", buf.Bytes(), 0644)
	return 0
}

func main() {
	usage := `CertCat.

Usage:
  certcat <bundlepath>
  certcat <bundlepath> <outdir>

Options:
  -h --help	   Show this screen.`

	arguments, _ := docopt.ParseDoc(usage)

	// Parse the incoming arguments.
	bundlepath, _ := arguments.String("<bundlepath>")

	var outdir string
	if val, ok := arguments["<outdir>"].(string); ok {
		outdir = val
	} else {
		outdir = "."
	}

	// Perform Action.
	os.Exit(doExtract(bundlepath, outdir))
}
