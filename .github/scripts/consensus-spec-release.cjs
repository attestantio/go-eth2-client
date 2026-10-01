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

function alreadyReported(issues, fork, tag) {
  return issues.some(issue => !issue.pull_request && issue.user?.login === 'github-actions[bot]'
    && issue.body?.includes(`<!-- consensus-spec:${fork}:${tag} -->`));
}

function supportedVersion(versions, fork) {
  if (!/^[a-z][a-z0-9]*$/.test(fork)) throw new Error('Invalid fork variable');
  if (!Object.hasOwn(versions, fork)) throw new Error(`No supported version declared for ${fork}`);
  if (!parseVersion(versions[fork])) throw new Error('Invalid supported version');
  return versions[fork];
}
function issueBody({ fork, tag, supported, sha, runURL, result }) {
  return `## Review consensus-spec ${tag} support for ${fork}

Declared supported version: ${supported}
Tested go-eth2-client commit: ${sha}
Result: ${result}

[Release notes](https://github.com/ethereum/consensus-specs/releases/tag/${tag})
[Vector run and log artifact](${runURL})

A passing static-vector run verifies covered YAML/SSZ encodings and hash-tree roots,
not full consensus-spec conformance. Review the release delta before updating support.
A setup failure or missing execution evidence is not a vector incompatibility result.

## Reproduce against a target

From a clean go-eth2-client clone, set TARGET to a branch, tag or full commit SHA.
Set TARGET_REMOTE to origin or a fork's repository URL. For local-only targets,
use TARGET_REMOTE=. and resolve abbreviated SHAs with git rev-parse first.
To test a GitHub PR through origin, use refs/pull/<number>/head as TARGET.
Copy the runner before checkout so targets without these CI files can also be tested.
Go, Python 3.12+, authenticated GitHub CLI and curl are required.
The mainnet download is approximately 1 GB.

\`\`\`sh
(
set -eu
TARGET='<branch-tag-or-commit>'
TARGET_REMOTE='<repository-url-or-origin>'
runner="$(mktemp -d)"
trap 'rm -rf "$runner"' EXIT
git fetch origin ${sha}
git show FETCH_HEAD:.github/scripts/consensus_vectors.py > "$runner/consensus_vectors.py"
git fetch "$TARGET_REMOTE" "$TARGET"
git checkout --detach FETCH_HEAD
echo "Tested go-eth2-client commit: $(git rev-parse HEAD)"
GOTOOLCHAIN=auto python3 "$runner/consensus_vectors.py" ${tag} ${fork} "$runner/vectors"
)
\`\`\`

The runner uses normal \`go test\`. Local temporary files are removed on subshell
exit; redirect this snippet's output to a file if you need a permanent local log.
The recorded automation commit must remain available from origin. If it was
removed and GitHub no longer serves it, restore that commit before reproducing.

- [ ] Review the release changes relevant to this API library.
- [ ] Run the ${fork} mainnet static vectors against the intended target and record the commit/result.
- [ ] Resolve incompatibilities or explain why no library update is needed.
- [ ] Update only the ${fork} entry in .github/consensus-spec-versions.json after verification.

<!-- consensus-spec:${fork}:${tag} -->`;
}

module.exports = { selectRelease, alreadyReported, supportedVersion, issueBody };
