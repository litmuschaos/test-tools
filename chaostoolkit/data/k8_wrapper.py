import subprocess
import sys
import os
import argparse
import logging
import json
from jinja2 import Environment, FileSystemLoader, select_autoescape
import yaml
from utils import *
from report import *

__author__ = 'Sumit_Nagal@intuit.com'

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
parser.add_argument('-report', action='store',
                    dest='report',
                    default="false",
                    help='Option to upload the result to report server')
parser.add_argument('-report_endpoint', action='store',
                    dest='report_endpoint',
                    default="none",
                    help='Endpoint where the report data will be uploaded')                                        


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
        REPORT=results.report,
        REPORT_ENDPOINT=results.report_endpoint,
        EXP=results.exp
    )

# check (&set) env based on input and/or default values
for key in env_params:
    if key in os.environ.keys():
        logging.debug("Environment exists for key: %s", key)
        env_params[key] = os.environ[key]
    else:
        os.environ[key] = env_params[key]

filename = os.environ['FILE']
namespace = os.environ['NAME_SPACE']
experiment = os.environ['EXP']
report = os.environ['REPORT']
report_endpoint = os.environ['REPORT_ENDPOINT']


# if the env CHAOSENGINE is defined, suffix it standard experiment name 
# to generate the fully-qualified chaos experiment/chaosresult name 

if 'CHAOSENGINE' in os.environ.keys():
    experiment_name = os.environ['CHAOSENGINE'] + '-' + experiment
else:
    experiment_name = experiment

# create chaosresult custom resource with phase=Running, verdict=Awaited
Utils().chaos_result_tracker(experiment_name, 'Running', 'Awaited', namespace)

# replacing percentage in json file as set_vars takes only string and int is needed for chaos
with open(filename, 'r') as file:
    json_data = json.load(file)
    if json_data['method'][0]['provider']['arguments']['mode'] == 'percentage':
        json_data['method'][0]['provider']['arguments']['qty'] = int(os.environ['PERCENTAGE'])
        with open(filename, 'w') as outfile:
            json.dump(json_data, outfile)

# run chaos and store status into journal.json
chaos_command_list = ['chaos', '--verbose', 'run', '--journal-path', 'journal.json', filename]
Utils().run_shell_task(chaos_command_list)

if report == 'true':
    json_data = {}
    logging.info('report end point is : %s', report_endpoint)
    json_data = Report().run(env_params, "journal.json")
    logger.info("Output data in main:---")
    logging.info(json_data)
    logger.info("----End of output data in main")


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
    Utils().chaos_result_tracker(experiment_name, 'Completed', 'Pass', namespace)
    logging.info('The chaos experiment verdict is: Pass')
else:
    # patch chaosresult custom resource with phase=Completed, verdict=Fail
    Utils().chaos_result_tracker(experiment_name, 'Completed', 'Fail', namespace)
    logging.info('The chaos experiment verdict is: Fail')

