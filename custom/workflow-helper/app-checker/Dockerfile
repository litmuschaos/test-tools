FROM python:3

LABEL maintainer="LitmusChaos"

ADD main.py /
RUN pip install --upgrade pip
RUN pip install requests
RUN pip install uuid

CMD [ "python", "./main.py" ]
