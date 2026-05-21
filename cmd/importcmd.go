package cmd

import (
	"context"
	"fmt"

	"github.com/O6lvl4/gformiac/engine"
	"github.com/O6lvl4/gformiac/locale"
	"github.com/spf13/cobra"
)

// outputFile is the path for the imported YAML spec; defaults to --file when empty.
var outputFile string

var importCmd = &cobra.Command{
	Use:   "import [formID]",
	Short: locale.M.ImportShort,
	Args:  cobra.ExactArgs(1),
	RunE:  runImport,
}

func init() {
	importCmd.Flags().StringVarP(&outputFile, "output", "o", "", locale.M.FlagOutput)
	rootCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, args []string) error {
	formID := args[0]
	out := importOutputPath()

	ctx := context.Background()
	client, err := engine.NewClient(ctx, credentialsFile, tokenFile)
	if err != nil {
		return err
	}

	spec, state, err := client.ImportForm(ctx, formID)
	if err != nil {
		return err
	}

	return saveImportResult(out, spec, state)
}

func importOutputPath() string {
	if outputFile != "" {
		return outputFile
	}
	return specFile
}

func saveImportResult(out string, spec *engine.FormSpec, state *engine.State) error {
	if err := engine.SaveSpec(out, spec); err != nil {
		return fmt.Errorf("%s: %w", locale.M.ErrSpecSave, err)
	}
	if err := engine.SaveState(stateFile, state); err != nil {
		return fmt.Errorf("%s: %w", locale.M.ErrStateSave, err)
	}

	fmt.Println(locale.M.Imported)
	fmt.Printf(locale.M.SpecFileLabel+"\n", out)
	fmt.Printf(locale.M.StateLabel+"\n", stateFile)
	fmt.Printf(locale.M.ItemCount+"\n", len(spec.Items))
	return nil
}
