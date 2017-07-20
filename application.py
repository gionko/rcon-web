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

    # Get status

    with valve.rcon.RCON(server, password) as rcon:
        data = rcon('status')

    title = regex_single(data, 'hostname *: *(.*)')

    players_status = regex_single(data, 'players *: *(\d* humans, \d* bots \(\d*/\d* max\)).*')
    players = []
    for name, steamid, connected, ping in regex_all(data, '# .*\d* .*\d* \"(.*)\" .*(STEAM.*?) (\d*:?\d*:\d*) (\d*) .*'):
        player = {
            'name':      name,
            'steamid':   steamid,
            'connected': connected,
            'ping':      ping
        }
        players.append(player)

    maps_status = regex_single(data, 'map *: *(.*)')
    maps = [{'name': name, 'value': config['maps'][name]['value']} for name in config['maps']]

    # Get configuration

    with valve.rcon.RCON(server, password) as rcon:
        data = rcon('bot_damage')

    config_damage = regex_single(data, '.*= *\"([\d.]*)\".*')

    with valve.rcon.RCON(server, password) as rcon:
        data = rcon('ins_bot_difficulty')

    config_difficulty = regex_single(data, '.*= *\"([\d.]*)\".*')

    config_status = 'damage {:.3}, difficulty {}'.format(config_damage, config_difficulty)

    # Render template

    return flask.render_template('index.html',
                                 title=title,
                                 players_status=players_status,
                                 players=players,
                                 maps_status=maps_status,
                                 maps=maps,
                                 config_status=config_status)


@app.route('/login', methods=['GET', 'POST'])
def login():
    if flask.request.method == 'POST':
        if flask.request.form['password'] == password:
            flask.session['login'] = True
            return flask.redirect(flask.url_for('index'))

    flask.session.pop('login', None)
    return flask.render_template('login.html')


@app.route('/ban', methods=['POST'])
@secured_ajax
def ajax_ban():
    data = flask.request.get_json()
    with valve.rcon.RCON(server, password) as rcon:
        rcon('banid {} {}'.format(data['period'], data['player']))
        rcon('kickid {} {}'.format(data['player'], data['message']))
    return flask.jsonify(ok=True, error=None)


@app.route('/config', methods=['POST'])
@secured_ajax
def ajax_config():
    data = flask.request.get_json()
    with valve.rcon.RCON(server, password) as rcon:
        rcon('{} {}'.format(data['name'], data['value']))
    return flask.jsonify(ok=True, error=None)


@app.route('/kick', methods=['POST'])
@secured_ajax
def ajax_kick():
    data = flask.request.get_json()
    with valve.rcon.RCON(server, password) as rcon:
        rcon('kickid {} {}'.format(data['player'], data['message']))
    return flask.jsonify(ok=True, error=None)


@app.route('/map', methods=['POST'])
@secured_ajax
def ajax_map():
    data = flask.request.get_json()
    with valve.rcon.RCON(server, password) as rcon:
        rcon('changelevel {}'.format(data['map']))
    return flask.jsonify(ok=True, error=None)
