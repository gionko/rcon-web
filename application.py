from flask import Flask
app = Flask(__name__)

@app.route("/")
def main():
	return "<html><body><h1>RCON-Web</h1></body></html>"
