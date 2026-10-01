# Consensus-spec vector checks

The supported release is declared per fork in
[`.github/consensus-spec-versions.json`](../.github/consensus-spec-versions.json).
Only **Gloas** declares support for **v1.7.0-beta.2**. This declaration covers the
mainnet static-vector types exercised by `spec/gloas/TestConsensusSpec`, not
consensus state transitions, gossip validation, or live HTTP compatibility.
It does not declare beta.2 support for older forks or the minimal preset.

## CI

- `consensus-spec-vectors` reads the Gloas declaration and verifies that exact
  release on pull requests, manual dispatch, and pushes to `master`.
- `consensus-spec-release` checks upstream every Friday at 10:00 UTC, including
  alpha, beta and rc releases. It compares version precedence, not publication
  dates, and tests only the fork selected by the repository Actions variable
  `CONSENSUS_SPEC_FORK`, initially `gloas`. The selected fork must have an entry
  in the support declaration.
- A newer release opens a GitHub issue with the tested commit, result, run/log
  link, release-review checklist and target-specific reproduction instructions.
  Passing vectors does not automatically update the support declaration.
- Open and closed bot issues with the same fork/release marker prevent duplicate
  issues and repeat scheduled tests. To retest a reported release, use the local
  runner below. Download, build, missing-vector and skipped-test problems without
  evidenced vector failures are reported as `not-run`. Evidenced vector failures
  are reported as `fail`, even if other cases were skipped. Cancelled runs do not
  open issues; an unreported release is checked again on the next scheduled run.

Both scheduled runs and manual dispatch require the workflow file on GitHub's
[default branch](https://docs.github.com/en/actions/how-tos/manage-workflow-runs/manually-run-a-workflow).
This repository's default is `master`, so merging into `gloas` alone does not
activate either trigger. Once the workflows reach `master`, manual dispatch can
select another branch with `--ref`. Pull-request vector checks do not require
that default-branch promotion.

Set or inspect the watcher variable with:

```sh
gh variable set CONSENSUS_SPEC_FORK --body gloas
gh variable get CONSENSUS_SPEC_FORK
```

## Local verification

Requires Go, Python 3.12+, GitHub CLI authenticated for public API access, and
curl. The runner uses normal `go test`; CI selects Go from the checked-out
`go.mod`, so the existing Go updater also controls vector-test toolchains.
Run from the target checkout:

```sh
(
set -eu
work="$(mktemp -d)"
trap 'rm -rf "$work"' EXIT
GOTOOLCHAIN=auto python3 .github/scripts/consensus_vectors.py \
  v1.7.0-beta.2 gloas "$work/vectors"
)
```

The runner downloads the exact release's `mainnet.tar.gz`, checks GitHub's
published SHA-256 digest, and extracts only the selected fork's `ssz_static`
files. Allow approximately 1 GB for the compressed download, plus extracted
files. A passing result requires executed vector subtests with no skips.
CI retains the log as an artifact. Locally, the snippet prints the log and removes
temporary files on subshell exit. Redirect its output to a file to keep a local log.

For a target without the runner, copy the script to a temporary location before
checking out the target. The release-review issue includes this procedure. It
fetches the recorded automation commit explicitly, then fetches the target from
origin or a fork URL and checks out the fetched commit detached, avoiding stale
local branches. Use a branch, tag or full commit SHA, not an abbreviated SHA.
For local-only targets, use `TARGET_REMOTE=.` and resolve abbreviated SHAs with
`git rev-parse` first. GitHub PR targets can use `refs/pull/<number>/head` from origin.
The snippet prints the tested commit. The automation commit must remain available
from origin; if GitHub no longer serves it, restore it before reproducing.

After reviewing a new release and resolving any incompatibilities, change only
the verified fork's version in the declaration through a PR. The pinned workflow
then verifies the new baseline; the watcher never promotes it automatically.
