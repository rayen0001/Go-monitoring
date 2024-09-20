from flask import Flask, request, jsonify
from prometheus_client import start_http_server, Counter
import time
import random

app = Flask(__name__)

# Define a Counter metric for visitor count
visitor_counter = Counter('webapp1_visitors_total', 'Total number of visitors to webapp1')

@app.route('/')
def home():
    visitor_counter.inc()  # Increment the counter for each visit
    # Simulate fluctuating CPU and memory usage
    simulate_load()
    visitor_count = visitor_counter._value.get()  # Get the current value of the gauge
    return f"Hello from Webapp1! "

@app.route('/count', methods=['GET'])
def count():
    return jsonify({"visitor_count": visitor_counter._value.get()})

def simulate_load():
    # Simulate CPU load by busy-waiting
    duration = random.uniform(0.1, 0.5)  # Random duration between 0.1 and 0.5 seconds
    end_time = time.time() + duration
    while time.time() < end_time:
        pass

    # Simulate memory load by allocating random amount of memory
    memory_load = random.randint(1, 5) * 1024 * 1024  # Random memory between 1 and 5 MB
    _ = [0] * memory_load

if __name__ == '__main__':
    start_http_server(5001)  # Serve metrics on port 5001
    app.run(host='0.0.0.0', port=8081)  # Serve the web app on port 8081
