import unittest
import jenkins
import os
from unittest.mock import MagicMock, call

os.environ['MINUTES'] = "0.1"
os.environ['SERVICE'] = "service"
os.environ['NAMESPACE'] = "namespace"
os.environ['USER'] = "user"
os.environ['PASSWORD'] = "password"
load_simulator = __import__("jenkins-load-simulator")


class TestJobSimulator(unittest.TestCase):
    def setUp(self):
        self.server = MagicMock()
        pass
    
    def test_load_simulator(self):
        load_simulator.job_simulator(self.server)
        self.server.create_job.assert_called_with(
            "empty",  jenkins.EMPTY_CONFIG_XML
        )
        self.assertEqual(self.server.create_job.call_count, 1)
        
        self.assertEqual(self.server.jobs_count.call_count, 2)
        self.assertEqual(self.server.get_jobs.call_count, 2)
        
        self.server.get_job_config.assert_called_with("empty")
        self.assertEqual(self.server.get_job_config.call_count, 1)
        
        self.server.copy_job.assert_called_with("empty", "empty_copy")
        self.server.enable_job.assert_called_with("empty_copy")
        self.assertEqual(self.server.copy_job.call_count, 1)
        self.assertEqual(self.server.enable_job.call_count, 1)
        
        self.server.delete_job.has_calls([call("empty"), call("empty_copy")])
        

if __name__ == '__main__':
    unittest.main()
