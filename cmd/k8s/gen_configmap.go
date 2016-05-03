package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/fd/k8s-utils/internal/generate"
	"gopkg.in/alecthomas/kingpin.v2"
)

func genConfigMap(app *kingpin.Application, configDir, configMapName string) {
	configDir, err := filepath.Abs(configDir)
	if err != nil {
		app.FatalIfError(err, "%v", err)
	}

	if configMapName == "" {
		configMapName = filepath.Base(configDir)
	}

	var configMap = ConfigMap{
		Name: configMapName,
	}

	fis, err := ioutil.ReadDir(configDir)
	if err != nil {
		app.FatalIfError(err, "%v", err)
	}
	for _, fi := range fis {
		if !fi.Mode().IsRegular() {
			continue
		}

		path := filepath.Join(configDir, fi.Name())
		name := filepath.Base(fi.Name())
		// ignore system files: .DS_Store
		if name == ".DS_Store" {
			continue
		}
		if len(name) > 253 || !nameConfigMapRE.MatchString(name) {
			app.Fatalf("%q must have at most 253 characters and match regex %s", name, nameConfigMapRE.String())
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			app.FatalIfError(err, "%v", err)
		}

		data, err = generate.Generate(path, data)
		if err != nil {
			app.FatalIfError(err, "%v", err)
		}

		configMap.Files = append(configMap.Files, ConfigMapFile{
			Name:  name,
			Value: string(data),
		})
	}

	err = configMap.writeTo(os.Stdout)
	if err != nil {
		app.FatalIfError(err, "%v", err)
	}
}

type ConfigMap struct {
	Name  string
	Files []ConfigMapFile
}

type ConfigMapFile struct {
	Name  string
	Value string
}

func (c *ConfigMap) writeTo(w io.Writer) error {
	return configMapTmpl.Execute(w, c)
}

var nameConfigMapRE = regexp.MustCompile(`^\.?[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`)

// kind: ConfigMap
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
var configMapTmpl = template.Must(template.New("").Funcs(template.FuncMap{
	"indent": indent,
}).Parse(`kind: ConfigMap
apiVersion: v1
metadata:
  name: {{ .Name }}
data:{{ range .Files }}
  {{ .Name }}: {{ indent .Value }}{{ end }}
`))

func indent(s string) string {
	s = strings.TrimSpace(s)
	if strings.IndexByte(s, '\n') >= 0 {
		s = "|-\n" + s
		s = strings.Replace(s, "\n", "\n     ", -1)
	} else {
		s = strconv.Quote(s)
	}
	return s
}
