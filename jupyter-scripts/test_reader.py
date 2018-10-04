import sys
import io
import os
import unittest
import generate_io as target


class TestReader(unittest.TestCase):

    def test_when_file_name_present_then_it_is_used(self):
        filename = "readertestfile"

        reader = target.Reader(filename)

        self.assertEqual(reader.file_name, filename)

    def test_when_file_is_read_then_output_forwarded_to_stdout(self):
        test_data = ["a", "b", "c"]
        test_file = "readertestfile"
        with open(test_file, "w") as fout:
            for line in test_data:
                fout.write("%s\n" % line)
        stdout = sys.stdout
        tmpout = io.StringIO()
        sys.stdout = tmpout

        reader = target.Reader(test_file)
        reader.run()
        sys.stdout = stdout

        result = tmpout.getvalue()
        tmpout.close()
        os.remove(test_file)

        for actual, expected in zip(result.strip().split("\n\n"), test_data):
            self.assertRegex(actual, "%s$" % expected)


if __name__ == '__main__':
    unittest.main()
