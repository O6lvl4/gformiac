package cmd

import (
	"context"
	"fmt"

	"github.com/O6lvl4/gformiac/engine"
	"github.com/spf13/cobra"
)

var outputFile string

var importCmd = &cobra.Command{
	Use:   "import [formID]",
	Short: "既存のGoogle FormをYAML定義にインポート",
	Args:  cobra.ExactArgs(1),
	RunE:  runImport,
}

func init() {
	importCmd.Flags().StringVarP(&outputFile, "output", "o", "", "出力ファイル（未指定時は --file の値）")
	rootCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, args []string) error {
	formID := args[0]
	out := specFile
	if outputFile != "" {
		out = outputFile
	}

	ctx := context.Background()
	client, err := engine.NewClient(ctx, credentialsFile, tokenFile)
	if err != nil {
		return err
	}

	spec, state, err := client.ImportForm(ctx, formID)
	if err != nil {
		return err
	}

	if err := engine.SaveSpec(out, spec); err != nil {
		return fmt.Errorf("定義ファイル保存失敗: %w", err)
	}

	if err := engine.SaveState(stateFile, state); err != nil {
		return fmt.Errorf("状態ファイル保存失敗: %w", err)
	}

	fmt.Printf("インポート完了!\n")
	fmt.Printf("  定義ファイル: %s\n", out)
	fmt.Printf("  状態ファイル: %s\n", stateFile)
	fmt.Printf("  項目数: %d\n", len(spec.Items))
	return nil
}
