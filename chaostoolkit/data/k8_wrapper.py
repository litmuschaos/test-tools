import subprocess
import sys
import os
import argparse
import logging
import json
from jinja2 import Environment, FileSystemLoader, select_autoescape
import yaml


__author__ = 'Sumit_Nagal@intuit.com'

####################################
#      Function definitions        #
####################################

"""
run_shell_task() runs a shell command and prints the output as it executes. 
It takes a list of strings that comprises the command itself, as the sole arg.
"""
def run_shell_task(cmd_arg_list):
    run_cmd = subprocess.Popen(cmd_arg_list, stdout=subprocess.PIPE, env=os.environ.copy())
    run_cmd.communicate()

"""
chaos_result_tracker() creates/patches the litmus chaosresult custom resource in the provided namespace.
Typically invoked before and after chaos, and takes the .spec.phase, .spec.verdict & namespace as as args.
"""
def chaos_result_tracker(exp_name, exp_phase, exp_verdict, ns):
    env_tmpl = Environment(loader = FileSystemLoader('./'), trim_blocks=True, lstrip_blocks=True, autoescape=select_autoescape(['yaml']))
    template = env_tmpl.get_template('chaos-result.j2')
    updated_chaosresult_template = template.render(c_experiment=exp_name, phase=exp_phase, verdict=exp_verdict)
    with open('chaosresult.yaml', "w+") as f:
        f.write(updated_chaosresult_template)
    chaosresult_update_cmd_args_list = ['kubectl','apply', '-f', 'chaosresult.yaml', '-n', ns]
    run_shell_task(chaosresult_update_cmd_args_list)

####################################
# Start of Python Chaos Experiment #
####################################

parser = argparse.ArgumentParser()

parser.add_argument("-file", action='store',
                    default="pod-app-kill-count.json",
                    dest = "file",
                    help="Chaos file to chose for execution"
                    )
parser.add_argument("-exp", action='store',
                    default="k8-pod-delete",
                    dest = "exp",
                    help="Chaos experiment to chose for execution"
                    )                    
parser.add_argument('-label', action='store',
                    dest='label',
                    default="app",
                    help='Store a label value')
parser.add_argument("-namespace", action='store',
                    default="default",
                    dest = "namespace",
                    help="namespace for application"
                    )
parser.add_argument('-app', action='store',
                    dest='app',
                    default="localhost",
                    help='Store the application health endpoint')
parser.add_argument('-percentage', action='store',
                    dest='percentage',
                    default="50",
                    help='Store the application health endpoint')


results = parser.parse_args()

# adopt log structure used by the chaostoolkit framework
logging.basicConfig(
    format="[%(asctime)s %(levelname)-2s] [%(module)s:%(lineno)s] %(message)s",
    level=logging.DEBUG,
    datefmt='%Y-%m-%d %H:%M:%S')

env_params = dict(
        LABEL_NAME=results.label,
        NAME_SPACE=results.namespace,
        APP_ENDPOINT=results.app,
        PERCENTAGE=int(results.percentage),
        FILE=results.file,
        EXP=results.exp
    )

# check (&set) env based on input and/or default values
for key in env_params:
    if key in os.environ.keys():
        logging.debug("Environment exists for key: %s", key)
    else:
        os.environ[key] = env_params[key]

filename = os.environ['FILE']
namespace = os.environ['NAME_SPACE']
experiment = os.environ['EXP']

# if the env CHAOSENGINE is defined, suffix it standard experiment name 
# to generate the fully-qualified chaos experiment/chaosresult name 

if 'CHAOSENGINE' in os.environ.keys():
    experiment_name = os.environ['CHAOSENGINE'] + '-' + experiment
else:
    experiment_name = experiment

# create chaosresult custom resource with phase=Running, verdict=Awaited
chaos_result_tracker(experiment_name, 'Running', 'Awaited', namespace)

# run chaos and store status into journal.json
chaos_command_list = ['chaos', '--verbose', 'run', '--journal-path', 'journal.json', filename]
run_shell_task(chaos_command_list)

# extract stage-wise success of experiment from the journal.json
with open('journal.json') as fp:
    data = json.load(fp)

pre_chaos_steady_state_check = data['steady_states']['before']['probes'][0]['status']
logging.info('status of pre-chaos steady_state_check is: %s', pre_chaos_steady_state_check)

run_status = data['run'][0]['status']
logging.info('status of chaos_injection_action is: %s', run_status)

post_chaos_steady_state_check = data['steady_states']['after']['probes'][0]['status']
logging.info('status of post-chaos steady_state_check is: %s', post_chaos_steady_state_check)

stage_level_verdicts = [pre_chaos_steady_state_check, run_status, post_chaos_steady_state_check]

# derive chaos experiment verdict as a logical AND of stage-wise results
if len(set(stage_level_verdicts)) == 1 and 'succeeded' in stage_level_verdicts:
    # patch chaosresult custom resource with phase=Completed, verdict=Pass
    chaos_result_tracker(experiment_name, 'Completed', 'Pass', namespace)
    logging.info('The chaos experiment verdict is: Pass')
else:
    # patch chaosresult custom resource with phase=Completed, verdict=Fail
    chaos_result_tracker(experiment_name, 'Completed', 'Fail', namespace)
    logging.info('The chaos experiment verdict is: Fail')
