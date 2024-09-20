# Platform Monitoring Service

This project is a platform monitoring service built using Go, Prometheus, and Docker. It aims to collect and visualize metrics from multiple web applications to ensure optimal performance and reliability.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Metrics](#metrics)
- [Future Work](#future-work)
- [Contributing](#contributing)
- [License](#license)

## Overview

The service collects CPU, memory, and visitor count metrics from multiple web applications (`webapp1` and `webapp2`) using Go and exposes these metrics to Prometheus for monitoring and visualization in Grafana.

## Features

- Collects CPU and memory usage metrics from Docker containers.
- Tracks visitor counts for web applications.
- Integrates with Prometheus for metrics scraping and storage.
- Displays metrics on Grafana dashboards.

## Architecture

The platform consists of the following components:

- **Go Service**: Collects metrics from the target containers.
- **Prometheus**: Scrapes metrics exposed by the Go service.
- **Grafana**: Visualizes metrics data.

## Prerequisites

- Docker and Docker Compose
- Go (version 1.21 or higher)
- Prometheus
- Grafana

## Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/your-username/platform-monitoring.git
    cd platform-monitoring
    ```

2. Build and run the Docker containers:
    ```bash
    docker-compose up --build
    ```

3. Access Prometheus at [http://localhost:9090](http://localhost:9090).

4. Access Grafana at [http://localhost:3000](http://localhost:3000).

## Usage

1. **Run the Go Monitoring Service**:
   The Go service collects metrics from the specified web applications and exposes them at `/metrics` endpoint.

2. **Check Metrics in Prometheus**:
   Prometheus scrapes metrics data from the Go service and stores it for analysis.

3. **Visualize in Grafana**:
   Create dashboards in Grafana to visualize CPU, memory, and visitor count metrics.

## Metrics

- **CPU Usage**: `container_cpu_usage_seconds_total`
- **Memory Usage**: `container_memory_usage_bytes`
- **Visitor Count**: `webapp_visitor_count`

## Future Work

- Implement more granular metrics for web application performance.
- Add alerting rules in Prometheus for abnormal behavior detection.
- Explore further integrations for enhanced monitoring capabilities.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License.
