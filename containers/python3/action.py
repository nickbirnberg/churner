import flask

app = flask.Flask(__name__)


@app.route('/run', methods=['POST'])
def run():
    action = flask.request.get_json(force=True, silent=True)

    if not action:
        flask.abort(403)

    # compile code given
    try:
        code = compile(action['Code'], '<string>', 'exec')
    except Exception as err:
        response = flask.jsonify({"Error": str(err)})
        response.status_code = 502
        return response

    # run compiled code
    func_resp = None
    global_vars = {'param': action['Param']}
    try:
        exec(code, global_vars)
        exec('func_resp = action_func(param)', global_vars)
        func_resp = global_vars['func_resp']
    except Exception as err:
        response = flask.jsonify({"Error": str(err)})
        response.status_code = 502
        return response

    response = flask.jsonify({"Response": func_resp, "Error": None})
    response.status_code = 200
    response.mimetype = "application/json"
    return response

if __name__ == '__main__':
    app.run(host='0.0.0.0')
