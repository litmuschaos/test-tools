import subprocess
import sys
import os
import argparse
import logging
import json
import threading

import requests
# from kubernetes.utils.utils import Utils
# from kubernetes.utils.report import Report

from jinja2 import Environment, FileSystemLoader, select_autoescape
import yaml

__author__ = 'Sumit_Nagal@intuit.com'

logger = logging.getLogger(__name__)


def run(self, serializer, journal, report_endpoint):
	kafka_thread = threading.Thread(target=self.report_post, args=(serializer, journal, report_endpoint))
	kafka_thread.start()


def get_value(self, element, *keys):
	_element = element
	for key in keys:
		try:
			_element = _element[key]
			if isinstance(_element, list):
				_element = _element[0]
		except Exception:
			return 'NA'
	return str(_element)


def json_parser(self, data, serializer, responsetext=None):
	output_data = {}
	# Get namespace, appEndpoint, label and application from serializer
	output_data['action'] = serializer['EXP']
	output_data['appEndpoint'] = serializer['APP_ENDPOINT']
	output_data['custom'] = {}
	output_data['custom']['namespace'] = serializer['NAME_SPACE']
	output_data['custom']['application'] = serializer['LABEL_NAME']

	output_data['scenarioName'] = self.get_value(data, 'experiment', 'method', 'provider', 'module')
	output_data['transaction'] = self.get_value(data, 'experiment', 'method', 'provider', 'func')
	output_data['timestamp'] = data.get('start') or 'NA'
	output_data['status'] = data.get('status') or 'NA'
	output_data['run_id'] = "journal" + serializer['EXP']

	if output_data['status'] == 'failed':
		output_data['run_status'] = 'failed'

		# if filename:
		#     if filename.find('.json'):
		#         output_data['run_id'] = filename[0: len(filename)-len('.json')]
		#     else:
		#         output_data['run_id'] = 'NA'
		# else:
		output_data['run_id'] = 'NA'

	output_data['custom']['steady_status_before_name'] = self.get_value(data, 'steady_states', 'before', 'probes',
																		'activity', 'name')
	output_data['custom']['steady_status_after_name'] = self.get_value(data, 'steady_states', 'after', 'probes',
																	   'activity', 'name')
	output_data['custom']['before_steady_state_met'] = self.get_value(data, 'steady_states', 'before',
																	  'steady_state_met')
	output_data['custom']['after_steady_state_met'] = self.get_value(data, 'steady_states', 'after',
																	 'steady_state_met')
	output_data['run_name'] = self.get_value(data, 'run', 'activity', 'name')
	output_data['run_status'] = self.get_value(data, 'run', 'status')
	output_data['custom']['rollback'] = data.get('rollbacks') or 'NA'

	logger.info("Output kubernetes:---")
	logger.info(output_data)
	logger.info("----End of output kubernetes")

	return output_data


# Try to find the journal file, try every 5 seconds, try 50 times
# After find the journal file, according to the timeout value, wait for that long then do the kafka post
# Format for Json upload
# {
# "run_id":"journal",
# "status":"completed",
# "scenarioName":"chaosk8s.pod.actions",
# "custom":{
#     "application":"<Label name>",
#     "after_steady_state_met":"True",
#     "rollback":"NA",
#     "before_steady_state_met":"True",
#     "steady_status_before_name":"application-must-respond-normally",
#     "steady_status_after_name":"application-must-respond-normally",
#     "namespace":"<namespace>"
# },
# "reportEndpoint":"<kafka endpoint>",
# "action":"k8-pod-delete",
# "run_name":"Terminate_pod",
# "timestamp":"2020-03-21T01:02:46.477100",
# "run_status":"succeeded",
# "transaction":"terminate_pods",
# "appEndpoint":"<app health endpoint>"
# }
def report_post(self, serializer, journal, report_endpoint):
	json_data = self.json_parser(journal, serializer)
	# EVENT_BUS_ENDPOINT = report_endpoint
	# namespace = json_data['custom']['namespace']
	# experiment_name = json_data['scenarioName']

	if report_endpoint != "none":
		try:
			logger.info("Start making http post request to Kafka %s", report_endpoint)
			headers = {'Content-type': 'application/json'}
			response = requests.post(report_endpoint, data=json.dumps(json_data, indent=4), headers=headers)
			logger.info(
				"Successfully posted kubernetes to Kafka end point " + report_endpoint + " with status code {0}".format(
					response.status_code))
		except Exception as e:
			logger.info("Posting kubernetes to Kafka failed: {0}".format(e))

	else:
		raise Exception("Unable to find report endpoint, value is -> " + report_endpoint)


"""
	run_shell_task() runs a shell command and prints the output as it executes.
	It takes a list of strings that comprises the command itself, as the sole arg.
	"""


def run_shell_task(self, cmd_arg_list):
	run_cmd = subprocess.Popen(cmd_arg_list, stdout=subprocess.PIPE, env=os.environ.copy())
	run_cmd.communicate()

	"""
	chaos_result_tracker() creates/patches the litmus chaosresult custom resource in the provided namespace.
	Typically invoked before and after chaos, and takes the .spec.phase, .spec.verdict & namespace as as args.
	"""


def chaos_result_tracker(self, exp_name, exp_phase, exp_verdict, ns):
	env_tmpl = Environment(loader=FileSystemLoader('./'), trim_blocks=True, lstrip_blocks=True,
						   autoescape=select_autoescape(['yaml']))
	template = env_tmpl.get_template('chaos-result.j2')
	updated_chaosresult_template = template.render(c_experiment=exp_name, phase=exp_phase, verdict=exp_verdict)
	with open('chaosresult.yaml', "w+") as f:
		f.write(updated_chaosresult_template)
	chaosresult_update_cmd_args_list = ['kubectl', 'apply', '-f', 'chaosresult.yaml', '-n', ns]
	self.run_shell_task(chaosresult_update_cmd_args_list)


####################################
# Start of Python Chaos Experiment #
####################################

parser = argparse.ArgumentParser()

parser.add_argument("-file", action='store',
					default="pod-app-kill-count-aks.json",
					dest="file",
					help="Chaos file to chose for execution"
					)
parser.add_argument("-exp", action='store',
					default="k8-pod-delete",
					dest="exp",
					help="Chaos experiment to chose for execution"
					)
parser.add_argument('-label', action='store',
					dest='label',
					default="app",
					help='Store a label value')
parser.add_argument("-namespace", action='store',
					default="default",
					dest="namespace",
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
