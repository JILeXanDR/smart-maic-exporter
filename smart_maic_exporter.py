import logging
import os
import sys

import requests
from flask import Flask, Response
from prometheus_client import CollectorRegistry, Gauge, generate_latest

# Configure Logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s"
)

# Configuration from environment variables
DATA_SOURCE_URL = os.getenv("DATA_SOURCE_URL")
EXPORTER_PORT = os.getenv("EXPORTER_PORT", "8000")

# Validate Configuration
if not DATA_SOURCE_URL or not (DATA_SOURCE_URL.startswith("http://") or DATA_SOURCE_URL.startswith("https://")):
    logging.error("Invalid or missing DATA_SOURCE_URL. Must start with 'http://' or 'https://'.")
    sys.exit(1)

try:
    EXPORTER_PORT = int(EXPORTER_PORT)
    if EXPORTER_PORT <= 0 or EXPORTER_PORT > 65535:
        raise ValueError("Port must be between 1 and 65535")
except ValueError as e:
    logging.error(f"Invalid EXPORTER_PORT: {EXPORTER_PORT}. {e}")
    sys.exit(1)

# Flask App
app = Flask(__name__)

# Initialize Prometheus Registry
registry = CollectorRegistry()

# Define Prometheus Metrics with a common prefix
metric_prefix = "smart_maic_"

metrics = {
    # Gauges
    "voltage": Gauge(f"{metric_prefix}voltage", "Voltage (V)", ["line"], registry=registry),
    "current": Gauge(f"{metric_prefix}current", "Current (A)", ["line"], registry=registry),
    "power_factor": Gauge(f"{metric_prefix}power_factor", "Power Factor", ["line"], registry=registry),
    "frequency": Gauge(f"{metric_prefix}frequency", "Frequency (Hz)", ["line"], registry=registry),
    "temperature": Gauge(f"{metric_prefix}temperature", "Device Temperature (°C)", registry=registry),
    "energy": Gauge(f"{metric_prefix}energy", "Energy (Wh)", ["line"], registry=registry),
    "total_energy": Gauge(f"{metric_prefix}total_energy", "Total Energy (Wh)", registry=registry),
}


# Helper Function to Fetch and Update Metrics
def fetch_and_update_metrics():
    logging.info("Fetching data from the Smart Maic device...")
    try:
        response = requests.get(DATA_SOURCE_URL, timeout=3)
        response.raise_for_status()
        data = response.json()

        device_data = data.get("data", {})
        if not device_data:
            logging.warning("No data found in the response. Skipping metrics update.")
            return

        # Update Prometheus Metrics for each line
        for line in ["1", "2", "3"]:
            metrics["voltage"].labels(line=line).set(float(device_data[f"V{line}"]["value"]))  # Last voltage
            metrics["current"].labels(line=line).set(float(device_data[f"A{line}"]["value"]))  # Last current
            metrics["power_factor"].labels(line=line).set(float(device_data[f"PF{line}"]["value"]))  # Last power factor
            metrics["frequency"].labels(line=line).set(float(device_data[f"Fr{line}"]["value"]))  # Last frequency
            metrics["energy"].labels(line=line).set(float(device_data[f"Wh{line}"]["value"]))  # Last total energy

        # Update general metrics
        metrics["temperature"].set(float(device_data["T"]["value"]))  # Last temperature of device itself
        metrics["total_energy"].set(float(device_data["TWh"]["value"]))  # Last total energy

        logging.info("Successfully updated metrics.")

    except requests.exceptions.RequestException as e:
        logging.error(f"Failed to fetch data from {DATA_SOURCE_URL}: {e}")
    except KeyError as e:
        logging.error(f"Missing expected key in the response data: {e}")
    except ValueError as e:
        logging.error(f"Error processing data: {e}")


@app.route("/metrics")
def metrics_endpoint():
    fetch_and_update_metrics()
    return Response(generate_latest(registry), mimetype="text/plain")


# Run Flask App
if __name__ == "__main__":
    logging.info(f"Starting Smart Maic Exporter on port {EXPORTER_PORT}...")
    app.run(host="0.0.0.0", port=EXPORTER_PORT)
