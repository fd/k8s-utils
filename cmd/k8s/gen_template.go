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

func genTemplate(app *kingpin.Application, sourceFileName string) {
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
