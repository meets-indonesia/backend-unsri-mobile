#!/bin/bash

# Script untuk install Kubernetes dan kubectl di Ubuntu
# Usage: ./scripts/install-kubernetes.sh

set -e

echo "🚀 Installing Kubernetes and kubectl..."

# Update system
echo "📦 Updating system packages..."
sudo apt-get update

# Install required packages
echo "📦 Installing required packages..."
sudo apt-get install -y apt-transport-https ca-certificates curl gpg

# Add Kubernetes GPG key
echo "🔑 Adding Kubernetes GPG key..."
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.28/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg

# Add Kubernetes repository
echo "📚 Adding Kubernetes repository..."
echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.28/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list

# Update package list
echo "🔄 Updating package list..."
sudo apt-get update

# Install kubectl, kubelet, kubeadm
echo "📦 Installing kubectl, kubelet, kubeadm..."
sudo apt-get install -y kubectl kubelet kubeadm

# Hold packages to prevent auto-update
echo "🔒 Holding Kubernetes packages..."
sudo apt-mark hold kubelet kubeadm kubectl

# Verify installation
echo "✅ Verifying installation..."
kubectl version --client --output=yaml

echo ""
echo "✅ Kubernetes and kubectl installed successfully!"
echo ""
echo "📝 Next steps:"
echo "   1. For single-node cluster (minikube):"
echo "      - Install minikube: curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64"
echo "      - sudo install minikube-linux-amd64 /usr/local/bin/minikube"
echo "      - minikube start"
echo ""
echo "   2. For multi-node cluster:"
echo "      - Initialize master: sudo kubeadm init"
echo "      - Setup kubeconfig: mkdir -p $HOME/.kube && sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config && sudo chown $(id -u):$(id -g) $HOME/.kube/config"
echo ""
echo "   3. Install container runtime (if not using Docker):"
echo "      - Install containerd: sudo apt-get install -y containerd"
echo "      - Configure: sudo mkdir -p /etc/containerd && containerd config default | sudo tee /etc/containerd/config.toml"
echo "      - Restart: sudo systemctl restart containerd && sudo systemctl enable containerd"
echo ""

