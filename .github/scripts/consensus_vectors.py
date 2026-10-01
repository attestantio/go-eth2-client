import argparse
import hashlib
import json
import os
from pathlib import Path
import re
import shutil
import tarfile
import sys
import subprocess


def validate_inputs(tag, fork):
    if not re.fullmatch(r'v[0-9]+\.[0-9]+\.[0-9]+(?:-(?:alpha|beta|rc)\.[0-9]+)?', tag):
        raise ValueError('Invalid consensus-spec release tag')
    if not re.fullmatch(r'[a-z][a-z0-9]*', fork):
        raise ValueError('Invalid consensus fork')


def extract_vectors(archive, fork, destination):
    prefix = f'tests/mainnet/{fork}/ssz_static/'
    count = 0
    with tarfile.open(archive, 'r|gz') as tar:
        for member in tar:
            if not member.name.startswith(prefix) or member.isdir():
                continue
            if not member.isfile() or '..' in Path(member.name).parts:
                raise ValueError(f'Unsafe vector archive member: {member.name}')
            path = Path(destination) / member.name
            path.parent.mkdir(parents=True, exist_ok=True)
            with tar.extractfile(member) as source, path.open('wb') as target:
                shutil.copyfileobj(source, target)
            count += 1
    if not count:
        raise ValueError(f'No mainnet static vectors found for {fork}')
    return count


def select_asset(release):
    for asset in release['assets']:
        if asset['name'] == 'mainnet.tar.gz':
            if not re.fullmatch(r'sha256:[a-f0-9]{64}', asset.get('digest') or ''):
                raise ValueError('Mainnet asset has no published SHA-256 digest')
            if not asset['browser_download_url'].startswith('https://github.com/ethereum/consensus-specs/releases/download/'):
                raise ValueError('Unexpected vector download URL')
            return asset
    raise ValueError('Release has no mainnet.tar.gz asset')


def verify_digest(archive, digest):
    with Path(archive).open('rb') as source:
        actual = 'sha256:' + hashlib.file_digest(source, 'sha256').hexdigest()
    if actual != digest:
        raise ValueError('Vector archive checksum mismatch')


def go_test_environment(environment, destination):
    excluded = {'GH_TOKEN', 'GITHUB_TOKEN', 'ACTIONS_RUNTIME_TOKEN',
                'ACTIONS_ID_TOKEN_REQUEST_TOKEN', 'GITHUB_OUTPUT'}
    return {**{key: value for key, value in environment.items() if key not in excluded},
            'CONSENSUS_SPEC_TESTS_DIR': str(Path(destination).resolve())}


def classify_result(returncode, log):
    if not re.search(r'^\s*--- (PASS|FAIL): TestConsensusSpec/[^/]+/[^\s]+', log, re.MULTILINE):
        return 'not-run'
    if (returncode != 0 and re.search(r'^--- FAIL: TestConsensusSpec ', log, re.MULTILINE)
            and re.search(r'^\s*--- FAIL: TestConsensusSpec/[^/]+/[^\s]+', log, re.MULTILINE)):
        return 'fail'
    if re.search(r'^\s*--- SKIP: TestConsensusSpec', log, re.MULTILINE):
        return 'not-run'
    if returncode == 0 and re.search(r'^--- PASS: TestConsensusSpec ', log, re.MULTILINE):
        return 'pass'
    return 'not-run'


def main():
    if os.environ.get('GITHUB_OUTPUT'):
        with open(os.environ['GITHUB_OUTPUT'], 'a') as output:
            output.write('result=not-run\n')
    parser = argparse.ArgumentParser(description='Verify one fork against mainnet consensus-spec vectors')
    parser.add_argument('tag')
    parser.add_argument('fork')
    parser.add_argument('destination', type=Path)
    args = parser.parse_args()
    validate_inputs(args.tag, args.fork)
    if not Path(f'spec/{args.fork}').is_dir():
        raise ValueError(f'Target checkout has no spec/{args.fork} package')
    args.destination.mkdir(parents=True, exist_ok=False)
    release = json.loads(subprocess.check_output([
        'gh', 'api', f'repos/ethereum/consensus-specs/releases/tags/{args.tag}',
    ], text=True))
    asset = select_asset(release)
    archive = args.destination / 'mainnet.tar.gz'
    subprocess.run(['curl', '--fail', '--location', '--retry', '3', '--silent', '--show-error',
                    '--output', str(archive), asset['browser_download_url']], check=True)
    verify_digest(archive, asset['digest'])
    count = extract_vectors(archive, args.fork, args.destination)
    archive.unlink()
    print(f'Verified {args.tag} SHA-256; extracted {count} files for {args.fork}', flush=True)
    process = subprocess.run([
        'go', 'test', '-v', '-count=1', '-timeout=15m',
        '-run', '^TestConsensusSpec$', f'./spec/{args.fork}/...',
    ], stdout=subprocess.PIPE, stderr=subprocess.STDOUT, text=True,
        env=go_test_environment(os.environ, args.destination))
    print(process.stdout, end='')
    result = classify_result(process.returncode, process.stdout)
    if os.environ.get('GITHUB_OUTPUT'):
        with open(os.environ['GITHUB_OUTPUT'], 'a') as output:
            output.write(f'result={result}\n')
    print(f'Vector result: {result}')
    return 0 if result == 'pass' else 1


if __name__ == '__main__':
    try:
        sys.exit(main())
    except Exception as error:
        print(f'Vector verification setup failed: {error}', file=sys.stderr)
        sys.exit(1)
