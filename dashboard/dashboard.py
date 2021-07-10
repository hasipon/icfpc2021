import json
import glob
import math
import os
import subprocess
import pathlib
import shutil
from typing import *
from flask import Flask, request, render_template, abort

static_path = pathlib.Path(__file__).resolve().parent / 'static'
repo_path = pathlib.Path(__file__).resolve().parent.parent
problems_path = repo_path / "problems"
app = Flask(__name__, static_folder=str(static_path), static_url_path='')
app.config['SEND_FILE_MAX_AGE_DEFAULT'] = 0

# global cache
problem_details = {}


@app.after_request
def add_header(response):
    if 'Expires' in response.headers:
        del response.headers['Expires']
    response.headers['Cache-Control'] = 'no-store'
    return response


def gen_problem_svg(name: str, problem_detail: Dict):
    return render_template(
        'thumbnail.jinja2',
        name=name,
        maxx=max(max(h[0] for h in problem_detail["hole"]), max(f[0] for f in problem_detail["figure"]["vertices"])),
        maxy=max(max(h[1] for h in problem_detail["hole"]), max(f[1] for f in problem_detail["figure"]["vertices"])),
        hole=problem_detail["hole"],
        figure=problem_detail["figure"])


def load_problem_details(problem_files):
    details = {}
    for prob in problem_files:
        with open(problems_path / prob) as fp:
            details.update({prob: json.load(fp)})
        details[prob]["bonus_from"] = []
        details[prob]["bonus_to"] = []
        details[prob]["base_score"] = 1000 * math.log2(
            len(details[prob]["hole"]) *
            len(details[prob]["figure"]["vertices"]) *
            len(details[prob]["figure"]["edges"]) / 6.0)

        svg_path = static_path / (prob + ".svg")
        if not svg_path.exists():
            svg_path.write_text(gen_problem_svg(prob, details[prob]), encoding="utf-8")

    for prob in problem_files:
        if "bonuses" in details[prob]:
            for bonus in details[prob]["bonuses"]:
                details[prob]["bonus_to"].append((str(bonus["problem"]), bonus))

    for prob in problem_files:
        for to in details[prob]["bonus_to"]:
            to_id, to_bonus = to[0], to[1]
            details[str(to_id)]["bonus_from"].append((prob, to_bonus))

    return details


def filter_problems(problems):
    dislike_min = request.args.get("dislike-min")
    dislike_max = request.args.get("dislike-max")
    top_dislike_min = request.args.get("top-dislike-min")
    top_dislike_max = request.args.get("top-dislike-max")
    dislike_ratio_min = request.args.get("dislike-ratio-min")
    dislike_ratio_max = request.args.get("dislike-ratio-max")

    def fix_dislike(d):
        if type(d) == int:
            return d
        if not d:
            return 9999999999999999999
        return d

    def f(p):
        if dislike_max and int(dislike_max) < fix_dislike(p["dislike"]):
            return False
        if dislike_min and int(dislike_min) > fix_dislike(p["dislike"]):
            return False
        if top_dislike_max and int(top_dislike_max) < fix_dislike(p["dislike_min"]):
            return False
        if top_dislike_min and int(top_dislike_min) > fix_dislike(p["dislike_min"]):
            return False
        if dislike_ratio_max and int(dislike_ratio_max) < p["dislike_ratio"] * 100:
            return False
        if dislike_ratio_min and int(dislike_ratio_min) > p["dislike_ratio"] * 100:
            return False
        return True

    return list(filter(f, problems))


@app.route('/')
def index():
    problem_files = [os.path.relpath(x, problems_path) for x in glob.glob(str(problems_path / "*"))]
    problem_files.sort(key=lambda x: int(x))

    global problem_details
    if len(problem_details) != len(problem_files):
        problem_details = load_problem_details(problem_files)

    problems_json = json.loads((static_path / "problems.json").read_text(encoding='utf-8'))
    dislikes = {x: (None, None, 0) for x in problem_files}
    dislikes.update({
        x[0]: (
            int(x[1]) if x[1].isdigit() else None,  # 自分のdislike
            int(x[2]) if x[2].isdigit() else None,  # TOPのdislike
            (int(x[2]) + 1) / (int(x[1]) + 1) if x[1].isdigit() and x[2].isdigit() else 0,
        ) for x in problems_json})

    problems = [
        {
            "name": x,
            "hole": len(problem_details[x]["hole"]),
            "eps": problem_details[x]["epsilon"],
            "edges": len(problem_details[x]["figure"]["edges"]),
            "vertices": len(problem_details[x]["figure"]["vertices"]),
            "dislike": dislikes[x][0],
            "dislike_min": dislikes[x][1],
            "dislike_ratio": dislikes[x][2],
            "topscore": math.ceil(problem_details[x]["base_score"]),
            "score": math.ceil(problem_details[x]["base_score"] * math.sqrt(dislikes[x][2])),
            "bonus_from": problem_details[x]["bonus_from"],
            "bonus_to": problem_details[x]["bonus_to"],
        }
        for x in problem_files
    ]

    problems = filter_problems(problems)

    return render_template('index.html',
                           is_search=request.args.get("search"),
                           problems=problems)


@app.route('/filter')
def get_filter():
    return render_template('filter.jinja2')


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
