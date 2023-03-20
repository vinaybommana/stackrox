/*
Copyright 2021 Red Hat.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pkg/errors"
	platform "github.com/stackrox/rox/operator/apis/platform/v1alpha1"
	centralReconciler "github.com/stackrox/rox/operator/pkg/central/reconciler"
	securedClusterReconciler "github.com/stackrox/rox/operator/pkg/securedcluster/reconciler"
	"github.com/stackrox/rox/operator/pkg/utils"
	"github.com/stackrox/rox/pkg/buildinfo"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/fileutils"
	"github.com/stackrox/rox/pkg/profiling"
	"github.com/stackrox/rox/pkg/version"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	setupLog                   = ctrl.Log.WithName("setup")
	scheme                     = runtime.NewScheme()
	enableWebhooks             = env.RegisterBooleanSetting("ENABLE_WEBHOOKS", true)
	enableProfiling            = env.RegisterBooleanSetting("ENABLE_PROFILING", false)
	profilingThresholdFraction = env.RegisterSetting("PROFILING_THRESHOLD_FRACTION", env.WithDefault("0.8"))
	memLimit                   = env.RegisterIntegerSetting("MEMORY_LIMIT_BYTES", 0)
	// Default place to put the heap dump is the /tmp directory because the container process has rights
	// to write to this directory without creating and mounting a PVC
	heapDumpDir = env.RegisterSetting("HEAP_DUMP_PARENT_DIR", env.WithDefault("/tmp"))
	// Default place where controller-runtime looks for TLS artifacts.
	// see https://github.com/kubernetes-sigs/controller-runtime/blob/v0.8.3/pkg/webhook/server.go#L96-L104
	defaultCertDir  = filepath.Join(os.TempDir(), "k8s-webhook-server", "serving-certs")
	defaultTLSPaths = []string{filepath.Join(defaultCertDir, "tls.crt"), filepath.Join(defaultCertDir, "tls.key")}
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(platform.AddToScheme(scheme))
}

func main() {
	if err := run(); err != nil {
		setupLog.Error(err, "fatal error")
		os.Exit(1)
	}
}

func run() error {
	setupLog.Info("Starting RHACS Operator", "version", version.GetMainVersion())

	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: !buildinfo.ReleaseBuild,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(utils.GetRHACSConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "bf7ea6a2.stackrox.io",
		LeaderElectionConfig:   utils.GetRHACSConfigOrDie(),
	})
	if err != nil {
		return errors.Wrap(err, "unable to create manager")
	}

	if !enableWebhooks.BooleanSetting() {
		setupLog.Info("skipping webhook setup, ENABLE_WEBHOOKS==false")
	} else {
		maybeUseLegacyTLSFileLocation(mgr)
		if err = (&platform.Central{}).SetupWebhookWithManager(mgr); err != nil {
			return errors.Wrap(err, "unable to create Central webhook")
		}
	}

	if enableProfiling.BooleanSetting() {
		thresholdS := profilingThresholdFraction.Setting()
		thresholdF, err := strconv.ParseFloat(thresholdS, 32)
		if err != nil {
			return errors.Wrapf(err, "unable to parse PROFILING_THRESHOLD set to '%s' as a float", thresholdS)
		}
		heapProfiler := profiling.NewHeapProfiler(thresholdF, uint64(memLimit.IntegerSetting()), heapDumpDir.Setting(), profiling.DefaultHeapProfilerBackoff)
		ctx, cancelProfiler := context.WithCancel(context.Background())
		go heapProfiler.DumpHeapOnThreshhold(ctx, time.Second)
		defer cancelProfiler()
	}

	// The following comment marks the place where `operator-sdk` inserts new scaffolded code.
	//+kubebuilder:scaffold:builder

	if err = centralReconciler.RegisterNewReconciler(mgr); err != nil {
		return errors.Wrap(err, "unable to set up Central reconciler")
	}

	if err = securedClusterReconciler.RegisterNewReconciler(mgr); err != nil {
		return errors.Wrap(err, "unable to set up SecuredCluster reconciler")
	}

	if err = mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		return errors.Wrap(err, "unable to set up health check")
	}
	if err = mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		return errors.Wrap(err, "unable to set up readiness check")
	}

	setupLog.Info("starting manager")
	if err = mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		return errors.Wrap(err, "problem running manager")
	}
	return nil
}

func maybeUseLegacyTLSFileLocation(mgr manager.Manager) {
	// OLM before version 0.17.0 (such as the one shipped with OpenShift 4.6) does not
	// provide the TLS certificate/key in the location referenced by default by the controller runtime
	// (i.e. /tmp/k8s-webhook-server/serving-certs/...).
	// See ROX-8304 for details.
	// If the files are missing at the default location, then we explicitly set the settings as follows
	// to force usage of the legacy location, which is provided both by old and new OLM, but not
	// by the "make deploy" scaffolding.
	if ok, _ := fileutils.AllExist(defaultTLSPaths...); ok {
		return
	}
	setupLog.Info("Webhook key and/or certificate missing at default paths, attempting use of legacy path.", "defaultTLSPaths", defaultTLSPaths)
	server := mgr.GetWebhookServer()
	server.CertDir = "/apiserver.local.config/certificates"
	server.CertName = "apiserver.crt"
	server.KeyName = "apiserver.key"
}
