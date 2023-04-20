package flags

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/errox"
	"github.com/stackrox/rox/pkg/sync"
)

var (
	endpoint        string
	endpointChanged *bool

	serverName string
	directGRPC bool
	forceHTTP1 bool

	plaintext    bool
	plaintextSet *bool
	insecure     bool

	insecureSkipTLSVerify    bool
	insecureSkipTLSVerifySet *bool

	caCertFile string

	connectionFlagsOnce sync.Once
	connectionFlagSet   *flag.FlagSet
)

// AddConnectionFlags adds connection-related flags to roxctl.
func AddConnectionFlags(c *cobra.Command) {
	connectionFlagsOnce.Do(func() {
		connectionFlagSet = flag.NewFlagSet("", flag.ContinueOnError)
		connectionFlagSet.StringVarP(&endpoint, "endpoint", "e", "localhost:8443",
			"endpoint for service to contact. Alternatively, set the endpoint via the ROX_ENDPOINT environment variable")
		endpointChanged = &connectionFlagSet.Lookup("endpoint").Changed
		connectionFlagSet.StringVarP(&serverName, "server-name", "s", "", "TLS ServerName to use for SNI (if empty, derived from endpoint)")
		connectionFlagSet.BoolVar(&directGRPC, "direct-grpc", false, "Use direct gRPC (advanced; only use if you encounter connection issues)")
		connectionFlagSet.BoolVar(&forceHTTP1, "force-http1", false, "Always use HTTP/1 for all connections (advanced; only use if you encounter connection issues)")

		connectionFlagSet.BoolVar(&plaintext, "plaintext", false, "Use a plaintext (unencrypted) connection; only works in conjunction with --insecure")
		plaintextSet = &connectionFlagSet.Lookup("plaintext").Changed
		connectionFlagSet.BoolVar(&insecure, "insecure", false, "Enable insecure connection options (DANGEROUS; USE WITH CAUTION)")
		connectionFlagSet.BoolVar(&insecureSkipTLSVerify, "insecure-skip-tls-verify", false, "Skip TLS certificate validation")
		insecureSkipTLSVerifySet = &connectionFlagSet.Lookup("insecure-skip-tls-verify").Changed
		connectionFlagSet.StringVar(&caCertFile, "ca", "", "Custom CA certificate to use (PEM format)")
	})
	c.PersistentFlags().AddFlagSet(connectionFlagSet)
}

// EndpointAndPlaintextSetting returns the Central endpoint to connect to, as well as a bool indicating whether to
// connect in plaintext mode.
func EndpointAndPlaintextSetting() (string, bool, error) {
	endpoint = flagOrSettingValue(endpoint, *endpointChanged, env.EndpointEnv)
	if !strings.Contains(endpoint, "://") {
		return endpoint, plaintext, nil
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return "", false, errors.Wrap(err, "malformed endpoint URL")
	}

	if u.Path != "" && u.Path != "/" {
		return "", false, errox.InvalidArgs.New("endpoint URL must not include a path component")
	}

	var usePlaintext bool
	switch u.Scheme {
	case "http":
		usePlaintext = true
	case "https":
		usePlaintext = false
	default:
		return "", false, errox.InvalidArgs.Newf("invalid scheme %q in endpoint URL, the scheme should be: http(s)://<endpoint>:<port>", u.Scheme)
	}

	if *plaintextSet {
		if plaintext != usePlaintext {
			return "", false, errox.InvalidArgs.Newf("endpoint URL scheme %q is incompatible with --plaintext=%v setting", u.Scheme, plaintext)
		}
	}

	return u.Host, usePlaintext, nil
}

// ServerName returns the specified ServerName.
func ServerName() string {
	return serverName
}

// UseDirectGRPC returns whether to use gRPC directly, i.e., without a proxy.
func UseDirectGRPC() bool {
	return directGRPC
}

// ForceHTTP1 indicates that the HTTP/1 should be used for all outgoing connections.
func ForceHTTP1() bool {
	return forceHTTP1
}

// UseInsecure returns whether to use insecure connection behavior.
func UseInsecure() bool {
	return insecure
}

// SkipTLSValidation returns a bool that indicates the value of the `--insecure-skip-tls-verify` flag, with `nil`
// indicating that it was left at its default value.
func SkipTLSValidation() *bool {
	if !*insecureSkipTLSVerifySet {
		return nil
	}
	return &insecureSkipTLSVerify
}

// CAFile returns the file for custom CA certificates.
func CAFile() string {
	return caCertFile
}
