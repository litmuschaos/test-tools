import jenkins
import time
import os
import sys

minutes=os.environ['MINUTES']
sv=os.environ['SERVICE']
ns=os.environ['NAMESPACE']
user=os.environ['USER']
password=os.environ['PASSWORD']

timeout = time.time() + 60*float(minutes)

def job_simulator(server):        
    while True:
        if time.time() > timeout:
            break   
        server.create_job('empty', jenkins.EMPTY_CONFIG_XML)
        time.sleep(4)
        print(server.jobs_count())
        jobs = server.get_jobs()
        print(jobs)
        sys.stdout.flush()
        my_job = server.get_job_config('empty')
        print(my_job) # prints XML configuration
        sys.stdout.flush()
        server.copy_job('empty', 'empty_copy')
        server.enable_job('empty_copy')
        time.sleep(4)
        jobs = server.get_jobs()
        print(jobs)
        sys.stdout.flush()
        print(server.jobs_count())
        sys.stdout.flush()
        server.delete_job('empty')
        server.delete_job('empty_copy')

if __name__== "__main__":
     # connect to the Jenkins server
    try:
        url = "http://"+sv+"."+ns+"."+"svc.cluster.local"
        server = jenkins.Jenkins(url, username=user, password=password)
        user = server.get_whoami()
        version = server.get_version()
        print('Hello %s from Jenkins %s' % (user['fullName'], version))
        job_simulator(server)
    except Exception as e:
        print('Error:', e)
