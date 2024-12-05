#!/bin/bash

# Check if script is being run with root privileges
if [ "$(id -u)" -ne 0 ]; then
  echo "Please run this script as root or with sudo."
  exit 1
fi

# Set etcd version (You can change the version here if needed)
ETCD_VERSION="v3.5.5"
ETCD_DOWNLOAD_URL="https://github.com/etcd-io/etcd/releases/download/${ETCD_VERSION}/etcd-${ETCD_VERSION}-darwin-arm64.zip"

# Set directory for installation
INSTALL_DIR="/usr/local/bin"

# Download the etcd binary
echo "Downloading etcd ${ETCD_VERSION} for macOS..."
curl -LO $ETCD_DOWNLOAD_URL

# Unzip the downloaded file
echo "Extracting etcd binaries..."
unzip -q "etcd-${ETCD_VERSION}-darwin-arm64.zip" -d /tmp

# Move the binaries to /usr/local/bin
echo "Installing etcd binaries to $INSTALL_DIR..."
mv /tmp/etcd-${ETCD_VERSION}-darwin-arm64/etcd /tmp/etcd-${ETCD_VERSION}-darwin-arm64/etcdctl $INSTALL_DIR

# Clean up downloaded and extracted files
echo "Cleaning up..."
rm -rf /tmp/etcd-${ETCD_VERSION}-darwin-arm64*
rm -f "etcd-${ETCD_VERSION}-darwin-arm64.zip"

# Check if installation was successful
if command -v etcd &>/dev/null && command -v etcdctl &>/dev/null; then
  echo "etcd and etcdctl installed successfully!"
else
  echo "Failed to install etcd. Please check the logs above."
  exit 1
fi

# Start etcd server in the background
echo "Starting etcd server..."
etcd &

# Check if etcd is running
sleep 2  # Give it a second to start
if pgrep -x "etcd" > /dev/null; then
  echo "etcd is running on http://127.0.0.1:2379"
else
  echo "Failed to start etcd server."
  exit 1
fi

# Provide usage instructions
echo "To interact with etcd, use 'etcdctl' commands."
echo "Example: 'etcdctl put mykey 'Hello, etcd!'"

# to stop etcd run pkill etcd
# once installed run etcd by calling etcd &