# Installation

## Set-up Nginx web server

### Install Nginx web server

```
sudo apt-get install nginx
```

### Add rcon-web Nginx site configuration into `/etc/nginx/sites-available/rcon-web`

```
server {
	listen 80;
	server_name localhost;

	location / {
		try_files $uri @app;
	}
	location @app {
		include uwsgi_params;
		uwsgi_pass unix:/tmp/rcon-web.sock;
	}
}
```

### Enable rcon-web Nginx site

```
ln -s /etc/nginx/sites-available/rcon-web /etc/nginx/sites-enabled/rcon-web
```

### Restart Nginx web server

```
sudo service nginx restart
```

## Install required packages

```
sudo apt-get install build-essential python-dev python-pip
```

## Install Python dependencies

```
sudo pip install -r requirements.txt
```

## Auto-start rcon-web on boot

Add following line into `/etc/rc.local` before `exit 0`:

```
/usr/local/bin/uwsgi --ini /srv/rcon-web/uwsgi.ini --uid www-data --gid www-data --daemonize /var/log/uwsgi.log --touch-reload /tmp/rcon-web-restart
```

## Start rcon-web

```
sudo /usr/local/bin/uwsgi --ini /srv/rcon-web/uwsgi.ini --uid www-data --gid www-data --daemonize /var/log/uwsgi.log --touch-reload /tmp/rcon-web-restart
```

# Restart

It is not enough to simply restart the Nginx web server, since uWSGI application restart is done manually:

```
touch /tmp/rcon-web-restart
```
