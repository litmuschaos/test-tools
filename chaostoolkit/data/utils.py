import json
import subprocess
import sys
import os
import time
import requests
import logging
import threading
from jinja2 import Environment, FileSystemLoader, select_autoescape
import yaml
logger = logging.getLogger(__name__)


class Utils(object):

    ####################################
    #      Function definitions        #
    ####################################
   
    """
    run_shell_task() runs a shell command and prints the output as it executes. 
    It takes a list of strings that comprises the command itself, as the sole arg.
    """
    def run_shell_task(self,cmd_arg_list):
        run_cmd = subprocess.Popen(cmd_arg_list, stdout=subprocess.PIPE, env=os.environ.copy())
        run_cmd.communicate()

    """
    chaos_result_tracker() creates/patches the litmus chaosresult custom resource in the provided namespace.
    Typically invoked before and after chaos, and takes the .spec.phase, .spec.verdict & namespace as as args.
    """
    def chaos_result_tracker(self, exp_name, exp_phase, exp_verdict, ns):
        env_tmpl = Environment(loader = FileSystemLoader('./'), trim_blocks=True, lstrip_blocks=True, autoescape=select_autoescape(['yaml']))
        template = env_tmpl.get_template('chaos-result.j2')
        updated_chaosresult_template = template.render(c_experiment=exp_name, phase=exp_phase, verdict=exp_verdict)
        with open('chaosresult.yaml', "w+") as f:
            f.write(updated_chaosresult_template)
        chaosresult_update_cmd_args_list = ['kubectl','apply', '-f', 'chaosresult.yaml', '-n', ns]
        self.run_shell_task(chaosresult_update_cmd_args_list)


