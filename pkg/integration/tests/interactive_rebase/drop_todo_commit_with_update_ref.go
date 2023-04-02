package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DropTodoCommitWithUpdateRef = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Drops a commit during interactive rebase when there is an update-ref in the git-rebase-todo file",
	ExtraCmdArgs: "",
	Skip:         false,
	GitVersion:   From("2.38.0"),
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(3).
			NewBranch("mybranch").
			CreateNCommitsStartingAt(3, 4)

		shell.SetConfig("rebase.updateRefs", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("(*) commit 06").IsSelected(),
				Contains("commit 05"),
				Contains("commit 04"),
				Contains("(*) commit 03"),
				Contains("commit 02"),
				Contains("commit 01"),
			).
			// Once "e" is fixed we can just hit "e", but for now we need to
			// manually do a command-line rebase
			// NavigateToLine(Contains("commit 01")).
			// Press(keys.Universal.Edit).
			Tap(func() {
				t.GlobalPress(keys.Universal.ExecuteCustomCommand)
				t.ExpectPopup().Prompt().
					Title(Equals("Custom Command:")).
					Type(`git -c core.editor="perl -i -lpe 'print \"break\" if $.==1'" rebase -i HEAD~5`).
					Confirm()
			}).
			Focus().
			Lines(
				Contains("pick").Contains("(*) commit 06"),
				Contains("pick").Contains("commit 05"),
				Contains("pick").Contains("commit 04"),
				Contains("update-ref").Contains("master"),
				Contains("pick").Contains("(*) commit 03"),
				Contains("pick").Contains("commit 02"),
				Contains("<-- YOU ARE HERE --- commit 01"),
			).
			NavigateToLine(Contains("commit 05")).
			Press(keys.Universal.Remove)

		t.Common().ContinueRebase()

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("(*) commit 06"),
				Contains("commit 04"),
				Contains("(*) commit 03"),
				Contains("commit 02"),
				Contains("commit 01"),
			)
	},
})
