#!/bin/bash

# Variables
EXECUTABLE="/usr/local/bin/nlip"
PLIST_PATH="/Library/LaunchDaemons/com.nlip.plist"
BUILD_PATH="./nlip"
CERT_NAME="NLIPSigningCert"
KEYCHAIN_PASSWORD=$(cat scripts/.keychain_password)


if [ -f .env ]; then
  echo "[INFO] .env file found. Exporting environment variables."
else
  echo "[ERROR] .env file not found in the current directory. Exiting."
  exit 1
fi

# Persist variables to /etc/environment
while IFS= read -r line; do
  # Skip comments and empty lines
  if [[ ! "$line" =~ ^# ]] && [[ "$line" =~ = ]]; then
    key=$(echo "$line" | cut -d '=' -f 1)
    if grep -q "^$key=" /etc/environment; then
      # Update existing variable
      sudo sed -i '' "s|^$key=.*|$line|" /etc/environment
    else
      # Add new variable
      sudo sh -c "echo $line >> /etc/environment"
    fi
  fi
done < .env

echo "[SUCCESS] Environment variables persisted to /etc/environment."
echo "Step 2: Building the Go project..."
go build -o $BUILD_PATH

if [ $? -ne 0 ]; then
  echo "[ERROR] Build failed. Exiting."
  exit 1
fi
echo "[SUCCESS] Build succeeded."

echo "Step 3: Unlocking the keychain..."
security unlock-keychain -p "$KEYCHAIN_PASSWORD" ~/Library/Keychains/login.keychain-db

if [ $? -ne 0 ]; then
  echo "[ERROR] Failed to unlock keychain. Exiting."
  exit 1
fi
echo "[SUCCESS] Keychain unlocked."

echo "Step 4: Signing the executable with certificate '$CERT_NAME'..."
codesign -s "$CERT_NAME" ./nlip

if [ $? -ne 0 ]; then
  echo "[ERROR] Code signing failed. Exiting."
  exit 1
fi
echo "[SUCCESS] Code signing completed."

echo "Step 5: Moving the executable to $EXECUTABLE..."
sudo mv $BUILD_PATH $EXECUTABLE

if [ $? -ne 0 ]; then
  echo "[ERROR] Failed to move the executable. Exiting."
  exit 1
fi
echo "[SUCCESS] Executable moved to $EXECUTABLE."
