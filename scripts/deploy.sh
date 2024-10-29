#!/bin/bash

# Variables
EXECUTABLE="/usr/local/bin/nlip"
PLIST_PATH="/Library/LaunchDaemons/com.nlip.plist"
BUILD_PATH="./nlip"
CERT_NAME="NLIPSigningCert"
KEYCHAIN_PASSWORD=$(cat scripts/.keychain_password)

# Build the Go project
echo "Building the Go project..."
go build -o $BUILD_PATH

if [ $? -ne 0 ]; then
  echo "Build failed. Exiting."
  exit 1
fi
echo "Build succeeded."

# Unlock the keychain
echo "Unlocking the keychain..."
security unlock-keychain -p "$KEYCHAIN_PASSWORD" ~/Library/Keychains/login.keychain-db

if [ $? -ne 0 ]; then
  echo "Failed to unlock keychain. Exiting."
  exit 1
fi

# Sign the executable
echo "Signing the executable..."
codesign -s "$CERT_NAME" ./nlip

if [ $? -ne 0 ]; then
  echo "Code signing failed. Exiting."
  exit 1
fi
echo "Code signing succeeded."

# Move the executable to the install path
echo "Moving the executable to $EXECUTABLE..."
sudo mv $BUILD_PATH $EXECUTABLE

# Set permissions for the executable
# echo "Setting permissions for the executable..."
# sudo chmod +x $EXECUTABLE
# sudo chown root:wheel $EXECUTABLE

# This seems necessary to process changes correctly
sleep 0.5

# Reload the launchd service
echo "Reloading the launchd service..."
sudo launchctl unload $PLIST_PATH 2>/dev/null
sudo launchctl load $PLIST_PATH

echo "Deployment complete. The service has been reloaded."
