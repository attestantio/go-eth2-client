from contextlib import chdir, redirect_stdout
import hashlib
import io
import json
import os
import subprocess
import sys
from pathlib import Path
import tarfile
import tempfile
import unittest

from consensus_vectors import validate_inputs, extract_vectors, covered_cases, classify_result, select_asset, verify_digest, go_test_environment, run_vectors


class ConsensusVectorsTest(unittest.TestCase):
    def test_rejects_invalid_release_or_fork(self):
        validate_inputs('v1.7.0-beta.2', 'gloas')
        for tag, fork in [('latest', 'gloas'), ('v1.7.0-beta.2', '../gloas'),
                          ('v1.7.0;echo bad', 'gloas'), ('v1.7.0', '')]:
            with self.subTest(tag=tag, fork=fork), self.assertRaises(ValueError):
                validate_inputs(tag, fork)

    def test_extracts_only_the_selected_forks_static_vectors(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            archive = root / 'vectors.tar.gz'
            with tarfile.open(archive, 'w:gz') as tar:
                for fork in ['gloas', 'fulu']:
                    for name in ['value.yaml', 'roots.yaml', 'serialized.ssz_snappy']:
                        member = tarfile.TarInfo(f'tests/mainnet/{fork}/ssz_static/BeaconState/ssz_random/case_0/{name}')
                        member.size = 6
                        tar.addfile(member, io.BytesIO(b'vector'))
            destination = root / 'extracted'
            self.assertEqual(extract_vectors(archive, 'gloas', destination), {'TestConsensusSpec/BeaconState/case_0'})
            self.assertEqual(len(list(destination.rglob('*.*'))), 3)
            self.assertFalse((destination / 'tests/mainnet/fulu').exists())
            self.assertEqual((destination / 'tests/mainnet/gloas/ssz_static/BeaconState/ssz_random/case_0/value.yaml').read_bytes(), b'vector')

    def test_rejects_missing_vectors_traversal_and_links(self):
        for name, kind in [
            ('tests/mainnet/fulu/ssz_static/value.yaml', tarfile.REGTYPE),
            ('tests/mainnet/gloas/ssz_static/../../escape', tarfile.REGTYPE),
            ('tests/mainnet/gloas/ssz_static/link', tarfile.SYMTYPE),
            ('tests/mainnet/gloas/ssz_static/link', tarfile.LNKTYPE),
        ]:
            with self.subTest(name=name, kind=kind), tempfile.TemporaryDirectory() as directory:
                archive = Path(directory) / 'vectors.tar.gz'
                with tarfile.open(archive, 'w:gz') as tar:
                    member = tarfile.TarInfo(name)
                    member.type = kind
                    member.linkname = '/etc/passwd' if kind != tarfile.REGTYPE else ''
                    tar.addfile(member)
                with self.assertRaises(ValueError):
                    extract_vectors(archive, 'gloas', Path(directory) / 'extracted')

    def test_rejects_incomplete_or_unrecognized_case_layouts(self):
        for case, files in [
            ('Fork/ssz_random/case_0', ['value.yaml']),
            ('Fork/ssz_other/case_0', ['value.yaml', 'roots.yaml', 'serialized.ssz_snappy']),
            ('Fork/ssz_random/case 0', ['value.yaml', 'roots.yaml', 'serialized.ssz_snappy']),
        ]:
            with self.subTest(case=case), tempfile.TemporaryDirectory() as directory:
                archive = Path(directory) / 'vectors.tar.gz'
                with tarfile.open(archive, 'w:gz') as tar:
                    for name in files:
                        tar.addfile(tarfile.TarInfo(f'tests/mainnet/gloas/ssz_static/{case}/{name}'))
                with self.assertRaises(ValueError):
                    extract_vectors(archive, 'gloas', Path(directory) / 'extracted')

    def test_classifies_json_terminal_events_not_printed_test_output(self):
        for code, statuses, expected in [
            (0, [('pass', 'TestConsensusSpec/Fork/case_0'), ('pass', 'TestConsensusSpec'), ('pass', '')], 'pass'),
            (1, [('fail', 'TestConsensusSpec/Fork/case_0'), ('fail', 'TestConsensusSpec'), ('fail', '')], 'fail'),
            (0, [('skip', 'TestConsensusSpec'), ('pass', '')], 'not-run'),
            (0, [('pass', 'TestConsensusSpec'), ('pass', '')], 'not-run'),
            (1, [('fail', '')], 'not-run'),
            (1, [('pass', 'TestConsensusSpec/Fork/case_0'), ('fail', 'TestConsensusSpec'), ('fail', '')], 'not-run'),
            (1, [('fail', 'TestConsensusSpec/Fork/case_0'), ('skip', 'TestConsensusSpec/Fork/case_1'),
                 ('fail', 'TestConsensusSpec'), ('fail', '')], 'fail'),
            (0, [('pass', 'TestConsensusSpec/Fork/case_0'), ('skip', 'TestConsensusSpec/Fork/case_1'),
                 ('pass', 'TestConsensusSpec'), ('pass', '')], 'not-run'),
            (0, [('pass', 'TestConsensusSpec/Fork/case_0'), ('pass', 'TestConsensusSpec')], 'not-run'),
            (0, [('output', 'TestConsensusSpec'), ('pass', 'TestConsensusSpec'), ('pass', '')], 'not-run'),
        ]:
            log = '\n'.join(json.dumps({'Package': 'gloas', 'Action': action, 'Test': name,
                                       'Output': '--- PASS: TestConsensusSpec/Fork/case_0'})
                            for action, name in statuses)
            with self.subTest(code=code, statuses=statuses):
                self.assertEqual(classify_result(code, log, {'TestConsensusSpec/Fork/case_0'}), expected)
        self.assertEqual(classify_result(0, 'malformed JSON', {'TestConsensusSpec/Fork/case_0'}), 'not-run')

    def test_pass_requires_every_inventory_case_exactly_once(self):
        case = 'TestConsensusSpec/Fork/case_0'
        for executed, inventory in [
            ([case], {case, 'TestConsensusSpec/FutureContainer/case_0'}),
            ([case], {'TestConsensusSpec/RenamedFork/case_0'}),
            ([case, case], {case}),
            ([case, 'TestConsensusSpec/Unexpected/case_0'], {case}),
            ([case], set()),
        ]:
            events = [{'Package': 'gloas', 'Action': 'pass', 'Test': name}
                      for name in [*executed, 'TestConsensusSpec', '']]
            log = '\n'.join(map(json.dumps, events))
            with self.subTest(executed=executed, inventory=inventory):
                self.assertEqual(classify_result(0, log, inventory), 'not-run')

    def test_only_declared_exclusions_are_removed_from_required_coverage(self):
        excluded_types = {'DataColumnSidecar', 'DataColumnsByRootIdentifier', 'Eth1Block',
                          'LightClientBootstrap', 'LightClientFinalityUpdate', 'LightClientHeader',
                          'LightClientOptimisticUpdate', 'LightClientUpdate', 'MatrixEntry',
                          'NewPayloadRequest', 'PartialDataColumnGroupID', 'PartialDataColumnPartsMetadata',
                          'PartialDataColumnSidecar', 'PowBlock', 'SigningData'}
        included = {'TestConsensusSpec/Fork/case_0', 'TestConsensusSpec/FutureContainer/case_0'}
        all_cases = included | {f'TestConsensusSpec/{name}/case_0' for name in excluded_types}
        output = io.StringIO()
        with redirect_stdout(output):
            self.assertEqual(covered_cases(all_cases, 'gloas'), included)
        self.assertIn('15 excluded cases', output.getvalue())
        for name in excluded_types:
            self.assertIn(name, output.getvalue())
        self.assertEqual(covered_cases(all_cases, 'fulu'), all_cases)
        self.assertEqual(len(all_cases), 17)

    def test_invalid_or_cross_package_events_cannot_prove_execution(self):
        case = 'TestConsensusSpec/Fork/case_0'
        for events in [
            [None], [[]],
            [{'Action': 'pass', 'Test': name} for name in [case, 'TestConsensusSpec', '']],
            [{'Package': 'gloas', 'Action': 'pass', 'Test': 123}],
            [{'Package': 'other', 'Action': 'pass', 'Test': case},
             {'Package': 'gloas', 'Action': 'pass', 'Test': 'TestConsensusSpec'},
             {'Package': 'gloas', 'Action': 'pass'}],
        ]:
            with self.subTest(events=events):
                self.assertEqual(classify_result(0, '\n'.join(map(json.dumps, events)), {case}), 'not-run')

    def test_requires_mainnet_asset_with_published_checksum_and_official_url(self):
        asset = {'name': 'mainnet.tar.gz', 'digest': 'sha256:' + 'a' * 64,
                 'browser_download_url': 'https://github.com/ethereum/consensus-specs/releases/download/v1.7.0-beta.2/mainnet.tar.gz'}
        self.assertEqual(select_asset({'assets': [asset]}), asset)
        for assets in [[], [{**asset, 'digest': None}], [{**asset, 'browser_download_url': 'https://example.com/archive'}]]:
            with self.subTest(assets=assets), self.assertRaises(ValueError):
                select_asset({'assets': assets})

    def test_compares_archive_bytes_with_upstreams_digest(self):
        with tempfile.TemporaryDirectory() as directory:
            archive = Path(directory) / 'archive'
            archive.write_bytes(b'release vector bytes')
            digest = 'sha256:' + hashlib.sha256(archive.read_bytes()).hexdigest()
            verify_digest(archive, digest)
            archive.write_bytes(b'corrupted bytes')
            with self.assertRaisesRegex(ValueError, 'checksum mismatch'):
                verify_digest(archive, digest)

    def test_test_process_gets_vectors_but_no_github_credentials_or_output_file(self):
        environment = {'PATH': '/usr/bin', 'GH_TOKEN': 'secret', 'GITHUB_TOKEN': 'secret',
                       'ACTIONS_RUNTIME_TOKEN': 'secret', 'ACTIONS_ID_TOKEN_REQUEST_TOKEN': 'secret',
                       'GITHUB_OUTPUT': '/tmp/output', 'GITHUB_ENV': '/tmp/env', 'GITHUB_PATH': '/tmp/path',
                       'GITHUB_STEP_SUMMARY': '/tmp/summary', 'GITHUB_STATE': '/tmp/state'}
        result = go_test_environment(environment, Path('/tmp/vectors'))
        self.assertEqual(result, {'PATH': '/usr/bin', 'CONSENSUS_SPEC_TESTS_DIR': str(Path('/tmp/vectors').resolve())})
        self.assertIn('GH_TOKEN', environment)

    def test_go_command_diagnostics_do_not_enter_the_json_stream(self):
        with tempfile.TemporaryDirectory() as directory, chdir(directory):
            Path('go.mod').write_text('module example.com/target\n\ngo 1.20\n')
            Path('spec/gloas').mkdir(parents=True)
            process = run_vectors('gloas', Path(directory) / 'vectors')
            self.assertNotEqual(process.returncode, 0)
            self.assertIsInstance(process.stderr, str)
            self.assertIn('matched no packages', process.stderr)
            self.assertNotIn('matched no packages', process.stdout)
            self.assertEqual(classify_result(process.returncode, process.stdout, set()), 'not-run')

    def test_cli_rejects_invalid_inputs_and_reports_no_conformance_result(self):
        with tempfile.TemporaryDirectory() as directory:
            output = Path(directory) / 'output'
            process = subprocess.run([sys.executable, str(Path(__file__).with_name('consensus_vectors.py')),
                                      'latest', 'gloas', directory], capture_output=True, text=True,
                                     env={**os.environ, 'GITHUB_OUTPUT': str(output)})
            self.assertNotEqual(process.returncode, 0)
            self.assertIn('Invalid consensus-spec release tag', process.stdout + process.stderr)
            self.assertIn('result=not-run', output.read_text())


if __name__ == '__main__':
    unittest.main()
