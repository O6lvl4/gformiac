package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/O6lvl4/gformiac/engine"
	"github.com/spf13/cobra"
)

var autoApprove bool

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "フォーム定義を適用",
	RunE:  runApply,
}

func init() {
	applyCmd.Flags().BoolVar(&autoApprove, "auto-approve", false, "確認をスキップ")
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
		return fmt.Errorf("状態ファイル読み込み失敗: %w", err)
	}

	ctx := context.Background()
	if state == nil {
		return applyNew(ctx, spec)
	}
	return applyUpdate(ctx, spec, state)
}

func applyNew(ctx context.Context, spec *engine.FormSpec) error {
	fmt.Println(engine.NewFormSummary(spec))
	if !autoApprove && !confirm("適用しますか？") {
		fmt.Println("キャンセルしました")
		return nil
	}

	client, err := engine.NewClient(ctx, credentialsFile, tokenFile)
	if err != nil {
		return err
	}

	fmt.Println("フォーム作成中...")
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
		fmt.Println("変更なし — フォームは最新です")
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
	if !confirm("\n適用しますか？") {
		fmt.Println("キャンセルしました")
		return false
	}
	return true
}

func saveAndReport(state *engine.State) error {
	if err := engine.SaveState(stateFile, state); err != nil {
		return fmt.Errorf("状態保存失敗: %w", err)
	}
	fmt.Println()
	fmt.Println("適用完了!")
	fmt.Printf("  フォームID:  %s\n", state.FormID)
	fmt.Printf("  回答URL:     %s\n", state.ResponderURL)
	fmt.Printf("  状態ファイル: %s\n", stateFile)
	return nil
}

func confirm(prompt string) bool {
	fmt.Printf("%s [y/N] ", prompt)
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}
