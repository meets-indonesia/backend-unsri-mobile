#!/bin/bash

# Script untuk install Minikube (single-node Kubernetes cluster)
# Usage: ./scripts/install-minikube.sh

set -e

echo "🚀 Installing Minikube..."

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl not found. Please install kubectl first:"
    echo "   ./scripts/install-kubernetes.sh"
    exit 1
fi

# Download Minikube
echo "📥 Downloading Minikube..."
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64

# Install Minikube
echo "📦 Installing Minikube..."
sudo install minikube-linux-amd64 /usr/local/bin/minikube

# Cleanup
rm minikube-linux-amd64

# Verify installation
echo "✅ Verifying installation..."
minikube version

# Start Minikube
echo "🚀 Starting Minikube cluster..."
minikube start --driver=docker

# Verify cluster
echo "✅ Verifying cluster..."
kubectl cluster-info
kubectl get nodes

echo ""
echo "✅ Minikube installed and started successfully!"
echo ""
echo "📝 Useful commands:"
echo "   - Start cluster: minikube start"
echo "   - Stop cluster: minikube stop"
echo "   - Delete cluster: minikube delete"
echo "   - Dashboard: minikube dashboard"
echo "   - Status: minikube status"
echo ""

