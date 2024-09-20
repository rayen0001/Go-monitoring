# Use an official Python image as the base
FROM python:3.9-slim

# Set the working directory to /app
WORKDIR /app

# Install any necessary Python dependencies
RUN pip install flask
RUN pip install prometheus_client

# Copy the application code
COPY . .

# Expose the port for the web app
EXPOSE 8081

# Run the command to start the Flask app
CMD ["python", "app.py"]
