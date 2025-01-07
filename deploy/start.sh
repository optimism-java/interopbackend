#!/bin/bash

# Check Foundry is installed
if ! command -v forge &> /dev/null
then
    echo "Foundry is not installed. Please install Foundry first."
    echo "You can install Foundry using:"
    echo "curl -L https://foundry.paradigm.xyz | bash"
    exit 1
fi

if ! command -v supersim &> /dev/null
then
    echo "SuperSim is not installed. Please install SuperSim first."
    echo "Installation instructions for SuperSim..."
    exit 1
fi

# Check Foundry version
FOUNDRY_VERSION=$(forge --version)
echo "Foundry version: $FOUNDRY_VERSION"

# Starting SuperSim in background
echo "Starting SuperSim..."
nohup supersim --interop.autorelay > supersim.log 2>&1 &

# Wait for SuperSim to start completely
sleep 5

# Then start the docker-compose service
echo "Starting Docker Compose services..."
docker-compose -f docker-compose.yaml up -d

# Validate network
docker network ls | grep interop_network
