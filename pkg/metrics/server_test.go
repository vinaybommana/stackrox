package metrics

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stackrox/rox/pkg/buildinfo"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricsServerAddressEnvs(t *testing.T) {
	cases := map[string]struct {
		metricsPort       string
		secureMetricsPort string
	}{
		"default": {
			metricsPort:       "",
			secureMetricsPort: "",
		},
		"only metricsPort set": {
			metricsPort:       ":8008",
			secureMetricsPort: "",
		},
		"only secureMetricsPort set": {
			metricsPort:       "",
			secureMetricsPort: ":8009",
		},
		"metrisPort and secureMetricsPort set": {
			metricsPort:       "8008",
			secureMetricsPort: ":8009",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Setenv(env.MetricsPort.EnvVar(), c.metricsPort)
			t.Setenv(env.SecureMetricsPort.EnvVar(), c.secureMetricsPort)

			server := NewServer(CentralSubsystem)

			require.NotNil(t, server)
			assert.Equal(t, env.MetricsPort.Setting(), server.metricsServer.Addr)
			assert.Equal(t, env.SecureMetricsPort.Setting(), server.secureMetricsServer.Addr)
		})
	}
}

func TestMetricsServerPanic(t *testing.T) {
	cases := map[string]struct {
		metricsPort       string
		secureMetricsPort string
		releaseBuild      bool
	}{
		"metrics error - debug build panics": {
			metricsPort:       "error",
			secureMetricsPort: "",
			releaseBuild:      false,
		},
		"metrics error - release build does not panic": {
			metricsPort:       "error",
			secureMetricsPort: "",
			releaseBuild:      true,
		},
		"secureMetrics error - debug build panics": {
			metricsPort:       "",
			secureMetricsPort: "error",
			releaseBuild:      false,
		},
		"secureMetrics error - release build does not panic": {
			metricsPort:       "",
			secureMetricsPort: "error",
			releaseBuild:      true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if buildinfo.ReleaseBuild != c.releaseBuild {
				t.SkipNow()
			}
			t.Setenv(env.MetricsPort.EnvVar(), c.metricsPort)
			t.Setenv(env.SecureMetricsPort.EnvVar(), c.secureMetricsPort)
			server := NewServer(CentralSubsystem)

			if c.releaseBuild {
				assert.NotPanics(t, func() { server.RunForever() })
			} else {
				assert.Panics(t, func() { server.RunForever() })
			}
			server.Stop(context.TODO())
		})
	}
}

func TestMetricsServerHTTPRequest(t *testing.T) {
	t.Setenv(env.SecureMetricsPort.EnvVar(), "disabled")
	server := NewServer(CentralSubsystem)
	server.RunForever()

	request := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	recorder := httptest.NewRecorder()
	server.metricsServer.Handler.ServeHTTP(recorder, request)
	resp := recorder.Result()
	_, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	resp.Body.Close()
	server.Stop(context.TODO())
}


// TOOD: test https server
