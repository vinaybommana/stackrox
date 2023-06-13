package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/stackrox/rox/central/clusters"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/image"
	"github.com/stackrox/rox/image/sensor"
	"github.com/stackrox/rox/pkg/helm/util"
	"github.com/stackrox/rox/pkg/images/defaults"
	"github.com/stackrox/rox/pkg/version"
	"github.com/stackrox/rox/pkg/zip"
	"helm.sh/helm/v3/pkg/chartutil"
)

var installCmd = &cobra.Command{
	Use:  "install",
	Long: "Install the StackRox Secured Cluster",
	RunE: func(cmd *cobra.Command, args []string) error {

		//ctx := cmd.Context()
		//
		//cliEnvironment := environment.CLIEnvironment()
		//grpcConn, err := cliEnvironment.GRPCConnection()
		//if err != nil {
		//	return err
		//}
		//
		//service := v1.NewClusterInitServiceClient(grpcConn)
		//bundle, err := service.GenerateInitBundle(ctx, &v1.InitBundleGenRequest{
		//	Name: "test",
		//})
		//if err != nil {
		//	return err
		//}

		v := version.Versions{
			CollectorVersion:      "4.0.1",
			GitCommit:             "abc123",
			GoVersion:             runtime.Version(),
			MainVersion:           "4.0.1",
			Platform:              runtime.GOOS + "/" + runtime.GOARCH,
			ScannerVersion:        "4.0.1",
			ChartVersion:          "4.0.1",
			Database:              "",
			DatabaseServerVersion: "",
		}
		flavor := defaults.ImageFlavor{
			MainRegistry:       "registry.redhat.io/advanced-cluster-security",
			MainImageName:      "rhacs-main-rhel8",
			MainImageTag:       v.MainVersion,
			CentralDBImageTag:  v.MainVersion,
			CentralDBImageName: "rhacs-central-db-rhel8",

			CollectorRegistry:      "registry.redhat.io/advanced-cluster-security",
			CollectorImageName:     "rhacs-collector-rhel8",
			CollectorImageTag:      v.CollectorVersion,
			CollectorSlimImageName: "rhacs-collector-slim-rhel8",
			CollectorSlimImageTag:  v.CollectorVersion,

			ScannerImageName:       "rhacs-scanner-rhel8",
			ScannerSlimImageName:   "rhacs-scanner-slim-rhel8",
			ScannerImageTag:        v.ScannerVersion,
			ScannerDBImageName:     "rhacs-scanner-db-rhel8",
			ScannerDBSlimImageName: "rhacs-scanner-db-slim-rhel8",

			ChartRepo: defaults.ChartRepo{
				URL:     "https://mirror.openshift.com/pub/rhacs/charts",
				IconURL: "https://raw.githubusercontent.com/stackrox/stackrox/master/image/templates/helm/shared/assets/Red_Hat-Hat_icon.png",
			},
			ImagePullSecrets: defaults.ImagePullSecrets{
				AllowNone: true,
			},
			Versions: v,
		}

		c := &storage.Cluster{
			Type:      storage.ClusterType_OPENSHIFT4_CLUSTER,
			MainImage: "registry.redhat.io/advanced-cluster-security/rhacs-main-rhel8:4.0.1",
		}

		metaValues, err := clusters.FieldsFromClusterAndRenderOpts(c, &flavor, clusters.RenderOptions{})
		if err != nil {
			return err
		}
		metaValues.ClusterName = "test"

		helmImage := image.GetDefaultImage()
		ch, err := helmImage.GetSensorChart(metaValues, &sensor.Certs{Files: map[string][]byte{}})
		if err != nil {
			return errors.Wrap(err, "pre-rendering sensor chart")
		}
		ch.Values["config"] = map[string]interface{}{
			"createSecrets": false,
		}

		m, err := util.Render(ch, nil, util.Options{
			ReleaseOptions: chartutil.ReleaseOptions{
				Name:      "test",
				Namespace: "stackrox",
				IsInstall: true,
			},
		})
		if err != nil {
			return err
		}

		var renderedFiles []*zip.File
		// For kubectl files, we don't want to have the templates path so we trim it out
		for k, v := range m {
			if strings.TrimSpace(v) == "" {
				continue
			}
			var flags zip.FileFlags
			renderedFiles = append(renderedFiles, zip.NewFile(filepath.Base(k), []byte(v), flags))
		}

		if err := os.MkdirAll("manifests/out", 0755); err != nil {
			return err
		}

		for _, f := range renderedFiles {
			// check if extension is yaml
			if filepath.Ext(f.Name) != ".yaml" {
				continue
			}
			file, err := os.Create(filepath.Join("manifests", "out", f.Name))
			if err != nil {
				return err
			}
			if _, err := file.Write(f.Content); err != nil {
				file.Close()
				return err
			}
			file.Close()
		}

		return nil
	},
}

func main() {
	installCmd.Execute()
}
