# **IoT Gateway & Simulator**

A high-performance, end-to-end IoT simulation system featuring a **Multi-protocol Gateway** and a **Device Simulator**. This project demonstrates modern Go patterns, secure communication, and automated infrastructure provisioning.

---

## **🏛️ System Architecture**

The system follows a modern IoT telemetry pipeline:

![System Architecture](assets/Architecture-diagram-iot-gateway-sim.svg)

1. **Simulator:** Generates mock sensor data (Temperature/Voltage).
2. **Gateway:** Ingests data via TCP/UDP listeners.
3. **Broker:** Eclipse Mosquitto handles message distribution.
4. **Time-Series Stack:** Telegraf consumes MQTT data and persists it to InfluxDB.


___

## **✨ Key Features**

* **Multi-Protocol Ingestion:** Concurrent **TCP** and **UDP** listeners for diverse sensor simulations.
* **Buffered Processing:** Implements internal batching to optimize network throughput and reduce broker overhead.
* **End-to-End Security:** Strict **mTLS (Mutual TLS)** implementation between the Gateway and MQTT Broker for encrypted, authenticated communication.
* **Automated Provisioning:** Zero-config dashboard setup using InfluxDB Templates and shell automation.
* **Advanced Concurrency:** Controlled lifecycle management using `context.Context` and `sync.WaitGroup` to ensure Graceful Shutdown and prevent goroutine leaks.
* **Production-Ready Structure:** Clean Architecture following the standard Go project layout.

___

## **📂 Directory Structure**

```text
├── cmd/                # Entry points for Gateway and Simulator
├── internal/           # Private business logic (Ingestion, Processing, Transport)
├── pkg/                # Public shared models
├── infra/              # Infrastructure configs (Mosquitto, Telegraf, InfluxDB)
├── gateway-certs/      # mTLS Client certificates
├── bin/                # Compiled binaries
├── logs/               # Application runtime logs
└── Makefile            # Orchestration and automation
```

___

## **📊 Dashboard Preview**

The system includes a pre-configured InfluxDB dashboard to monitor sensor metrics in real-time.

![Iot Dashboard Monitoring](assets/iot-dasboard-monitoring.png)

- InfluxDB Dashboard : http://localhost:8086
- Username : admin
- Password : adminpassword
___

## **🔧 Prerequisites**

* **Go** (v1.22 or higher)
* **Docker** & **Docker Compose**
* **Make** (Standard on Linux/WSL)

___

## **🚀 Quick Start**

**1. Clone the Repository**

```bash
  git clone https://github.com/Nuryanto19/iot-gateway-sim.git
  cd iot-gateway-sim
```

**2. Generate Secure Certificates (mTLS)** Initialize the Certificate Authority and generate keys for both the Broker and the Gateway.

```bash
chmod +x generate-certs.sh
./generate-certs.sh
```

**3. Run the Entire Stack** Build the Go binaries, spin up the Docker infrastructure, and start the services in the background.

```bash
make run
```

**4. Operations & Monitoring**

* **Check System Health:** `make stats` (Verify PIDs and Container status)
* **View Real-time Logs:** `make logs` 
* **Stop Everything:** `make stop`

**5. Access InfluxDB**
Navigate to `http://localhost:8086`. The dashboard is automatically provisioned upon startup via `infra/influxdb/provisioning/init.sh`.

___

## **🎛️ Makefile Reference**

| **Command**             | **Action**                                              |
|:----------------------- |:------------------------------------------------------- |
| `make all` / `make run` | Build binaries, start infra, and launch apps.           |
| `make build`            | Compile Go source code into bin/.                       |
| `make start-infra`      | Launch Docker containers (Mosquitto, Influx, Telegraf). |
| `make stats`            | Display health status of Go processes and Docker.       |
| `make logs`             | Tail logs for both Gateway and Simulator.               |
| `make clean`            | Remove binaries, PIDs, and log files.                   |
| `make help`             | Showing available command 
___

## **📄 License**

This project is licensed under the **MIT License**.
