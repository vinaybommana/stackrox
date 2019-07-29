package deploy

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/renderer"
	"github.com/stackrox/rox/pkg/roxctl"
	"github.com/stackrox/rox/pkg/roximages/defaults"
	"github.com/stackrox/rox/pkg/utils"
	"github.com/stackrox/rox/roxctl/common/flags"
)

type persistentFlagsWrapper struct {
	*pflag.FlagSet
}

func (w *persistentFlagsWrapper) StringVar(p *string, name, value, usage string, groups ...string) {
	w.StringVarP(p, name, "", value, usage, groups...)
}

func (w *persistentFlagsWrapper) StringVarP(p *string, name, shorthand, value, usage string, groups ...string) {
	w.FlagSet.StringVarP(p, name, shorthand, value, usage)
	utils.Must(w.SetAnnotation(name, groupAnnotationKey, groups))
}

func (w *persistentFlagsWrapper) BoolVar(p *bool, name string, value bool, usage string, groups ...string) {
	w.FlagSet.BoolVar(p, name, value, usage)
	utils.Must(w.SetAnnotation(name, groupAnnotationKey, groups))
}

func (w *persistentFlagsWrapper) Var(value pflag.Value, name, usage string, groups ...string) {
	w.FlagSet.Var(value, name, usage)
	utils.Must(w.SetAnnotation(name, groupAnnotationKey, groups))
}

func orchestratorCommand(shortName, longName string) *cobra.Command {
	c := &cobra.Command{
		Use:   shortName,
		Short: fmt.Sprintf("%s specifies that you are going to launch StackRox Central in %s.", shortName, longName),
		Long: fmt.Sprintf(`%s specifies that you are going to launch StackRox Central in %s.
Output is a zip file printed to stdout.`, shortName, longName),
		Annotations: map[string]string{
			categoryAnnotation: "Enter orchestrator",
		},
		RunE: func(*cobra.Command, []string) error {
			return fmt.Errorf("storage type must be specified")
		},
	}
	if !roxctl.InMainImage() {
		c.PersistentFlags().Var(newOutputDir(&cfg.OutputDir), "output-dir", "the directory to output the deployment bundle to")
	}
	return c
}

func k8sBasedOrchestrator(k8sConfig *renderer.K8sConfig, shortName, longName string, cluster storage.ClusterType) *cobra.Command {
	c := orchestratorCommand(shortName, longName)
	c.PersistentPreRun = func(*cobra.Command, []string) {
		cfg.K8sConfig = k8sConfig
		cfg.ClusterType = cluster
	}

	c.AddCommand(externalVolume())
	c.AddCommand(hostPathVolume())
	c.AddCommand(noVolume())

	flagWrap := &persistentFlagsWrapper{FlagSet: c.PersistentFlags()}

	// Adds k8s specific flags
	flagWrap.StringVarP(&k8sConfig.MainImage, "main-image", "i", defaults.MainImage(), "main image to use", "central")
	flagWrap.BoolVar(&k8sConfig.OfflineMode, "offline", false, "whether to run StackRox in offline mode, which avoids reaching out to the Internet", "central")

	// Monitoring Flags
	flagWrap.StringVar(&k8sConfig.Monitoring.Endpoint, "monitoring-endpoint", "monitoring.stackrox:443", "monitoring endpoint", "monitoring", "monitoring-type=on-prem")
	flagWrap.Var(&monitoringWrapper{Monitoring: &k8sConfig.Monitoring.Type}, "monitoring-type", "where to host the monitoring (on-prem, none)", "monitoring")

	flagWrap.StringVarP(&k8sConfig.Monitoring.Password, "monitoring-password", "p", "", "a monitoring password (default: autogenerated)", "monitoring", "monitoring-type=on-prem")
	utils.Must(
		flagWrap.SetAnnotation("monitoring-password", flags.PasswordKey, []string{"true"}))
	flagWrap.StringVar(&k8sConfig.MonitoringImage, "monitoring-image", "", "monitoring image to use (default: same repository as main)", "monitoring", "monitoring-type=on-prem")

	// Monitoring Persistence flags
	flagWrap.Var(&persistenceTypeWrapper{PersistenceType: &k8sConfig.Monitoring.PersistenceType}, "monitoring-persistence-type", "monitoring persistence type (none, hostpath, pvc)", "monitoring", "monitoring-type=on-prem")

	flagWrap.StringVar(&k8sConfig.Monitoring.External.Name, "monitoring-persistence-name", "monitoring-db", "external volume name", "monitoring", "monitoring-type=on-prem", "monitoring-persistence-type=pvc")
	flagWrap.StringVar(&k8sConfig.Monitoring.External.StorageClass, "monitoring-persistence-storage-class", "", "monitoring storage class name (optional if you have a default StorageClass configured)", "monitoring", "monitoring-type=on-prem", "monitoring-persistence-type=pvc")

	flagWrap.StringVar(&k8sConfig.Monitoring.HostPath.HostPath, "monitoring-persistence-hostpath", "/var/lib/stackrox/monitoring", "monitoring path on the host", "monitoring", "monitoring-type=on-prem", "monitoring-persistence-type=hostpath")
	flagWrap.StringVar(&k8sConfig.Monitoring.HostPath.NodeSelectorKey, "monitoring-node-selector-key", "", "monitoring node selector key (e.g. kubernetes.io/hostname)", "monitoring", "monitoring-type=on-prem", "monitoring-persistence-type=hostpath")
	flagWrap.StringVar(&k8sConfig.Monitoring.HostPath.NodeSelectorValue, "monitoring-node-selector-value", "", "monitoring node selector value", "monitoring", "monitoring-type=on-prem", "monitoring-persistence-type=hostpath")

	// Scanner
	flagWrap.StringVar(&k8sConfig.ScannerImage, "scanner-image", defaults.ScannerImage(), "Scanner image to use", "scanner")
	flagWrap.BoolVar(&k8sConfig.EnableScannerV2, "enable-scanner-v2", false, "Whether to enable scanner v2", "scanner")
	utils.Must(flagWrap.MarkHidden("enable-scanner-v2"))
	flagWrap.StringVar(&k8sConfig.ScannerV2DBImage, "scanner-db-image", defaults.ScannerV2DBImage(), "Scanner V2 DB image to use", "scanner")
	utils.Must(flagWrap.MarkHidden("scanner-db-image"))

	return c
}

func newK8sConfig() *renderer.K8sConfig {
	return &renderer.K8sConfig{
		Monitoring: renderer.MonitoringConfig{
			HostPath: &renderer.HostPathPersistence{},
			External: &renderer.ExternalPersistence{},
		},
	}
}

func k8s() *cobra.Command {
	k8sConfig := newK8sConfig()
	c := k8sBasedOrchestrator(k8sConfig, "k8s", "Kubernetes", storage.ClusterType_KUBERNETES_CLUSTER)
	flagWrap := &persistentFlagsWrapper{FlagSet: c.PersistentFlags()}

	flagWrap.Var(&loadBalancerWrapper{LoadBalancerType: &k8sConfig.LoadBalancerType}, "lb-type", "the method of exposing Central (lb, np, none)", "central")

	flagWrap.Var(&fileFormatWrapper{DeploymentFormat: &k8sConfig.DeploymentFormat}, "output-format", "the deployment tool to use (kubectl, helm)", "central")

	flagWrap.Var(&loadBalancerWrapper{LoadBalancerType: &k8sConfig.Monitoring.LoadBalancerType}, "monitoring-lb-type", "the method of exposing Monitoring (lb, np, none)", "monitoring", "monitoring-type=on-prem")

	return c
}

func openshift() *cobra.Command {
	k8sConfig := newK8sConfig()
	c := k8sBasedOrchestrator(k8sConfig, "openshift", "Openshift", storage.ClusterType_OPENSHIFT_CLUSTER)

	flagWrap := &persistentFlagsWrapper{FlagSet: c.PersistentFlags()}

	flagWrap.Var(&loadBalancerWrapper{LoadBalancerType: &k8sConfig.LoadBalancerType}, "lb-type", "the method of exposing Central (route, lb, np, none)", "central")

	flagWrap.Var(&loadBalancerWrapper{LoadBalancerType: &k8sConfig.Monitoring.LoadBalancerType}, "monitoring-lb-type", "the method of exposing Monitoring (route, lb, np, none)", "monitoring", "monitoring-type=on-prem")

	return c
}
