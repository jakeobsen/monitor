# System Monitor with MQTT Integration

## Overview
This repository implements a Go-based system monitoring tool capable of collecting and publishing system resource metrics to an MQTT broker. It provides detailed insights into both memory and CPU usage for distributed monitoring.

## Features
This system monitor includes the following features:
- **Memory Monitoring**: Captures detailed memory statistics such as total, free, cached, available, and shared memory by leveraging data from the `/proc/meminfo` file.
- **CPU Monitoring**: Continuously tracks CPU utilization for all cores as well as the overall system by reading and processing data from `/proc/stat`. The utilization is calculated as a percentage for overall and per-core CPU activity.
- **MQTT Integration**: Publishes gathered metrics in JSON format to specific MQTT topics, enabling easy integration into subscriber-based monitoring and alert systems.
- **Modular and Extensible Design**: The codebase is designed to allow rapid inclusion of additional system monitoring features.

## Data Published

### Memory Metrics
The monitoring tool publishes memory-related metrics to the MQTT topic `memory`. The data includes:
- Total memory, representing the full system memory capacity in kilobytes.
- Used memory, calculated as total memory minus free and cached memory.
- Free memory, indicating the unused system memory in kilobytes.
- Memory used for cache purposes to enhance system performance.
- Available memory, providing an estimate of memory that can be allocated to new processes.
- Shared memory between processes.

### CPU Metrics
CPU utilization statistics are published to the MQTT topic `cpu`. This includes:
- Overall CPU usage as a percentage under the key `Total`.
- Per-core CPU usage as percentages under keys corresponding to individual CPUs, such as `CPU0`, `CPU1`, and so forth.

## Prerequisites
- **Go Language**: Version 1.19.8 or later is required for building and running the project.
- **MQTT Broker**: A functional MQTT broker is required to test and observe published system metrics. This has been tested against Mosquitto MQTT.

## How to Run
1. Clone the repository and navigate to the project directory.
2. Use the provided `build.sh` script to compile the project and package it as a Debian package.
3. Install and run the program using the systemd configuration included in the package. Start the service using `sudo systemctl start monitor.service`.
   - You may need to setup environment variables.

## Configuration
- MQTT connection settings, including the broker address, port, username, and passwors, are configured with either environment variables or command line arguments.
- MQTT topics can be adjusted to suit your requirements or existing monitoring infrastructure.

## Project Structure
- The `monitor` package contains the implementations for memory and CPU monitoring.
- The `main.go` file functions as the program's entry point, initializing and running the monitoring daemon.
- A `build.sh` script is included for Debian package creation and installation.

## Contributing
Contributions are welcome, whether for new monitoring capabilities (such as disk or network usage) or other enhancements. Pull requests and suggestions via the repository's issue tracker are gratefully accepted.

## License
This project is open source and licensed under the AGPLv3 License. For additional details, please refer to the LICENSE file.

## Contact
If you encounter any issues or have suggestions for improvement, feel free to raise an issue or discuss with the maintainers.
