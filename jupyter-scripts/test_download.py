import sys
import os
import unittest
import requests
import generate_io


class MockedRequests:

    def get(self, url):
        return MockedResponse(url)


class MockedResponse:

    def __init__(self, url):
        self.content = str.encode(url)


class TestReader(unittest.TestCase):

    def test_download(self):
        test_url = "someurl"
        test_file = "testdownloadfile"
        sys.modules["requests"] = MockedRequests()

        generate_io.download(test_url, test_file)

        sys.modules["requests"] = requests

        with open(test_file) as fin:
            download = fin.read()
        os.remove(test_file)

        self.assertEqual(test_url, download)


if __name__ == '__main__':
    unittest.main()
