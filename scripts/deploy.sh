#!/bin/bash

# Variables
EXECUTABLE="/usr/local/bin/nlip"
PLIST_PATH="/Library/LaunchDaemons/com.nlip.plist"

# Build the Go project
echo "Building the Go project..."
sudo go build -o $EXECUTABLE

if [ $? -ne 0 ]; then
  echo "Build failed. Exiting."
  exit 1
fi
echo "Build succeeded."

# Set permissions for the executable
echo "Setting permissions for the executable..."
sudo chmod +x $EXECUTABLE
sudo chown root:wheel $EXECUTABLE

# This might be needed to process changes correctly
sleep 0.5

# Reload the launchd service
echo "Reloading the launchd service..."
sudo launchctl unload $PLIST_PATH 2>/dev/null
sudo launchctl load $PLIST_PATH

echo "Deployment complete. The service has been reloaded."
