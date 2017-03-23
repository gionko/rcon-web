import flask
import flask_kvsession
import simplekv.memory


# Create Flask application
app = flask.Flask(__name__)

# Create session store
store = simplekv.memory.DictStore()
flask_kvsession.KVSessionExtension(store, app)


@app.route('/')
def index():
    return flask.render_template('index.html')
