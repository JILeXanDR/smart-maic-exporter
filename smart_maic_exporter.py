import logging
import os

import requests
from flask import Flask, Response
from prometheus_client import CollectorRegistry, Gauge, generate_latest

# Configure Logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s",
    handlers=[logging.StreamHandler()]
)

# Configuration from environment variables
DATA_SOURCE_URL = os.getenv("DATA_SOURCE_URL", "http://192.168.10.55/?page=getwdata")
EXPORTER_PORT = int(os.getenv("EXPORTER_PORT", 8000))

if not DATA_SOURCE_URL.startswith("http://") and not DATA_SOURCE_URL.startswith("https://"):
    logging.error("Invalid DATA_SOURCE_URL. Must start with 'http://' or 'https://'.")
    exit(1)

# Flask App
app = Flask(__name__)

# Initialize Prometheus Registry
registry = CollectorRegistry()

# Define Prometheus Metrics with a common prefix and labels
metric_prefix = "smart_maic_"
metrics = {
    "voltage": Gauge(f"{metric_prefix}voltage", "Voltage per line (V)", ["line"], registry=registry),
    "current": Gauge(f"{metric_prefix}current", "Current per line (A)", ["line"], registry=registry),
    "power": Gauge(f"{metric_prefix}power", "Active Power per line (W)", ["line"], registry=registry),
    "energy": Gauge(f"{metric_prefix}energy", "Energy per line (Wh)", ["line"], registry=registry),
    "power_factor": Gauge(f"{metric_prefix}power_factor", "Power Factor per line", ["line"], registry=registry),
    "frequency": Gauge(f"{metric_prefix}frequency", "Frequency per line (Hz)", ["line"], registry=registry),
    "total_current": Gauge(f"{metric_prefix}total_current", "Total Current (A)", [], registry=registry),
    "total_power": Gauge(f"{metric_prefix}total_power", "Total Active Power (W)", [], registry=registry),
    "total_energy": Gauge(f"{metric_prefix}total_energy", "Total Energy (Wh)", [], registry=registry),
    "temperature": Gauge(f"{metric_prefix}temperature", "Device Temperature (°C)", [], registry=registry),
    "device_api_status": Gauge(
        f"{metric_prefix}device_api_status",
        "Device API Status (0 = Offline, 1 = OK, 2 = Too Many Requests)",
        [],
        registry=registry
    ),
}


# Helper Function to Fetch and Update Metrics
def fetch_and_update_metrics():
    try:
        logging.info("Fetching data from the Smart Maic device...")
        response = requests.get(DATA_SOURCE_URL, timeout=3)

        # Handle HTTP 429 (Too Many Requests)
        if response.status_code == 429:
            logging.warning("Received 429 Too Many Requests from the API.")
            metrics["device_api_status"].set(2)
            return

        response.raise_for_status()
        data = response.json()

        device_data = data.get("data", {})
        if not device_data:
            logging.warning("No data found in the response. Skipping metrics update.")
            metrics["device_api_status"].set(0)
            return

        # Set request status to OK
        metrics["device_api_status"].set(1)

        # Update Prometheus Metrics for each line
        for line in ["1", "2", "3"]:
            metrics["voltage"].labels(line=line).set(float(device_data.get(f"V{line}", {}).get("value", 0)))
            metrics["current"].labels(line=line).set(float(device_data.get(f"A{line}", {}).get("value", 0)))
            metrics["power"].labels(line=line).set(float(device_data.get(f"W{line}", {}).get("value", 0)))
            metrics["energy"].labels(line=line).set(float(device_data.get(f"Wh{line}", {}).get("value", 0)))
            metrics["power_factor"].labels(line=line).set(float(device_data.get(f"PF{line}", {}).get("value", 0)))
            metrics["frequency"].labels(line=line).set(float(device_data.get(f"Fr{line}", {}).get("value", 0)))

        # Update total metrics
        metrics["total_current"].set(float(device_data.get("A", {}).get("value", 0)))
        metrics["total_power"].set(float(device_data.get("W", {}).get("value", 0)))
        metrics["total_energy"].set(float(device_data.get("TWh", {}).get("value", 0)))
        metrics["temperature"].set(float(device_data.get("T", {}).get("value", 0)))

        logging.info("Successfully updated metrics.")
    except requests.exceptions.RequestException as e:
        logging.error(f"Failed to fetch data: {e}")
        # Set request status to Offline
        metrics["device_api_status"].set(0)


# Flask Endpoint for Prometheus Metrics
@app.route("/metrics")
def metrics_endpoint():
    fetch_and_update_metrics()
    return Response(generate_latest(registry), mimetype="text/plain")


# Run Flask App
if __name__ == "__main__":
    try:
        logging.info(f"Starting Smart Maic Exporter on port {EXPORTER_PORT}...")
        app.run(host="0.0.0.0", port=EXPORTER_PORT)
    except KeyboardInterrupt:
        logging.info("Shutting down Smart Maic Exporter...")
