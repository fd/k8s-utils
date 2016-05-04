package main

import (
	"os"

	"limbo.services/version"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		configDir      string
		configMapName  string
		secretDir      string
		secretMapName  string
		sourceFileName string
		gcpProjectID   string
		gkeClusterName string
		gceZoneName    string
	)

	app := kingpin.New("k8s", "kubernetes utilities").
		Version(version.Get().String()).
		Author(version.Get().ReleasedBy)

	gcpLoginCmd := app.Command("gcp-login", "Login to a GKE cluster")
	gcpLoginCmd.Flag("gcp-project", "GCP project id").Short('p').
		PlaceHolder("GCP_PROJECT_ID").Envar("GCP_PROJECT_ID").
		Required().StringVar(&gcpProjectID)
	gcpLoginCmd.Flag("gce-zone", "GCE zone").Short('z').
		PlaceHolder("GCE_ZONE").Envar("GCE_ZONE").
		Required().StringVar(&gceZoneName)
	gcpLoginCmd.Flag("gke-cluster", "GKE cluster name").Short('c').
		PlaceHolder("GKE_CLUSTER").Envar("GKE_CLUSTER").
		Required().StringVar(&gkeClusterName)

	genConfigMapCmd := app.Command("gen-configmap", "Generate a ConfigMap from a directory")
	genConfigMapCmd.Arg("config-directory", "Directory containing configuration").Default(".").ExistingDirVar(&configDir)
	genConfigMapCmd.Flag("name", "name of the ConfigMap").StringVar(&configMapName)

	genSecretMapCmd := app.Command("gen-secretmap", "Generate a SecretMap from a directory")
	genSecretMapCmd.Arg("secret-directory", "Directory containing secrets").Default(".").ExistingDirVar(&secretDir)
	genSecretMapCmd.Flag("name", "name of the Secret").StringVar(&secretMapName)

	genTemplateCmd := app.Command("gen-template", "Generate content from a template")
	genTemplateCmd.Arg("template", "Template file").Default(".").ExistingFileVar(&sourceFileName)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case gcpLoginCmd.FullCommand():
		gcpLogin(app, gcpProjectID, gceZoneName, gkeClusterName)
	case genConfigMapCmd.FullCommand():
		genConfigMap(app, configDir, configMapName)
	case genSecretMapCmd.FullCommand():
		genSecretMap(app, secretDir, secretMapName)
	case genTemplateCmd.FullCommand():
		genTemplate(app, sourceFileName)
	}
}
