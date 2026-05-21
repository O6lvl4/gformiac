package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version is set via ldflags at build time.
var Version = "dev"

var (
	specFile        string
	credentialsFile string
	tokenFile       string
	stateFile       string
)

var rootCmd = &cobra.Command{
	Use:     "gformiac",
	Short:   "Google Forms Infrastructure as Code",
	Long:    "YAML定義からGoogle Formsを宣言的に管理するIaCツール",
	Version: Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&specFile, "file", "f", "form.yaml", "フォーム定義ファイル")
	rootCmd.PersistentFlags().StringVar(&credentialsFile, "credentials", "credentials.json", "OAuth2認証情報ファイル")
	rootCmd.PersistentFlags().StringVar(&tokenFile, "token", "token.json", "OAuthトークンファイル")
	rootCmd.PersistentFlags().StringVar(&stateFile, "state", "gformiac.state.json", "状態ファイル")
}
