from flask import Flask, render_template
from services import Services

app = Flask(__name__)


@app.get("/")
def home():
    return render_template("index.html", services=Services.all)
