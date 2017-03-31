import configobj
import flask
import flask_kvsession
import flask_seasurf
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


@app.route('/')
def index():
    if 'login' not in flask.session:
        return flask.redirect(flask.url_for('login'))

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
