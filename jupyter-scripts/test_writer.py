import sys
import io
import os
import unittest
import re
import generate_io as target


class TestWriter(unittest.TestCase):

    def test_when_file_name_present_then_it_is_used(self):
        filename = "writertestfile"

        writer = target.Writer(filename)

        self.assertEqual(writer.file_name, filename)

    def test_file_is_written_and_reported_to_stdout(self):
        test_file = "writertestfile"

        stdout = sys.stdout
        tmpout = io.StringIO()
        sys.stdout = tmpout

        writer = target.Writer(test_file)
        writer.run()
        sys.stdout = stdout

        reported = [re.sub(r"\D", "", v) for v in tmpout.getvalue().strip().split("\n")]
        tmpout.close()
        with open(test_file, "r") as fin:
            written = fin.read()
        os.remove(test_file)

        i = 0
        for r in reported:
            rlen = len(r)
            self.assertEqual(r, written[i:i + rlen])
            i += rlen


if __name__ == '__main__':
    unittest.main()
