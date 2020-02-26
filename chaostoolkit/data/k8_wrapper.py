import subprocess
import sys
import os
import argparse
import logging


__author__ = 'Sumit_Nagal@intuit.com'

parser = argparse.ArgumentParser()

parser.add_argument("-file", action='store',
                    default="experiment.json",
                    dest = "file",
                    help="Chaos file to chose for execution"
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


env_params = dict(
        LABEL_NAME=results.label,
        NAME_SPACE=results.namespace,
        APP_ENDPOINT=results.app,
        PERCENTAGE=results.percentage,
        FILE=results.file
    )

for key in env_params:
    if key in os.environ.keys():
        print("Environment exists")
    else:
        os.environ[key] = env_params[key]

filename = os.environ['FILE']

execcommandlist = []
execcommandlist.append("chaos")
execcommandlist.append("--verbose")
execcommandlist.append("run")
execcommandlist.append(filename)
runchaoscmd = subprocess.Popen(execcommandlist, stdout=subprocess.PIPE, encoding='utf-8', env=os.environ.copy())

for line in runchaoscmd.stdout:
    logging.info(line)
