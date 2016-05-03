package main

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"text/template"

	"github.com/fd/k8s-utils/internal/generate"
	"gopkg.in/alecthomas/kingpin.v2"
)

func genSecretMap(app *kingpin.Application, secretDir, secretMapName string) {
	secretDir, err := filepath.Abs(secretDir)
	if err != nil {
		app.FatalIfError(err, "%v", err)
	}

	if secretMapName == "" {
		secretMapName = filepath.Base(secretDir)
	}

	var secretMap = SecretMap{
		Name: secretMapName,
	}

	fis, err := ioutil.ReadDir(secretDir)
	if err != nil {
		app.FatalIfError(err, "%v", err)
	}
	for _, fi := range fis {
		if !fi.Mode().IsRegular() {
			continue
		}

		path := filepath.Join(secretDir, fi.Name())
		name := filepath.Base(fi.Name())
		// ignore system files: .DS_Store
		if name == ".DS_Store" {
			continue
		}
		if len(name) > 253 || !nameSecretMapRE.MatchString(name) {
			app.Fatalf("%q must have at most 253 characters and match regex %s", name, nameSecretMapRE.String())
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			app.FatalIfError(err, "%v", err)
		}

		data, err = generate.Generate(path, data)
		if err != nil {
			app.FatalIfError(err, "%v", err)
		}

		secretMap.Files = append(secretMap.Files, SecretMapFile{
			Name:  name,
			Value: string(data),
		})
	}

	err = secretMap.writeTo(os.Stdout)
	if err != nil {
		app.FatalIfError(err, "%v", err)
	}
}

type SecretMap struct {
	Name  string
	Files []SecretMapFile
}

type SecretMapFile struct {
	Name  string
	Value string
}

func (c *SecretMap) writeTo(w io.Writer) error {
	return secretMapTmpl.Execute(w, c)
}

var nameSecretMapRE = regexp.MustCompile(`^\.?[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`)

// kind: SecretMap
// apiVersion: v1
// metadata:
//   name: example-config
// data:
//   example.property.1: hello
//   example.property.2: world
//   example.property.file: |-
//     property.1=value-1
//     property.2=value-2
//     property.3=value-3
var secretMapTmpl = template.Must(template.New("").Funcs(template.FuncMap{
	"b64": b64,
}).Parse(`kind: Secret
apiVersion: v1
metadata:
  name: {{ .Name }}
data:{{ range .Files }}
  {{ .Name }}: {{ b64 .Value }}{{ end }}
`))

func b64(s string) string {
	return strconv.Quote(base64.StdEncoding.EncodeToString([]byte(s)))
}
