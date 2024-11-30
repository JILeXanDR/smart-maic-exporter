import logging
import os

import requests
from flask import Flask, Response
from flask_caching import Cache
from prometheus_client import CollectorRegistry, Gauge, generate_latest

# Configuration from environment variables
DATA_SOURCE_URL = os.getenv("DATA_SOURCE_URL", "http://192.168.10.55/?page=getwdata")
EXPORTER_PORT = int(os.getenv("EXPORTER_PORT", 8000))
CACHE_TIMEOUT = 5  # Cache duration in seconds

# Flask App
app = Flask(__name__)

# Configure Cache
app.config["CACHE_TYPE"] = "SimpleCache"  # Use in-memory caching
app.config["CACHE_DEFAULT_TIMEOUT"] = CACHE_TIMEOUT
cache = Cache(app)

# Initialize Prometheus Registry
registry = CollectorRegistry()

# Define Prometheus Metrics with a common prefix
metric_prefix = "smart_maic_"
metrics = {
    "voltage": Gauge(f"{metric_prefix}voltage", "Voltage (V)", ["phase"], registry=registry),
    "current": Gauge(f"{metric_prefix}current", "Current (A)", ["phase"], registry=registry),
    "power": Gauge(f"{metric_prefix}power", "Active Power (W)", ["phase"], registry=registry),
    "energy": Gauge(f"{metric_prefix}energy", "Energy (Wh)", ["phase"], registry=registry),
    "power_factor": Gauge(f"{metric_prefix}power_factor", "Power Factor", ["phase"], registry=registry),
    "frequency": Gauge(f"{metric_prefix}frequency", "Frequency (Hz)", ["phase"], registry=registry),
    "total_current": Gauge(f"{metric_prefix}total_current", "Total Current (A)", registry=registry),
    "total_power": Gauge(f"{metric_prefix}total_power", "Total Active Power (W)", registry=registry),
    "total_energy": Gauge(f"{metric_prefix}total_energy", "Total Energy (Wh)", registry=registry),
    "temperature": Gauge(f"{metric_prefix}temperature", "Device Temperature (°C)", registry=registry),
    "request_status": Gauge(f"{metric_prefix}request_status", "Request Status (1 = OK, 0 = NOT OK)", registry=registry),
}


# Helper Function to Fetch and Update Metrics
def fetch_and_update_metrics():
    try:
        response = requests.get(DATA_SOURCE_URL)
        response.raise_for_status()
        data = response.json()

        device_data = data.get("data", {})

        # Set request status to OK
        metrics["request_status"].set(1)

        # Update Prometheus Metrics for each phase
        for phase in ["1", "2", "3"]:
            metrics["voltage"].labels(phase=phase).set(float(device_data[f"V{phase}"]["value"]))
            metrics["current"].labels(phase=phase).set(float(device_data[f"A{phase}"]["value"]))
            metrics["power"].labels(phase=phase).set(float(device_data[f"W{phase}"]["value"]))
            metrics["energy"].labels(phase=phase).set(float(device_data[f"Wh{phase}"]["value"]))
            metrics["power_factor"].labels(phase=phase).set(float(device_data[f"PF{phase}"]["value"]))
            metrics["frequency"].labels(phase=phase).set(float(device_data[f"Fr{phase}"]["value"]))

        # Update total metrics
        metrics["total_current"].set(float(device_data["A"]["value"]))
        metrics["total_power"].set(float(device_data["W"]["value"]))
        metrics["total_energy"].set(float(device_data["TWh"]["value"]))
        metrics["temperature"].set(float(device_data["T"]["value"]))

    except requests.RequestException as e:
        logging.error(f"Failed to fetch data: {e}")
        # Set request status to NOT OK for any error
        metrics["request_status"].set(0)


@app.route("/metrics")
@cache.cached(timeout=CACHE_TIMEOUT)
def metrics_endpoint():
    fetch_and_update_metrics()
    return Response(generate_latest(registry), mimetype="text/plain")


# Run Flask App
if __name__ == "__main__":
    app.run(host="0.0.0.0", port=EXPORTER_PORT)
