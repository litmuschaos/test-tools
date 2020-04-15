"""Chaos toolkit runner and report class"""

import datetime
import json
import logging
import os
import site
import sys

import click
from chaoslib.control import load_global_controls
from chaoslib.exceptions import InvalidSource
from chaoslib.experiment import run_experiment
from chaoslib.loader import load_experiment
from  chaostest.utils.report import Report
from chaostoolkit import encoder

logger = logging.getLogger(__name__)


class ChaosUtils(object):

    def run_chaos_engine(self, file, env_params: dict, report: str, report_endpoint: str) -> bool:
        settings = ({}, os.environ.get("settings_path"))[os.environ.get("settings_path") is not None]
        has_deviated = False
        has_failed = False
        load_global_controls(settings)
        jornal_file_suffix = file
        try:
            try:
                with open(file, "r"):
                    logger.info("File exists in local")
            except FileNotFoundError:
                logger.info("File is not available in the current directory, looking inside site packages")
                location = site.getsitepackages().__getitem__(0)
                file_found = False
                for root, dirs, files in os.walk(location):
                    if file in files:
                        file_found = True
                        file = os.path.join(root, file)
                        break
                if not file_found:
                    logger.error("File " + file + " not found in site packages too, quitting")
                    raise FileNotFoundError("Chaos file is not found")
            experiment = load_experiment(
                click.format_filename(file), settings)
        except InvalidSource as x:
            logger.error(str(x))
            logger.debug(x)
            sys.exit(1)
        journal = run_experiment(experiment, settings=settings)
        has_deviated = journal.get("deviated", False)
        has_failed = journal["status"] != "completed"
        json_file_name = "journal" + "-" + jornal_file_suffix
        with open(json_file_name, "w") as r:
            json.dump(
                journal, r, indent=2, ensure_ascii=False, default=encoder)
        r.close()
        if report == 'true':
            self.create_report(env_params, journal, report_endpoint)
        if has_failed or has_deviated:
            logger.error("Test Failed")
            return has_failed or has_deviated
        else:
            logger.info("Test Passed")
            return True

    @staticmethod
    def create_report(env_params: dict, journal_file_name: str, report_endpoint):
        logging.info('report end point is : %s', report_endpoint)
        json_data = Report().run(env_params, journal_file_name, report_endpoint)
        logger.info("Output kubernetes in main:---")
        logging.info(json_data)
        logger.info("----End of output kubernetes in main")
