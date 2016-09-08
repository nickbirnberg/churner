FROM ubuntu:xenial
MAINTAINER Nick Birnberg <nick@bitofnick.com>

RUN locale-gen en_US.UTF-8  
ENV LANG en_US.UTF-8  
ENV LANGUAGE en_US:en  
ENV LC_ALL en_US.UTF-8  

RUN apt-get -y update && \
    apt-get -y install python3.5 python3-pip
RUN apt-get clean

COPY . /action
WORKDIR /action

RUN pip3 install flask

ENV FLASK_APP action.py

EXPOSE 5000

ENTRYPOINT ["python3"]

CMD ["action.py"]