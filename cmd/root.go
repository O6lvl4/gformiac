// Package cmd implements the gformiac CLI commands using cobra.
// It wires together the engine and locale packages to provide plan, apply,
// and import subcommands with persistent flags for credentials and state.
package cmd

import (
	"fmt"
	"os"

	"github.com/O6lvl4/gformiac/locale"
	"github.com/spf13/cobra"
)

// Version is the CLI version string, set via ldflags at build time.
var Version = "dev"

// Persistent flag values shared across all subcommands.
var (
	specFile        string
	credentialsFile string
	tokenFile       string
	stateFile       string
	langFlag        string
)

var rootCmd = &cobra.Command{
	Use:     "gformiac",
	Short:   "Google Forms Infrastructure as Code",
	Version: Version,
}

// Execute runs the root cobra command and exits the process on error.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Long = locale.M.CmdLong
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if langFlag != "" {
			locale.Set(locale.Lang(langFlag))
		}
	}

	rootCmd.PersistentFlags().StringVarP(&specFile, "file", "f", "form.yaml", locale.M.FlagFile)
	rootCmd.PersistentFlags().StringVar(&credentialsFile, "credentials", ".gformiac/credentials.json", locale.M.FlagCredentials)
	rootCmd.PersistentFlags().StringVar(&tokenFile, "token", ".gformiac/token.json", locale.M.FlagToken)
	rootCmd.PersistentFlags().StringVar(&stateFile, "state", ".gformiac/state.json", locale.M.FlagState)
	rootCmd.PersistentFlags().StringVar(&langFlag, "lang", "", locale.M.FlagLang)
}
