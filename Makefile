# ===================================================================
# Configuration
# ===================================================================

GO := go
BIN_DIR := ./bin
LOGS_DIR := ./logs
GATEWAY_BIN := $(BIN_DIR)/gateway
SIMULATOR_BIN := $(BIN_DIR)/simulator
GATEWAY_PID := .gateway.pid
SIMULATOR_PID := .simulator.pid
INFLUX_URL := http://localhost:8086
DOCKER_COMPOSE = docker-compose

.PHONY: all build run strart-infra stop clean logs stats help


all: run  ## (Default) Run infrastructure and application


help:  ## Show help for make command
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# -------------------------------------------------------------------------------
# Target for Build & Run Applicaions
# -------------------------------------------------------------------------------

build:  ## Compiling program gateway dan simulator
	@echo "Building applications..."
	@mkdir -p $(BIN_DIR)
	@mkdir -p $(LOGS_DIR)
	@$(GO) build -o $(GATEWAY_BIN) ./cmd/gateway/main.go
	@$(GO) build -o $(SIMULATOR_BIN) ./cmd/simulator/main.go

run: build start-infra  ## Run gateway and simulator in background
	@echo "Starting Gateway..."
	@nohup $(GATEWAY_BIN) > $(LOGS_DIR)/gateway.log 2>&1 & echo $$! > $(GATEWAY_PID)
	@echo "Gateway started with PID: $$(cat $(GATEWAY_PID))"
	@sleep 1
	@echo "Starting Simulator..."
	@nohup $(SIMULATOR_BIN) > $(LOGS_DIR)/simulator.log 2>&1 & echo $$! > $(SIMULATOR_PID)
	@echo "Simulator started with PID: $$(cat $(SIMULATOR_PID))"
	@echo ""
	@echo "ALL PROGRAM AND INFRA ARE RUNNING IN BACKGROUND"
	@echo " - InfluxDB Dashboard : $(INFLUX_URL)"
	@echo " - make logs : To see program logs"
	@echo " - make stop : To stop all services"

stop: stop-apps stop-infra  ## Stop program and infrastructure
	@echo "All services stopped."

stop-apps:  ## (Internal) Only stop Go program
	@echo "Stopping Go applications..."
	@# Stop Gateway
	@if [ -f $(GATEWAY_PID) ]; then \
		PID=$$(cat $(GATEWAY_PID)); \
		echo "Stopping Gateway (PID $$PID)..."; \
		kill $$PID 2>/dev/null || true; \
		sleep 1; \
		if ps -p $$PID > /dev/null 2>&1; then \
			kill -9 $$PID 2>/dev/null || true; \
		fi; \
		rm -f $(GATEWAY_PID); \
	fi
	@# Stop Simulator
	@if [ -f $(SIMULATOR_PID) ]; then \
		PID=$$(cat $(SIMULATOR_PID)); \
		echo "Stopping Simulator (PID $$PID)..."; \
		kill $$PID 2>/dev/null || true; \
		sleep 1; \
		if ps -p $$PID > /dev/null 2>&1; then \
			kill -9 $$PID 2>/dev/null || true; \
		fi; \
		rm -f $(SIMULATOR_PID); \
	fi
	@# Cleanup by file name, but exclude 'make' process it self
	@pgrep -f "$(GATEWAY_BIN)" | grep -v "$$$$" | xargs kill -9 2>/dev/null || true
	@pgrep -f "$(SIMULATOR_BIN)" | grep -v "$$$$" | xargs kill -9 2>/dev/null || true
	@echo "Go applications stopped."

logs: ## Show log from gateway and simulator
	@echo "Showing logs... (Press Ctrl+C to close logs view)"
	@tail -f $(LOGS_DIR)/gateway.log $(LOGS_DIR)/simulator.log

# ------------------------------------------------------------------------------
# Target For Infrastructure (Docker)
# ------------------------------------------------------------------------------

start-infra: ## Run infrastructure Docker di background
	@echo "Starting infrastructure..."
	@docker compose up -d

stop-infra: ## Stop (pause) Docker container
	@echo "Stopping infrastructure..."
	@docker compose stop

down: stop-apps ## Stop and delete all (program/apps & infrastructure)
	@echo "Tearing down infrastructure..."
	@docker compose down -v
	@echo "Cleanup complete."

# ----------------------------------------------------------------------------
# Target Status
# ----------------------------------------------------------------------------

stats: ## Check process status for applications anda infrastructure
	@if [ -f $(GATEWAY_PID) ]; then \
    	PID=$$(cat $(GATEWAY_PID)); \
        if ps -p $$PID > /dev/null; then \
            echo "Gateway is running (PID: $$PID)"; \
			echo "Uptime: $$(ps -o etime= -p $$PID)"; \
		else \
			echo "Gateway PID file exists, but process $$PID is NOT running."; \
		fi \
	else \
		echo "Gateway is not running..."; \
    fi

	@if [ -f $(SIMULATOR_PID) ]; then \
	    PID=$$(cat $(SIMULATOR_PID)); \
       if ps -p $$PID > /dev/null; then \
           echo "Simulator is running (PID: $$PID)"; \
			echo "Uptime: $$(ps -o etime= -p $$PID)"; \
		else \
			echo "Simulator PID file exists, but process $$PID is NOT running."; \
		fi \
	else \
		echo "Simulator is not running..."; \
    fi

	@if [ $$(docker compose ps -q | wc -l) -gt 0 ]; then \
		echo "Infrastructure container are UP."; \
		echo "Suggest: Run 'docker compose ps' for container details."; \
	else \
		echo "Infrastructure container is DOWN."; \
	fi

# ----------------------------------------------------------------------------
# Target for Cleanup
# ----------------------------------------------------------------------------

clean: ## Clear binary file, log, and PID
	@echo "Cleaning up build artifacts and logs..."
	@rm -rf $(BIN_DIR)
	@rm -rf $(LOGS_DIR)
	@rm -f $(GATEWAY_PID) $(SIMULATOR_PID)
