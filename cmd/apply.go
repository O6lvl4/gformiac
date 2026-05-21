package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/O6lvl4/gformiac/engine"
	"github.com/O6lvl4/gformiac/locale"
	"github.com/spf13/cobra"
)

var autoApprove bool

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: locale.M.ApplyShort,
	RunE:  runApply,
}

func init() {
	applyCmd.Flags().BoolVar(&autoApprove, "auto-approve", false, locale.M.FlagAutoApprove)
	rootCmd.AddCommand(applyCmd)
}

func runApply(cmd *cobra.Command, args []string) error {
	spec, err := engine.LoadSpec(specFile)
	if err != nil {
		return err
	}
	if err := engine.Validate(spec); err != nil {
		return err
	}

	state, err := engine.LoadState(stateFile)
	if err != nil {
		return fmt.Errorf("%s: %w", locale.M.ErrStateRead, err)
	}

	ctx := context.Background()
	if state == nil {
		return applyNew(ctx, spec)
	}
	return applyUpdate(ctx, spec, state)
}

func applyNew(ctx context.Context, spec *engine.FormSpec) error {
	fmt.Println(engine.NewFormSummary(spec))
	if !autoApprove && !confirm(locale.M.ConfirmApply) {
		fmt.Println(locale.M.Cancelled)
		return nil
	}

	client, err := engine.NewClient(ctx, credentialsFile, tokenFile)
	if err != nil {
		return err
	}

	fmt.Println(locale.M.Creating)
	state, err := client.CreateForm(ctx, spec)
	if err != nil {
		return err
	}

	return saveAndReport(state)
}

func applyUpdate(ctx context.Context, spec *engine.FormSpec, state *engine.State) error {
	client, err := engine.NewClient(ctx, credentialsFile, tokenFile)
	if err != nil {
		return err
	}

	diff, err := client.Plan(ctx, state.FormID, spec, state)
	if err != nil {
		return err
	}
	if !diff.HasChanges() {
		fmt.Println(locale.M.NoChangesLong)
		return nil
	}
	if !confirmDiff(diff) {
		return nil
	}

	return executeUpdate(ctx, client, state.FormID, spec)
}

func executeUpdate(ctx context.Context, client *engine.Client, formID string, spec *engine.FormSpec) error {
	newState, err := client.UpdateForm(ctx, formID, spec)
	if err != nil {
		return err
	}
	return saveAndReport(newState)
}

func confirmDiff(diff *engine.DiffResult) bool {
	fmt.Println(diff.String())
	if autoApprove {
		return true
	}
	if !confirm("\n" + locale.M.ConfirmApply) {
		fmt.Println(locale.M.Cancelled)
		return false
	}
	return true
}

func saveAndReport(state *engine.State) error {
	if err := engine.SaveState(stateFile, state); err != nil {
		return fmt.Errorf("%s: %w", locale.M.ErrStateSave, err)
	}
	fmt.Println()
	fmt.Println(locale.M.Applied)
	fmt.Printf(locale.M.FormIDLabel+"\n", state.FormID)
	fmt.Printf(locale.M.URLLabel+"\n", state.ResponderURL)
	fmt.Printf(locale.M.StateLabel+"\n", stateFile)
	return nil
}

func confirm(prompt string) bool {
	fmt.Printf("%s [y/N] ", prompt)
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}
