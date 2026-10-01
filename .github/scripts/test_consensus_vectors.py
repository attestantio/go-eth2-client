import hashlib
import io
import os
import subprocess
import sys
from pathlib import Path
import tarfile
import tempfile
import unittest

from consensus_vectors import validate_inputs, extract_vectors, classify_result, select_asset, verify_digest, go_test_environment


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
            self.assertEqual(extract_vectors(archive, 'gloas', destination), 3)
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

    def test_success_requires_executed_vector_cases_without_skips(self):
        for code, log, expected in [
            (0, '    --- PASS: TestConsensusSpec/BeaconState/case_0\n--- PASS: TestConsensusSpec (1s)\n', 'pass'),
            (1, '    --- FAIL: TestConsensusSpec/BeaconState/case_0\n--- FAIL: TestConsensusSpec (1s)\n', 'fail'),
            (0, '--- SKIP: TestConsensusSpec (0s)\n', 'not-run'),
            (0, '--- PASS: TestConsensusSpec (0s)\n', 'not-run'),
            (1, 'build failed\n', 'not-run'),
            (1, '    --- PASS: TestConsensusSpec/BeaconState/case_0\n    --- SKIP: TestConsensusSpec/BeaconState/case_1\n--- FAIL: TestConsensusSpec (1s)\n', 'not-run'),
            (1, '    --- PASS: TestConsensusSpec/BeaconState/case_0\n--- FAIL: TestConsensusSpec (1s)\n', 'not-run'),
            (1, '    --- FAIL: TestConsensusSpec/BeaconState/case_0\n    --- SKIP: TestConsensusSpec/BeaconState/case_1\n--- FAIL: TestConsensusSpec (1s)\n', 'fail'),
            (0, '    --- PASS: TestConsensusSpec/BeaconState/case_0\n    --- SKIP: TestConsensusSpec/BeaconState/case_1\n--- PASS: TestConsensusSpec (1s)\n', 'not-run'),
        ]:
            with self.subTest(code=code, log=log):
                self.assertEqual(classify_result(code, log), expected)

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
                       'GITHUB_OUTPUT': '/tmp/output'}
        result = go_test_environment(environment, Path('/tmp/vectors'))
        self.assertEqual(result, {'PATH': '/usr/bin', 'CONSENSUS_SPEC_TESTS_DIR': str(Path('/tmp/vectors').resolve())})
        self.assertIn('GH_TOKEN', environment)

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
