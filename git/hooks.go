// Package git provides all supported git hooks as constants in order to work on a reliable set.
package git

const (
	HookApplyPatchMsg = "applypatch-msg"
	HookCommitMsg = "commit-msg"
	HookFsMonitorWatchman = "fsmonitor-watchman"
	HookPostUpdate = "post-update"
	HookPreApplyPatch = "pre-applypatch"
	HookPreCommit = "pre-commit"
	HookPreMergeCommit = "pre-merge-commit"
	HookPrePush = "pre-push"
	HookPreRebase = "pre-rebase"
	HookPreReceive = "pre-receive"
	HookPrepareCommitMsg = "prepare-commit-msg"
	HookUpdate = "update"
)

// Hooks provides an array of all available Git hooks
var Hooks = []string{
	HookApplyPatchMsg,
	HookCommitMsg,
	HookFsMonitorWatchman,
	HookPostUpdate,
	HookPreApplyPatch,
	HookPreCommit,
	HookPreMergeCommit,
	HookPrePush,
	HookPreRebase,
	HookPreReceive,
	HookPrepareCommitMsg,
	HookUpdate,
}
