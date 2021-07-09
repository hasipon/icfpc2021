import json
import glob
import os
import subprocess
import pathlib
import shutil
from flask import Flask, request, render_template, abort

static_path = pathlib.Path(__file__).resolve().parent / 'static'
repo_path = pathlib.Path(__file__).resolve().parent.parent
app = Flask(__name__, static_folder=str(static_path), static_url_path='')
app.config['SEND_FILE_MAX_AGE_DEFAULT'] = 0


@app.after_request
def add_header(response):
    if 'Expires' in response.headers:
        del response.headers['Expires']
    response.headers['Cache-Control'] = 'no-store'
    return response


@app.route('/')
def index():
    problems_path = repo_path / "problems"
    problems_json = json.loads((static_path / "problems.json").read_text(encoding='utf-8'))
    dislikes = {
        x[0]: (
            int(x[1]) if x[1].isdigit() else '',
            int(x[2]) if x[2].isdigit() else '',
            int(x[1]) - int(x[2]) if x[1].isdigit() and x[2].isdigit() else ''
        ) for x in problems_json}

    problem_files = [os.path.relpath(x, problems_path) for x in glob.glob(str(problems_path / "*"))]
    problem_files.sort(key=lambda x: int(x))
    problems = [
        {
            "name": x,
            "dislike": dislikes[x][0],
            "mindislike": dislikes[x][1],
            "difflike": dislikes[x][2]
        }
        for x in problem_files
    ]
    return render_template('index.html', problems=problems)


@app.route('/git_status')
def git_status():
    output = ""
    try:
        output += subprocess.check_output(["git", "status"], stderr=subprocess.STDOUT).decode('utf-8').strip()
    except subprocess.CalledProcessError as e:
        output += "Error:" + str(e)
    return render_template('output.html', output=output)


@app.route('/fetch_problems')
def fetch_problems():
    output = ""
    try:
        output += subprocess.check_output(["node", "main.js"], cwd=(repo_path / 'portal')).decode("utf-8").strip()
        shutil.copyfile(repo_path / "portal/problems.json", static_path / "problems.json")
    except subprocess.CalledProcessError as e:
        output += "Error:" + str(e)
    return render_template('output.html', output=output)


@app.route('/git_pull')
def git_pull():
    output = ""
    try:
        output += subprocess.check_output(["git", "pull"], stderr=subprocess.STDOUT).decode(
            'utf-8').strip()
    except subprocess.CalledProcessError as e:
        output += "Error:" + str(e)
    return render_template('output.html', output=output)


if __name__ == "__main__":
    app.run(host='0.0.0.0', port=5000, threaded=True, debug=True)
