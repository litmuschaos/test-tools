import json
import logging
import threading
import time
import requests

__author__ = 'Sumit_Nagal@intuit.com'

logger = logging.getLogger(__name__)


class Report(object):

    ####################################
    #      Function definitions        #
    ####################################

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
            logger.info("Not making http post request to Kafka %s", report_endpoint)
