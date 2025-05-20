import json
import hashlib
import random
import string
import time
from datetime import datetime
from flask import Flask, request, jsonify

app = Flask(__name__)

# Constants
DATA_FILE = 'data.json'
PASSWORD = 'YOUR_ADMIN_PASSWORD'  # Change this in production

# Helper functions
def load_data():
    try:
        with open(DATA_FILE, 'r') as f:
            return json.load(f)
    except FileNotFoundError:
        return {}

def save_data(data):
    with open(DATA_FILE, 'w') as f:
        json.dump(data, f, indent=4)

def generate_license(name):
    md5_hash = hashlib.md5(name.encode()).hexdigest()
    sha256_hash = hashlib.sha256(name.encode()).hexdigest()
    combined = ''.join(random.choices(md5_hash + sha256_hash, k=32))
    return ''.join(random.choices(string.ascii_uppercase + string.digits, k=32))

def is_valid_license_format(lic):
    return len(lic) == 32 and all(c in string.ascii_uppercase + string.digits for c in lic)

def is_valid_ip(ip):
    parts = ip.split('.')
    if len(parts) != 4:
        return False
    for part in parts:
        if not part.isdigit() or not 0 <= int(part) <= 255:
            return False
    return True

def format_expire_date(timestamp):
    return datetime.fromtimestamp(timestamp).strftime('%Y-%m-%d')

# API Routes
@app.route('/check-license', methods=['POST'])
def check_license():
    data = request.json
    lic = data.get('lic')
    ip = data.get('ip')

    if not lic or not ip:
        return jsonify({'error': 'Missing parameters'}), 400

    if not is_valid_license_format(lic):
        return jsonify({'error': 'Invalid license format'}), 400

    if not is_valid_ip(ip):
        return jsonify({'error': 'Invalid IP format'}), 400

    licenses = load_data()
    if lic not in licenses:
        return jsonify({'error': 'License not found'}), 404

    license_data = licenses[lic]
    if ip not in license_data['ip_list']:
        return jsonify({'error': 'IP not allowed'}), 403

    if time.time() > license_data['expire']:
        del licenses[lic]
        save_data(licenses)
        return jsonify({'error': 'License expired and has been removed'}), 403

    license_data['expire'] = format_expire_date(license_data['expire'])

    return jsonify(license_data), 200

@app.route('/create-lic', methods=['GET'])
def create_license():
    password = request.args.get('password')
    name = request.args.get('name')
    expire_days = request.args.get('expire', type=int)

    if password != PASSWORD:
        return jsonify({'error': 'Unauthorized'}), 401

    if not name or not expire_days:
        return jsonify({'error': 'Missing parameters'}), 400

    licenses = load_data()
    
    # Check for duplicate name
    if any(license_data['name'] == name for license_data in licenses.values()):
        return jsonify({'error': 'Name already exists'}), 409

    lic = generate_license(name)
    expire_timestamp = int(time.time() + expire_days * 86400)  # Convert days to seconds and ensure integer

    licenses[lic] = {
        'name': name,
        'ip_list': [],
        'expire': expire_timestamp
    }
    save_data(licenses)

    return jsonify({
        'license': lic,
        'name': name,
        'ip_list': [],
        'expire': format_expire_date(expire_timestamp)
    }), 201

@app.route('/add-ip', methods=['GET'])
def add_ip():
    password = request.args.get('password')
    lic = request.args.get('lic')
    ip = request.args.get('ip')

    if password != PASSWORD:
        return jsonify({'error': 'Unauthorized'}), 401

    if not lic or not ip:
        return jsonify({'error': 'Missing parameters'}), 400

    if not is_valid_ip(ip):
        return jsonify({'error': 'Invalid IP format'}), 400

    licenses = load_data()
    if lic not in licenses:
        return jsonify({'error': 'License not found'}), 404

    license_data = licenses[lic]
    if len(license_data['ip_list']) >= 2:
        return jsonify({'error': 'IP list is full'}), 403

    if ip in license_data['ip_list']:
        return jsonify({'error': 'IP already exists in the list'}), 409

    license_data['ip_list'].append(ip)
    save_data(licenses)

    return jsonify({'message': 'IP added successfully'}), 200

@app.route('/clear-ips', methods=['GET'])
def clear_ips():
    password = request.args.get('password')
    lic = request.args.get('lic')

    if password != PASSWORD:
        return jsonify({'error': 'Unauthorized'}), 401

    if not lic:
        return jsonify({'error': 'Missing parameters'}), 400

    licenses = load_data()
    if lic not in licenses:
        return jsonify({'error': 'License not found'}), 404

    license_data = licenses[lic]
    ip_list_length = len(license_data['ip_list'])

    if ip_list_length >= 2:
        license_data['ip_list'] = []
        save_data(licenses)
        return jsonify({'message': 'IP list cleared successfully'}), 200
    else:
        remaining_ips = 2 - ip_list_length
        return jsonify({'message': f'There are still {remaining_ips} IPs that can be added. No need to clear IP list.'}), 200

@app.route('/delete-lic', methods=['GET'])
def delete_lic():
    password = request.args.get('password')
    lic = request.args.get('lic')

    if password != PASSWORD:
        return jsonify({'error': 'Unauthorized'}), 401

    if not lic:
        return jsonify({'error': 'Missing parameters'}), 400

    licenses = load_data()
    if lic not in licenses:
        return jsonify({'error': 'License not found'}), 404

    del licenses[lic]
    save_data(licenses)
    return jsonify({'message': 'License deleted successfully'}), 200

@app.route('/all', methods=['GET'])
def get_all_licenses():
    password = request.args.get('password')

    if password != PASSWORD:
        return jsonify({'error': 'Unauthorized'}), 401

    licenses = load_data()
    for lic, license_data in licenses.items():
        license_data['expire'] = format_expire_date(license_data['expire'])
    
    return jsonify(licenses), 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8000,debug=True)
