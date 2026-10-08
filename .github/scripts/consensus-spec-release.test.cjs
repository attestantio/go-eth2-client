const assert = require('node:assert/strict');
const { test } = require('node:test');
const { execFileSync } = require('node:child_process');
const { mkdtempSync, mkdirSync, writeFileSync, readFileSync, rmSync } = require('node:fs');
const { tmpdir } = require('node:os');
const { join } = require('node:path');
const { selectRelease, reviewIssue, alreadyReported, supportedVersion, issueBody } = require('./consensus-spec-release.cjs');

test('selects the newest prerelease rather than an older stable patch', () => {
  const releases = [
    { tag_name: 'v1.6.1', draft: false, prerelease: false },
    { tag_name: 'v1.7.0-beta.2', draft: false, prerelease: true },
    { tag_name: 'v1.7.0-beta.10', draft: false, prerelease: true },
  ];
  assert.equal(selectRelease(releases, 'v1.7.0-beta.1'), releases[2]);
});

test('recognizes an existing bot issue even when closed or renamed', () => {
  const issues = [{
    state: 'closed', title: 'Renamed by a maintainer', user: { login: 'github-actions[bot]' },
    body: 'Result: pass\n<!-- consensus-spec:gloas:v1.7.0-beta.3 -->',
  }];
  assert.equal(alreadyReported(issues, 'gloas', 'v1.7.0-beta.3'), true);
  assert.equal(alreadyReported(issues, 'gloas', 'v1.7.0-beta.4'), false);
  assert.equal(alreadyReported(issues, 'fulu', 'v1.7.0-beta.3'), false);
  assert.equal(alreadyReported([{ ...issues[0], user: { login: 'someone' } }], 'gloas', 'v1.7.0-beta.3'), false);
});

test('not-run and missing or ambiguous verdicts remain retryable in open and closed issues', () => {
  const marker = '<!-- consensus-spec:gloas:v1.7.0-beta.3 -->';
  for (const state of ['open', 'closed']) {
    for (const verdict of ['Result: not-run', 'Result: not-run (setup failed)', '',
      'Result: unknown', 'Result: pass\nResult: not-run']) {
      assert.equal(alreadyReported([{ state, user: { login: 'github-actions[bot]' },
        body: `${verdict}\n${marker}` }], 'gloas', 'v1.7.0-beta.3'), false, `${state}/${verdict}`);
    }
    for (const verdict of ['pass', 'fail']) {
      assert.equal(alreadyReported([{ state, user: { login: 'github-actions[bot]' },
        body: `Result: ${verdict}\n${marker}` }], 'gloas', 'v1.7.0-beta.3'), true);
    }
  }
});

test('retry evidence updates the same issue without losing human review notes or checklist progress', () => {
  const data = { fork: 'gloas', tag: 'v1.7.0-beta.3', supported: 'v1.7.0-beta.2',
    sha: 'abc123', runURL: 'https://example.com/run/1', result: 'not-run' };
  const marker = '<!-- consensus-spec:gloas:v1.7.0-beta.3 -->';
  for (const previousBody of [issueBody(data), `Result: not-run\n- [x] Reviewed\n${marker}`]) {
    const previous = `${previousBody}\nHuman review note.`.replace('- [ ] Review the release', '- [x] Review the release');
    const existing = { number: 42, state: 'closed', user: { login: 'github-actions[bot]' }, body: previous };
    assert.equal(reviewIssue([existing], data.fork, data.tag), existing);
    const body = issueBody({ ...data, result: 'pass', runURL: 'https://example.com/run/2', previousBody: previous });
    assert.ok(body.includes('Human review note.'));
    assert.ok(body.includes('- [x] Review'));
    assert.ok(body.includes('run/2'));
    assert.equal(alreadyReported([{ ...existing, body }], data.fork, data.tag), true);
    assert.equal(reviewIssue([{ ...existing, body, pull_request: {} }], data.fork, data.tag), undefined);
  }
});

test('requires a declared version for the configured fork', () => {
  assert.equal(supportedVersion({ gloas: 'v1.7.0-beta.2' }, 'gloas'), 'v1.7.0-beta.2');
  assert.throws(() => supportedVersion({ gloas: 'v1.7.0-beta.2' }, 'fulu'), /No supported version/);
  assert.throws(() => supportedVersion({}, '../gloas'), /Invalid fork/);
  assert.throws(() => supportedVersion({ gloas: 'latest' }, 'gloas'), /Invalid supported version/);
  assert.throws(() => supportedVersion({ gloas: 'v1.7.0-beta.2\n' }, 'gloas'), /Invalid supported version/);
});

test('release discovery and the vector runner accept the same tag and fork examples', () => {
  const examples = [
    ['v1.7.0', 'gloas', true], ['v1.7.0-alpha.1', 'gloas', true],
    ['v1.7.0-beta.2', 'gloas', true], ['v1.7.0-rc.1', 'fork2', true],
    ['latest', 'gloas', false], ['v1.7.0-preview.1', 'gloas', false],
    ['v1.7.0-beta.2\n', 'gloas', false], ['v١.7.0', 'gloas', false],
    ['v1.7.0', '../gloas', false], ['v1.7.0', 'Gloas', false],
    ['v1.7.0', 'gloas\n', false], ['v1.7.0', '', false],
  ];
  const python = execFileSync('python3', ['-c', `
import json, sys
from consensus_vectors import validate_inputs
results = []
for tag, fork, _ in json.load(sys.stdin):
    try:
        validate_inputs(tag, fork)
        results.append(True)
    except ValueError:
        results.append(False)
print(json.dumps(results))
`], { cwd: __dirname, input: JSON.stringify(examples), encoding: 'utf8',
    env: { ...process.env, PYTHONDONTWRITEBYTECODE: '1' } });
  const results = JSON.parse(python);
  examples.forEach(([tag, fork, expected], index) => {
    let accepted = true;
    try { supportedVersion({ [fork]: tag }, fork); } catch { accepted = false; }
    assert.equal(accepted, expected, `discovery: ${tag}/${fork}`);
    assert.equal(results[index], expected, `runner: ${tag}/${fork}`);
  });
});

test('issue contains evidence, arbitrary-target instructions and a stable deduplication marker', () => {
  const body = issueBody({ fork: 'gloas', tag: 'v1.7.0-beta.3', supported: 'v1.7.0-beta.2',
    sha: 'abc123', runURL: 'https://github.com/attestantio/go-eth2-client/actions/runs/1', result: 'fail' });
  for (const expected of ['v1.7.0-beta.3', 'v1.7.0-beta.2', 'abc123', 'actions/runs/1',
    'Result: fail', 'compare/v1.7.0-beta.2...v1.7.0-beta.3',
    'after declared type exclusions', 'exclusion names and case counts',
    'git worktree add --detach', 'TARGET', 'full commit SHA', 'TARGET_REMOTE=.',
    'consensus_vectors.py', 'GOTOOLCHAIN=auto python3', '(\nset -eu',
    'native `go test -json`', 'redirect stdout and stderr',
    '<!-- consensus-spec:gloas:v1.7.0-beta.3 -->', 'not full consensus-spec conformance', 'GitHub CLI', 'curl']) {
    assert.ok(body.includes(expected), expected);
  }
  assert.ok(!body.includes('gosilent'));
  assert.ok(!body.includes('From a clean go-eth2-client clone'));
});

test('reproduction fetches the automation commit and tests the fetched target rather than a stale branch', () => {
  const root = mkdtempSync(join(tmpdir(), 'consensus-reproduction-'));
  const env = { ...process.env, GIT_CONFIG_GLOBAL: '/dev/null', GIT_CONFIG_NOSYSTEM: '1',
    GIT_AUTHOR_NAME: 'Test', GIT_AUTHOR_EMAIL: 'test@example.com',
    GIT_COMMITTER_NAME: 'Test', GIT_COMMITTER_EMAIL: 'test@example.com' };
  const git = (cwd, ...args) => execFileSync('git', ['-c', 'core.hooksPath=/dev/null', ...args],
    { cwd, env, encoding: 'utf8', stdio: ['pipe', 'pipe', 'pipe'] }).trim();
  try {
    const origin = join(root, 'origin');
    const target = join(root, 'target');
    const fork = join(root, 'fork');
    mkdirSync(origin);
    git(origin, 'init', '-b', 'feature');
    writeFileSync(join(origin, 'README'), 'initial');
    git(origin, 'add', '.');
    git(origin, 'commit', '-m', 'initial');
    git(root, 'clone', '--no-local', '--single-branch', origin, target);
    git(root, 'clone', '--no-local', '--single-branch', origin, fork);
    writeFileSync(join(fork, 'README'), 'fork-only change');
    git(fork, 'commit', '-am', 'fork target');
    const forkSHA = git(fork, 'rev-parse', 'HEAD');

    git(origin, 'checkout', '-b', 'automation');
    mkdirSync(join(origin, '.github/scripts'), { recursive: true });
    writeFileSync(join(origin, '.github/scripts/consensus_vectors.py'), 'print("runner")');
    git(origin, 'add', '.');
    git(origin, 'commit', '-m', 'automation');
    const sha = git(origin, 'rev-parse', 'HEAD');
    git(origin, 'checkout', 'feature');
    git(origin, 'branch', '-D', 'automation');
    writeFileSync(join(origin, 'README'), 'updated origin');
    git(origin, 'commit', '-am', 'origin target');
    const originSHA = git(origin, 'rev-parse', 'HEAD');
    git(target, 'checkout', '-b', 'local-only');
    writeFileSync(join(target, 'README'), 'local-only change');
    git(target, 'commit', '-am', 'local target');
    const localSHA = git(target, 'rev-parse', 'HEAD');
    writeFileSync(join(target, 'README'), 'dirty caller file');
    writeFileSync(join(target, 'untracked'), 'caller-owned data');
    const originalBranch = git(target, 'symbolic-ref', 'HEAD');
    const originalStatus = git(target, 'status', '--porcelain');

    for (const [remote, ref, expected, content] of [[fork, 'feature', forkSHA, 'fork-only change'],
      ['origin', 'feature', originSHA, 'updated origin'], ['.', localSHA, localSHA, 'local-only change']]) {
      const body = issueBody({ fork: 'gloas', tag: 'v1.7.0-beta.3', supported: 'v1.7.0-beta.2',
        sha, runURL: 'https://example.com/run', result: 'fail' });
      const snippet = /```sh\n([\s\S]*?)```/.exec(body)[1]
        .replace("TARGET='<branch-tag-or-commit>'", `TARGET='${ref}'`)
        .replace("TARGET_REMOTE='<repository-url-or-origin>'", `TARGET_REMOTE='${remote}'`);
      // Exercise the real Git preflight without downloading vectors or running Go.
      for (const failure of [false, true]) {
        const preflight = snippet.split('GOTOOLCHAIN=auto python3')[0]
          + `test "$(cat "$runner/consensus_vectors.py")" = 'print("runner")'\n`
          + `test "$(cat README)" = '${content}'\n${failure ? 'exit 1\n' : ''})\n`;
        let output;
        try {
          output = execFileSync('sh', ['-c', preflight],
            { cwd: target, env, encoding: 'utf8', stdio: ['pipe', 'pipe', 'pipe'] });
          assert.equal(failure, false);
        } catch (error) {
          assert.equal(failure, true);
          assert.equal(error.status, 1);
          output = error.stdout;
        }
        assert.equal(git(target, 'rev-parse', 'HEAD'), localSHA);
        assert.equal(git(target, 'symbolic-ref', 'HEAD'), originalBranch);
        assert.equal(git(target, 'status', '--porcelain'), originalStatus);
        assert.equal(git(target, 'worktree', 'list', '--porcelain').split('\n').filter(line => line.startsWith('worktree ')).length, 1);
        assert.ok(output.includes(expected), 'prints the tested commit');
      }
    }
  } finally {
    rmSync(root, { recursive: true, force: true });
  }
});

test('workflow policy cancels superseded PR runs without cancelling master or manual runs', () => {
  const workflow = readFileSync(join(__dirname, '../workflows/consensus-spec-vectors.yml'), 'utf8');
  assert.ok(workflow.includes("cancel-in-progress: ${{ github.event_name == 'pull_request' }}"));
  assert.ok(workflow.includes("group: ${{ github.workflow }}-${{ github.event_name == 'pull_request' && format('pr-{0}', github.event.pull_request.number) || github.run_id }}"));
});

test('ignores drafts and never proposes an equal or older version', () => {
  assert.equal(selectRelease([
    { tag_name: 'v1.8.0', draft: true },
    { tag_name: 'v1.7.0-beta.2', draft: false },
    { tag_name: 'v1.7.0-alpha.14', draft: false },
  ], 'v1.7.0-beta.2'), null);
});
