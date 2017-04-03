import configobj
import flask
import flask_kvsession
import flask_seasurf
import functools
import os
import re
import simplekv.memory
import valve.rcon


# Create Flask application
app = flask.Flask(__name__)

# Configure application
app.config.update(
    DEBUG=True,
    SECRET_KEY='die, bitch!'
)

# Create session store
store = simplekv.memory.DictStore()
flask_kvsession.KVSessionExtension(store, app)

# Initialize CSRF prevention
flask_seasurf.SeaSurf(app)

# Load configuration
__path__ = os.path.dirname(os.path.realpath(__file__))
config = configobj.ConfigObj(__path__ + '/rcon-web.conf')

# Set global variables
server = (config['server']['address'], int(config['server']['port']))
password = config['server']['password']

# Set default RCON encoding to UTF-8
valve.rcon.RCONMessage.ENCODING = 'utf-8'


def secured(f):
    @functools.wraps(f)
    def decorator(*args, **kwargs):
        if 'login' not in flask.session:
            return flask.redirect(flask.url_for('login'))
        return f(*args, **kwargs)
    return decorator


def secured_ajax(f):
    @functools.wraps(f)
    def decorator(*args, **kwargs):
        if 'login' not in flask.session:
            return flask.jsonify(ok=False, error="Not authorized")
        return f(*args, **kwargs)
    return decorator


def regex_all(string, regex):
    return re.findall(regex, string, re.IGNORECASE)


def regex_single(string, regex):
    return regex_all(string, regex)[0]


@app.route('/')
@secured
def index():
    with valve.rcon.RCON(server, password) as rcon:
        data = rcon('status')

    title = regex_single(data, 'hostname *: *(.*)')

    players_status = regex_single(data, 'players *: *(\d* humans, \d* bots \(\d*/\d* max\)).*')
    players = []
    for name, value in regex_all(data, '# .*\d* .*\d* \"(.*)\" .*(STEAM.*?) .*'):
        player = {
            'name': name,
            'value': value
        }
        players.append(player)

    maps_status = regex_single(data, 'map *: *(.*)')
    maps = [{'name': name, 'value': config['maps'][name]['value']} for name in config['maps']]

    return flask.render_template('index.html',
                                 title=title,
                                 players_status=players_status,
                                 players=players,
                                 maps_status=maps_status,
                                 maps=maps)


@app.route('/login', methods=['GET', 'POST'])
def login():
    if flask.request.method == 'POST':
        if flask.request.form['password'] == config['server']['password']:
            flask.session['login'] = True
            return flask.redirect(flask.url_for('index'))

    flask.session.pop('login', None)
    return flask.render_template('login.html')


@app.route('/ban', methods=['POST'])
@secured_ajax
def ban():
    data = flask.request.get_json()
    with valve.rcon.RCON(server, password) as rcon:
        rcon('banid {} {}'.format(data['period'], data['player']))
        rcon('kickid {} {}'.format(data['player'], data['message']))
    return flask.jsonify(ok=True, error=None)


@app.route('/kick', methods=['POST'])
@secured_ajax
def kick():
    data = flask.request.get_json()
    with valve.rcon.RCON(server, password) as rcon:
        rcon('kickid {} {}'.format(data['player'], data['message']))
    return flask.jsonify(ok=True, error=None)


@app.route('/map', methods=['POST'])
@secured_ajax
def map():
    data = flask.request.get_json()
    with valve.rcon.RCON(server, password) as rcon:
        rcon('changelevel {}'.format(data['map']))
    return flask.jsonify(ok=True, error=None)
