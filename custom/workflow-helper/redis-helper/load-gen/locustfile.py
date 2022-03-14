#!/usr/bin/python3
## pylint: disable = invalid-name, too-few-public-methods

from random import randint
import json
import time
import redis
from locust import User, events, TaskSet, task, constant
import gevent.monkey, os
gevent.monkey.patch_all()
from pyredis import Client

redis_host = os.getenv("REDIS_HOST", "redis.redis.svc")
redis_port = os.getenv("REDIS_PORT", "6379")
redis_pw = os.getenv("REDIS_PW", "")

# Locust helps in defining website user behavior with code and swarms your system with millions of simultaneous users.
# Logs of the metrics contain requests per second, total requests, failed requests, average, min, max, and failed per second.
class RedisClient(object):
    def __init__(self, host=redis_host, port=redis_port, password=redis_pw):
        self.rc = redis.StrictRedis(host=host, port=port, password=password)

    def query(self, key, command='GET'):
        """Function to Test GET operation on Redis"""
        result = None
        start_time = time.time()
        try:
            result = self.rc.get(key)
            if not result:
                result = ''
        except Exception as e:
            total_time = int((time.time() - start_time) * 1000)
            events.request_failure.fire(
                request_type=command, name=key, response_time=total_time, exception=e)
        else:
            total_time = int((time.time() - start_time) * 1000)
            length = len(result)
            events.request_success.fire(
                request_type=command, name=key, response_time=total_time, response_length=length)
        return result

    def write(self, key, value, command='SET'):
        """Function to Test SET operation on Redis"""
        result = None
        start_time = time.time()
        try:
            result = self.rc.set(key, value)
            if not result:
                result = ''
        except Exception as e:
            total_time = int((time.time() - start_time) * 1000)
            events.request_failure.fire(
                request_type=command, name=key, response_time=total_time, exception=e)
        else:
            total_time = int((time.time() - start_time) * 1000)
            length = 1
            events.request_success.fire(
                request_type=command, name=key, response_time=total_time, response_length=length)
        return result


class RedisLocust(User):
    wait_time = constant(0.1)
    key_range = 500

    def __init__(self, *args, **kwargs):
        super(RedisLocust, self).__init__(*args, **kwargs)
        self.client = RedisClient()
        self.key = 'key1'
        self.value = 'value1'

    @task(2)
    def get_time(self):
        for i in range(self.key_range):
            self.key = 'key'+str(i)
            self.client.query(self.key)

    @task(1)
    def write(self):
        for i in range(self.key_range):
            self.key = 'key'+str(i)
            self.value = 'value'+str(i)
            self.client.write(self.key, self.value)

    @task(1)
    def get_key(self):
        var = str(randint(1, self.key_range-1))
        self.key = 'key'+var
        self.value = 'value'+var