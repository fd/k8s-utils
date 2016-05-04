package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"text/template"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gopkg.in/alecthomas/kingpin.v2"

	raw "google.golang.org/api/container/v1"
)

func gcpLogin(app *kingpin.Application, gcpProjectID, gcpZone, gkeCluster string) {
	client, err := google.DefaultClient(oauth2.NoContext)
	app.FatalIfError(err, "k8s")

	svc, err := raw.New(client)
	app.FatalIfError(err, "k8s")

	cluster, err := svc.Projects.Zones.Clusters.Get(gcpProjectID, gcpZone, gkeCluster).Do()
	app.FatalIfError(err, "k8s")

	var buf bytes.Buffer
	err = k8sConfigTmpl.Execute(&buf, cluster)
	app.FatalIfError(err, "k8s")

	err = os.MkdirAll(path.Join(os.Getenv("HOME"), ".kube"), 0755)
	app.FatalIfError(err, "k8s")

	err = ioutil.WriteFile(path.Join(os.Getenv("HOME"), ".kube/config"), buf.Bytes(), 0600)
	app.FatalIfError(err, "k8s")
}

var k8sConfigTmpl = template.Must(template.New("").Parse(`
apiVersion: v1
kind: Config
preferences: {}
current-context: k8s-context
contexts:
- context:
    cluster: k8s-cluster
    user:    k8s-user
  name: k8s-context
clusters:
- cluster:
    certificate-authority-data: {{.MasterAuth.ClusterCaCertificate}}
    server:                     https://{{.Endpoint}}
  name: k8s-cluster
users:
- name: k8s-user
  user:
    client-certificate-data: {{.MasterAuth.ClientCertificate}}
    client-key-data:         {{.MasterAuth.ClientKey}}
    username: {{.MasterAuth.Username}}
    password: {{.MasterAuth.Password}}
`))
