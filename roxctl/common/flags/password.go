package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/sync"
)

var (
	password        string
	passwordChanged *bool

	passwordOnce    sync.Once
	passwordFlagSet *pflag.FlagSet
)

// AddPassword adds the password flag to the base command.
func AddPassword(c *cobra.Command) {
	passwordOnce.Do(func() {
		passwordFlagSet = pflag.NewFlagSet("", pflag.ContinueOnError)
		passwordFlagSet.StringVarP(&password, "password", "p", "",
			"password for basic auth. Alternatively, set the password via the ROX_ADMIN_PASSWORD environment variable")
		passwordChanged = &passwordFlagSet.Lookup("password").Changed
	})
	c.PersistentFlags().AddFlagSet(passwordFlagSet)
}

// Password returns the set password.
func Password() string {
	return flagOrSettingValue(password, *passwordChanged, env.PasswordEnv)
}
