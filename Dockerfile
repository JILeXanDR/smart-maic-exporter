# Use Python base image
FROM python:3.9-slim

# Set working directory
WORKDIR /app

# Copy requirements and install dependencies
COPY requirements.txt requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

# Copy application files
COPY smart_maic_exporter.py .

# Set environment variables for Flask
ENV FLASK_APP=smart_maic_exporter.py

# Expose port
EXPOSE 8000

# Command to run the application
CMD ["python", "smart_maic_exporter.py"]
