package flags

import "github.com/spf13/cobra"

// AddCentralConnectivityFlags adds flags when central connectivity should be configured and used, in particular:
// - admin password information
// - endpoint information and protocol selection
// - API token file
// The flags will be added as persistent flags, and shared with all sub commands of the given command.
func AddCentralConnectivityFlags(cmd *cobra.Command) {
	AddPassword(cmd)
	AddConnectionFlags(cmd)
	AddAPITokenFile(cmd)
}
