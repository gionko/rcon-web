import configobj
import flask
import flask_kvsession
import flask_seasurf
import functools
import os
import simplekv.memory


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


@app.route('/')
@secured
def index():
    title = 'Server name'
    players_status = '2 humans, 8 bots'
    players = [
        {
            'name': 'Player 1',
            'value': 'STEAM_1:1:12345678'
        },
        {
            'name': 'Player 2',
            'value': 'STEAM_1:1:87654321'
        }
    ]
    maps_status = 'buhriz_coop'
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


@app.route('/ban')
@secured_ajax
def ban():
    return flask.jsonify(ok=True, error=None)


@app.route('/kick')
@secured_ajax
def kick():
    return flask.jsonify(ok=True, error=None)


@app.route('/map')
@secured_ajax
def map():
    return flask.jsonify(ok=True, error=None)
