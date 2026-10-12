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


# These types are outside the declared Gloas vector coverage, never inferred from a run.
EXCLUDED_TYPES = {
    'gloas': {
        'DataColumnSidecar', 'DataColumnsByRootIdentifier', 'Eth1Block',
        'LightClientBootstrap', 'LightClientFinalityUpdate', 'LightClientHeader',
        'LightClientOptimisticUpdate', 'LightClientUpdate', 'MatrixEntry',
        'NewPayloadRequest', 'PartialDataColumnGroupID', 'PartialDataColumnPartsMetadata',
        'PartialDataColumnSidecar', 'PowBlock', 'SigningData',
    },
}


def validate_inputs(tag, fork):
    if not re.fullmatch(r'v[0-9]+\.[0-9]+\.[0-9]+(?:-(?:alpha|beta|rc)\.[0-9]+)?', tag):
        raise ValueError('Invalid consensus-spec release tag')
    if not re.fullmatch(r'[a-z][a-z0-9]*', fork):
        raise ValueError('Invalid consensus fork')


def extract_vectors(archive, fork, destination):
    prefix = f'tests/mainnet/{fork}/ssz_static/'
    cases = {}
    with tarfile.open(archive, 'r|gz') as tar:
        for member in tar:
            if not member.name.startswith(prefix) or member.isdir():
                continue
            if not member.isfile() or '..' in Path(member.name).parts:
                raise ValueError(f'Unsafe vector archive member: {member.name}')
            parts = Path(member.name[len(prefix):]).parts
            if (len(parts) != 4 or parts[1] != 'ssz_random'
                    or not all(re.fullmatch(r'[A-Za-z0-9_-]+', name) for name in (parts[0], parts[2]))):
                raise ValueError(f'Unrecognized vector case layout: {member.name}')
            case = f'TestConsensusSpec/{parts[0]}/{parts[2]}'
            cases.setdefault(case, set()).add(parts[3])
            path = Path(destination) / member.name
            path.parent.mkdir(parents=True, exist_ok=True)
            with tar.extractfile(member) as source, path.open('wb') as target:
                shutil.copyfileobj(source, target)
    if not cases:
        raise ValueError(f'No mainnet static vectors found for {fork}')
    required = {'value.yaml', 'roots.yaml', 'serialized.ssz_snappy'}
    for case, files in cases.items():
        if not required <= files:
            raise ValueError(f'Incomplete vector case {case}: missing {sorted(required - files)}')
    return set(cases)


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
                'ACTIONS_ID_TOKEN_REQUEST_TOKEN', 'GITHUB_OUTPUT', 'GITHUB_ENV',
                'GITHUB_PATH', 'GITHUB_STEP_SUMMARY', 'GITHUB_STATE'}
    return {**{key: value for key, value in environment.items() if key not in excluded},
            'CONSENSUS_SPEC_TESTS_DIR': str(Path(destination).resolve())}


def covered_cases(cases, fork):
    excluded_types = EXCLUDED_TYPES.get(fork, set())
    included = {case for case in cases if case.split('/')[1] not in excluded_types}
    print(f'Coverage for {fork}: {len(included)} required cases; {len(cases - included)} excluded cases')
    print(f'Declared excluded types: {sorted(excluded_types)}')
    return included


def classify_result(returncode, log, expected_cases):
    try:
        events = [json.loads(line) for line in log.splitlines()]
    except json.JSONDecodeError:
        return 'not-run'
    if not all(isinstance(event, dict) and isinstance(event.get('Action'), str)
               and isinstance(event.get('Test', ''), str) for event in events):
        return 'not-run'
    terminal = [(event.get('Package'), event.get('Test', ''), event['Action'])
                for event in events if event.get('Action') in {'pass', 'fail', 'skip'}]
    roots = [(package, action) for package, name, action in terminal if name == 'TestConsensusSpec']
    if len(roots) != 1:
        return 'not-run'
    package, root = roots[0]
    if not isinstance(package, str) or not package:
        return 'not-run'
    cases = [(name, action) for owner, name, action in terminal
             if owner == package and name.startswith('TestConsensusSpec/') and name.count('/') == 2]
    executed = {name for name, _ in cases}
    if expected_cases - executed:
        print(f'Unexecuted vector cases: {sorted(expected_cases - executed)}')
    if executed - expected_cases:
        print(f'Unexpected vector cases: {sorted(executed - expected_cases)}')
    if returncode != 0 and root == 'fail' and any(
            name in expected_cases and action == 'fail' for name, action in cases):
        return 'fail'
    if any(action == 'skip' for owner, name, action in terminal
           if owner == package and (name == 'TestConsensusSpec' or name.startswith('TestConsensusSpec/'))):
        return 'not-run'
    if (returncode == 0 and root == 'pass' and expected_cases
            and executed == expected_cases and len(cases) == len(expected_cases)
            and all(action == 'pass' for _, action in cases) and (package, '', 'pass') in terminal):
        return 'pass'
    return 'not-run'


def run_vectors(fork, destination):
    return subprocess.run([
        'go', 'test', '-json', '-count=1', '-timeout=15m',
        '-run', '^TestConsensusSpec$', f'./spec/{fork}/...',
    ], stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True,
        env=go_test_environment(os.environ, destination))


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
    cases = extract_vectors(archive, args.fork, args.destination)
    archive.unlink()
    print(f'Verified {args.tag} SHA-256; extracted {len(cases)} cases for {args.fork}', flush=True)
    cases = covered_cases(cases, args.fork)
    process = run_vectors(args.fork, args.destination)
    print(process.stderr, end='', file=sys.stderr)
    print(process.stdout, end='')
    result = classify_result(process.returncode, process.stdout, cases)
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
