package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fd/k8s-utils/internal/generate"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		sourceFileName string
	)

	app := kingpin.New("k8s-gen-template", "Generate content from a template")
	app.Arg("template", "Template file").Default(".").ExistingFileVar(&sourceFileName)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	sourceFileName, err := filepath.Abs(sourceFileName)
	if err != nil {
		app.FatalIfError(err, "%v", err)
	}

	data, err := ioutil.ReadFile(sourceFileName)
	if err != nil {
		app.FatalIfError(err, "%v", err)
	}

	data, err = generate.Generate(sourceFileName, data)
	if err != nil {
		app.FatalIfError(err, "%v", err)
	}

	_, err = io.Copy(os.Stdout, bytes.NewReader(data))
	if err != nil {
		app.FatalIfError(err, "%v", err)
	}
}
