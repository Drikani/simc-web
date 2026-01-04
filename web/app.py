from flask import Flask, request, render_template, jsonify
import redis, uuid, os, subprocess

app = Flask(__name__)
r = redis.Redis(host="redis", decode_responses=True)
DATA = "/data"

@app.route("/", methods=["GET", "POST"])
def index():
    if request.method == "POST":
        job_id = str(uuid.uuid4())
        r.rpush("jobs", f"{job_id}::{request.form['simc']}")

        return render_template("index.html", job_id=job_id)

    return render_template("index.html")

@app.route("/status/<job_id>")
def status(job_id):
    if r.exists(f"done:{job_id}"):
        return jsonify(done=True)
    return jsonify(done=False)

@app.route("/result/<job_id>")
def result(job_id):
    xml = f"{DATA}/{job_id}.xml"
    html = f"{DATA}/{job_id}.html"

    subprocess.run([
        "xsltproc",
        "-o", html,
        "templates/result.xslt",
        xml
    ])

    return open(html).read()

app.run(host="0.0.0.0", port=5000)