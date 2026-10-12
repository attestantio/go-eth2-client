const evidencePattern = /<!-- consensus-spec-evidence:start -->[\s\S]*?<!-- consensus-spec-evidence:end -->/;

function parseVersion(tag) {
  const match = /^v(\d+)\.(\d+)\.(\d+)(?:-(alpha|beta|rc)\.(\d+))?$/.exec(tag);
  if (!match) return null;
  return [Number(match[1]), Number(match[2]), Number(match[3]),
    { alpha: 0, beta: 1, rc: 2 }[match[4]] ?? 3, Number(match[5] ?? 0)];
}

function compare(left, right) {
  for (let i = 0; i < left.length; i++) {
    if (left[i] !== right[i]) return left[i] - right[i];
  }
  return 0;
}

function selectRelease(releases, supported) {
  const baseline = parseVersion(supported);
  return releases.filter(release => !release.draft && parseVersion(release.tag_name))
    .filter(release => compare(parseVersion(release.tag_name), baseline) > 0)
    .sort((left, right) => compare(parseVersion(right.tag_name), parseVersion(left.tag_name)))[0] ?? null;
}

function reviewIssue(issues, fork, tag) {
  return issues.find(issue => !issue.pull_request && issue.user?.login === 'github-actions[bot]'
    && issue.body?.includes(`<!-- consensus-spec:${fork}:${tag} -->`));
}

function alreadyReported(issues, fork, tag) {
  const body = reviewIssue(issues, fork, tag)?.body ?? '';
  const verdicts = (body.match(evidencePattern)?.[0] ?? body).match(/^Result: .*$/gm) ?? [];
  return verdicts.length === 1 && ['Result: pass', 'Result: fail'].includes(verdicts[0]);
}

function supportedVersion(versions, fork) {
  if (!/^[a-z][a-z0-9]*$/.test(fork)) throw new Error('Invalid fork variable');
  if (!Object.hasOwn(versions, fork)) throw new Error(`No supported version declared for ${fork}`);
  if (!parseVersion(versions[fork])) throw new Error('Invalid supported version');
  return versions[fork];
}
function issueBody({ fork, tag, supported, sha, runURL, result, previousBody }) {
  const body = `## Review consensus-spec ${tag} support for ${fork}

<!-- consensus-spec-evidence:start -->
## Latest vector run

Declared supported version: ${supported}
Tested go-eth2-client commit: ${sha}
Result: ${result}

[Release notes](https://github.com/ethereum/consensus-specs/releases/tag/${tag})
[Changes since declared support](https://github.com/ethereum/consensus-specs/compare/${supported}...${tag})
[Vector run and log artifact](${runURL})

A pass requires complete static-case execution after declared type exclusions.
The log artifact records exclusion names and case counts. It verifies covered YAML/SSZ
encodings and hash-tree roots, not full consensus-spec conformance.
Review the release delta before updating support. A setup failure or missing execution
evidence is not a vector incompatibility result.

## Reproduce against a target

From a go-eth2-client clone, set TARGET to a branch, tag or full commit SHA.
Set TARGET_REMOTE to origin or a fork's repository URL. For local-only targets,
use TARGET_REMOTE=. and resolve abbreviated SHAs with git rev-parse first.
To test a GitHub PR through origin, use refs/pull/<number>/head as TARGET.
Copy the runner before creating a temporary worktree, so targets without these
CI files can also be tested without changing the caller's checkout.
Go, Python 3.12+, authenticated GitHub CLI and curl are required.
The mainnet download is approximately 1 GB.

\`\`\`sh
(
set -eu
TARGET='<branch-tag-or-commit>'
TARGET_REMOTE='<repository-url-or-origin>'
repo="$(git rev-parse --show-toplevel)"
runner="$(mktemp -d)"
cleanup() {
  git -C "$repo" worktree remove --force "$runner/target" >/dev/null 2>&1 || true
  rm -rf "$runner"
}
trap cleanup EXIT
cd "$repo"
git fetch origin ${sha}
git show FETCH_HEAD:.github/scripts/consensus_vectors.py > "$runner/consensus_vectors.py"
git fetch "$TARGET_REMOTE" "$TARGET"
git worktree add --detach "$runner/target" FETCH_HEAD
cd "$runner/target"
echo "Tested go-eth2-client commit: $(git rev-parse HEAD)"
GOTOOLCHAIN=auto python3 "$runner/consensus_vectors.py" ${tag} ${fork} "$runner/vectors"
)
\`\`\`

The runner uses native \`go test -json\`. Local temporary files are removed on subshell
exit; redirect stdout and stderr to a file with \`> log 2>&1\` to keep a local log.
The recorded automation commit must remain available from origin. If it was
removed and GitHub no longer serves it, restore that commit before reproducing.
<!-- consensus-spec-evidence:end -->

- [ ] Review the release changes relevant to this API library.
- [ ] Run the ${fork} mainnet static vectors against the intended target and record the commit/result.
- [ ] Resolve incompatibilities or explain why no library update is needed.
- [ ] Update only the ${fork} entry in .github/consensus-spec-versions.json after verification.

<!-- consensus-spec:${fork}:${tag} -->`;
  if (!previousBody) return body;
  const evidence = body.match(evidencePattern)[0];
  if (evidencePattern.test(previousBody)) return previousBody.replace(evidencePattern, () => evidence);
  return `${evidence}\n\n## Original issue content\n\n${previousBody}`;
}

module.exports = { selectRelease, reviewIssue, alreadyReported, supportedVersion, issueBody };
