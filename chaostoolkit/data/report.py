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
from utils import *
logger = logging.getLogger(__name__)


class Report(object):

    ####################################
    #      Function definitions        #
    ####################################

    def run(self, serializer, journal_file):
        kafka_thread = threading.Thread(target=self.report_post, args=(serializer, journal_file))
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

    def file_parser(self, data, serializer, filename, responsetext=None):

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
        output_data['reportEndpoint'] = serializer['REPORT_ENDPOINT']

        if output_data['status'] == 'failed':
            output_data['run_status'] = 'failed'

        if filename:
            if filename.find('.json'):
                output_data['run_id'] = filename[0: len(filename)-len('.json')]
            else:
                output_data['run_id'] = 'NA'
        else:
            output_data['run_id'] = 'NA'

        output_data['custom']['steady_status_before_name'] = self.get_value(data, 'steady_states', 'before', 'probes', 'activity', 'name')
        output_data['custom']['steady_status_after_name'] = self.get_value(data, 'steady_states', 'after', 'probes', 'activity', 'name')
        output_data['custom']['steady_status_after'] = self.get_value(data, 'steady_states', 'after', 'probes', 'status')
        output_data['custom']['steady_status_before'] = self.get_value(data, 'steady_states', 'after', 'probes', 'status')
        output_data['custom']['before_steady_state_met'] = self.get_value(data, 'steady_states', 'before', 'steady_state_met')
        output_data['custom']['after_steady_state_met'] = self.get_value(data, 'steady_states', 'after', 'steady_state_met')
        output_data['run_name'] = self.get_value(data, 'run', 'activity', 'name')
        output_data['run_status'] = self.get_value(data, 'run', 'status')
        output_data['custom']['rollback'] = data.get('rollbacks') or 'NA'

        logger.info("Output data:---")
        logger.info(output_data)
        logger.info("----End of output data")

        return output_data

    
    # Try to find the journal file, try every 5 seconds, try 50 times
    # After find the journal file, accrrding to the timeout value, wait for that long then do the kafka post
    def report_post(self, serializer, file_name, attempts=0, timeout=50, sleep_int=5):
        if attempts < timeout:
            try:
                with open(file_name, encoding='utf-8', errors='replace') as f:
                    #time.sleep(timeout)
                    logger.info("Trying to load json data")
                    file = f.read()
                    data = json.loads(file)
                    logger.info("End of json loading")
                    logger.info(data)
                    json_data = self.file_parser(data, serializer, file_name)
                    EVENT_BUS_ENDPOINT = json_data['reportEndpoint']
                    namespace = json_data['custom']['namespace']
                    experiment_name = json_data['scenarioName'] 
                    ## Format for Json upload
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
                    if EVENT_BUS_ENDPOINT != "none":
                        try:                                                    
                            logger.info("Start making http post request to Kafka %s", json_data['reportEndpoint'] )                        
                            headers = {'Content-type': 'application/json'}
                            response = requests.post(EVENT_BUS_ENDPOINT, data=json.dumps(json_data, indent=4), headers=headers)
                            logger.info("Successfully posted data to Kafka end point "+ EVENT_BUS_ENDPOINT + " with status code {0}".format(response.status_code))
                        except Exception as e:
                            logger.info("Posting data to Kafka failed: {0}".format(e))
                        
            except FileNotFoundError:
                logger.info("Unable to find matching json file with file name " + file_name + ". Attempted to try " + str(attempts) + " times")
                time.sleep(sleep_int)
                self.report_post(serializer, file_name, attempts + 1)
            except Exception as e:
                logger.info("JSON file load failed: {0}".format(e))
        else:
            raise Exception("Unable to find matching json file " + file_name + " in 3 minutes")


