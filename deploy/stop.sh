#!/bin/bash

# Stop Docker Compose services
echo "Stopping Docker Compose services..."
docker-compose -f docker-compose.yaml down

# Stop SuperSim process
echo "Stopping SuperSim..."
pkill -f supersim

# Remove the SuperSim log file
if [ -f supersim.log ]; then
    rm supersim.log
fi

# Ask about deleting mysql folder
read -p "Do you want to delete the ./mysql folder? (y/n): " answer
if [ "$answer" = "y" ] || [ "$answer" = "Y" ]; then
    echo "Deleting ./mysql folder..."
    rm -rf ./mysql
    echo "MySQL folder deleted."
fi

echo "All services have been stopped."
