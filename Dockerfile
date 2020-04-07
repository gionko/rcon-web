FROM ubuntu:18.04

RUN apt-get update && \
	apt-get --yes --force-yes install build-essential python-dev python-pip nano nginx git && \
	apt-get clean

WORKDIR /usr/src/app
RUN git clone https://github.com/gionko/rcon-web.git
WORKDIR /usr/src/app/rcon-web
RUN cp /usr/src/app/rcon-web/rcon-web /etc/nginx/sites-available/rcon-web
RUN ln -s /etc/nginx/sites-available/rcon-web /etc/nginx/sites-enabled/rcon-web
RUN service nginx restart
RUN pip install -r requirements.txt
RUN cp -r /usr/src/app/rcon-web /srv


EXPOSE 9006
ENTRYPOINT /bin/bash 
