import flask
import flask_kvsession
import simplekv.memory


# Create Flask application
app = flask.Flask(__name__)

# Create session store
store = simplekv.memory.DictStore()
flask_kvsession.KVSessionExtension(store, app)


@app.route("/")
def index():
    return "<html><body><h1>RCON-Web</h1></body></html>"
