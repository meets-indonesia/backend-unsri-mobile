# 🚀 Quick Install Kubernetes & kubectl

Panduan cepat untuk install Kubernetes dan kubectl di server Ubuntu.

## ✅ Prerequisites Check

Anda sudah punya:
- ✅ Docker (version 26.1.3)
- ✅ Ubuntu 20.04
- ✅ Server access

## 🎯 Quick Install

### Step 1: Install kubectl

```bash
cd ~/sinergi/backend-unsri-mobile

# Download dan jalankan script
chmod +x scripts/install-kubernetes.sh
./scripts/install-kubernetes.sh
```

Atau install manual:

```bash
# Update system
sudo apt-get update

# Install required packages
sudo apt-get install -y apt-transport-https ca-certificates curl gpg

# Add Kubernetes GPG key
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.28/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg

# Add Kubernetes repository
echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.28/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list

# Update dan install
sudo apt-get update
sudo apt-get install -y kubectl

# Verify
kubectl version --client
```

### Step 2: Install Minikube (Optional - untuk local cluster)

Jika ingin setup local Kubernetes cluster untuk testing:

```bash
chmod +x scripts/install-minikube.sh
./scripts/install-minikube.sh
```

Atau manual:

```bash
# Download Minikube
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64

# Install
sudo install minikube-linux-amd64 /usr/local/bin/minikube

# Start cluster
minikube start --driver=docker

# Verify
kubectl get nodes
```

## ✅ Verify Installation

```bash
# Check kubectl
kubectl version --client

# Check Minikube (jika install)
minikube version
minikube status
```

## 🚀 Next Steps

Setelah Kubernetes terinstall:

1. **Deploy UNSRI Backend:**
```bash
cd ~/sinergi/backend-unsri-mobile

# Create namespace
kubectl create namespace unsri-backend

# Deploy services
kubectl apply -f deployments/kubernetes/
```

2. **Lihat dokumentasi lengkap:**
   - [Kubernetes Installation Guide](./docs/KUBERNETES_INSTALLATION.md)
   - [Deployment Guide](./docs/DEPLOYMENT.md)

## 📚 Useful Commands

```bash
# kubectl
kubectl get nodes
kubectl get pods -A
kubectl cluster-info

# Minikube
minikube start
minikube stop
minikube dashboard
minikube status
```

## ❓ Troubleshooting

**kubectl: command not found**
```bash
# Reinstall
sudo apt-get install -y kubectl
```

**Minikube start fails**
```bash
# Check Docker
docker ps

# Delete and recreate
minikube delete
minikube start --driver=docker
```

Lihat [docs/KUBERNETES_INSTALLATION.md](./docs/KUBERNETES_INSTALLATION.md) untuk troubleshooting lengkap.

