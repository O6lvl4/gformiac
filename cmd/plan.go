package cmd

import (
	"context"
	"fmt"

	"github.com/O6lvl4/gformiac/engine"
	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "変更のプレビュー（dry-run）",
	RunE:  runPlan,
}

func init() {
	rootCmd.AddCommand(planCmd)
}

func runPlan(cmd *cobra.Command, args []string) error {
	spec, err := engine.LoadSpec(specFile)
	if err != nil {
		return err
	}
	if err := engine.Validate(spec); err != nil {
		return err
	}

	state, err := engine.LoadState(stateFile)
	if err != nil {
		return fmt.Errorf("状態ファイル読み込み失敗: %w", err)
	}

	if state == nil {
		fmt.Println(engine.NewFormSummary(spec))
		return nil
	}

	return planExisting(spec, state)
}

func planExisting(spec *engine.FormSpec, state *engine.State) error {
	ctx := context.Background()
	client, err := engine.NewClient(ctx, credentialsFile, tokenFile)
	if err != nil {
		return err
	}

	diff, err := client.Plan(ctx, state.FormID, spec, state)
	if err != nil {
		return err
	}

	if !diff.HasChanges() {
		fmt.Println("変更なし — フォームは最新です")
		return nil
	}

	fmt.Println(diff.String())
	return nil
}
