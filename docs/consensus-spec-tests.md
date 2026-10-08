# Consensus-spec vector checks

The supported release is declared per fork in
[`.github/consensus-spec-versions.json`](../.github/consensus-spec-versions.json).
Only Gloas declares support for v1.7.0-beta.2. This declaration covers mainnet
static-vector YAML/SSZ encodings and hash-tree roots, not consensus state
transitions, gossip validation, live HTTP compatibility, older forks or the minimal preset.

## Coverage policy

A pass requires every nonexcluded archive case to execute exactly once and pass,
with successful root-test and package results and no skips. The runner reconciles
native `go test -json` terminal events with the extracted case inventory. Go command
diagnostics remain on stderr and do not enter the JSON classifier.
A genuine failed covered case reports `fail`, even if other cases did not execute.
Setup failures, malformed evidence and otherwise incomplete execution report `not-run`.
Missing or empty required type directories also fail the Go test itself.

The approved Gloas exclusions are explicit in
[the runner](../.github/scripts/consensus_vectors.py), never inferred from executed tests:

```text
DataColumnSidecar, DataColumnsByRootIdentifier, Eth1Block,
LightClientBootstrap, LightClientFinalityUpdate, LightClientHeader,
LightClientOptimisticUpdate, LightClientUpdate, MatrixEntry,
NewPayloadRequest, PartialDataColumnGroupID, PartialDataColumnPartsMetadata,
PartialDataColumnSidecar, PowBlock, SigningData
```

These types were already untested. Their exclusion is not a claim that their
production types are absent. The beta.2 archive contains 365 Gloas cases across
73 types: 75 cases across these 15 types are excluded, leaving 290 required cases
across 58 types. The log records declared exclusion names and included/excluded
case counts. New or renamed types are required unless separately approved as
exclusions; they cannot silently disappear from coverage.

## CI

- **C1:** `consensus-spec-vectors` verifies the exact Gloas declaration on pull
  requests, manual dispatch, and pushes to `master`. Superseded PR runs are
  cancelled; master and manual runs are not.
- **C2:** `consensus-spec-release` checks upstream every Friday at 10:00 UTC,
  including alpha, beta and rc releases. It selects the highest newer version by
  precedence, not publication date, and tests only `CONSENSUS_SPEC_FORK`, initially
  `gloas`. That fork must have a support declaration.
- **C3:** A newer release opens a bot-owned review issue with the tested commit,
  result, log artifact, release notes, cumulative supported-to-candidate comparison,
  checklist and target-specific reproduction instructions. Passing never promotes
  support automatically.
- **C4:** Exact bot-owned fork/release markers are checked before testing and again
  before writing. Only unambiguous `pass` or `fail` evidence suppresses retesting,
  whether the issue is open or closed. `not-run`, absent and ambiguous legacy
  verdicts remain retryable. Retries update the same issue's evidence section,
  preserving human notes, checklist progress and open/closed state. Persistent
  coverage gaps can therefore retry weekly until fixed or superseded by a newer release.
  Cancelled runs do not create issues.

Scheduled runs and manual dispatch require these workflows on `master`, the
repository's [default branch](https://docs.github.com/en/actions/how-tos/manage-workflow-runs/manually-run-a-workflow).
Manual dispatch can select another branch with `--ref`. PR checks do not require
that default-branch placement. Vector caching and paths exclusions are deferred.

Test jobs have read-only permissions; issue writing runs in a separate job.
The test subprocess does not inherit GitHub credentials or writable Actions
file-command variables. This is environment hygiene, not a sandbox for arbitrary code.

Set or inspect the watcher variable with:

```sh
gh variable set CONSENSUS_SPEC_FORK --body gloas
gh variable get CONSENSUS_SPEC_FORK
```

## Local verification

Requires Go, Python 3.12+, GitHub CLI authenticated for public API access, and curl.
The runner uses native `go test -json`; CI selects Go from the checked-out `go.mod`.
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

The runner downloads the exact release's `mainnet.tar.gz`, verifies GitHub's
published SHA-256 digest, and extracts only the selected fork's `ssz_static`
files. Allow approximately 1 GB for the download, plus extracted files.
CI retains the log artifact. The local snippet removes temporary files on exit;
redirect stdout and stderr to a file to keep evidence.

The release-review issue includes reproduction for arbitrary targets. It fetches
the recorded automation commit explicitly and copies its runner before creating
an owned temporary detached worktree from the fetched target. This supports targets
without CI files and preserves the caller's branch, HEAD and dirty files. Cleanup
removes the worktree registration and temporary files on success or failure.

Use a branch, tag or full commit SHA, not an abbreviated SHA. For local-only
targets, use `TARGET_REMOTE=.` and resolve abbreviated SHAs with `git rev-parse`
first. GitHub PR targets can use `refs/pull/<number>/head` from origin. The snippet
prints the tested commit. The automation commit must remain available from origin;
if GitHub no longer serves it, restore it before reproducing.

After release review and any fixes, change only the verified fork's version through
a PR. The baseline workflow verifies the new declaration; the watcher never updates it.
