# Kubernetes Installation Guide

Panduan lengkap untuk install Kubernetes dan kubectl di Ubuntu server.

## 📋 Prerequisites

- Ubuntu 20.04+ (atau Debian-based Linux)
- Docker sudah terinstall (sudah ada ✅)
- Root atau sudo access
- Minimum 2GB RAM (4GB+ recommended)
- Minimum 2 CPU cores

## 🚀 Quick Installation

### Option 1: Install kubectl Only (Recommended untuk development)

Jika Anda hanya perlu kubectl untuk manage cluster remote:

```bash
# Download script
cd ~/sinergi/backend-unsri-mobile
chmod +x scripts/install-kubernetes.sh

# Run installation
./scripts/install-kubernetes.sh
```

### Option 2: Install Minikube (Single-node cluster untuk testing)

Untuk local development/testing dengan single-node cluster:

```bash
# Install kubectl first
./scripts/install-kubernetes.sh

# Install Minikube
chmod +x scripts/install-minikube.sh
./scripts/install-minikube.sh
```

## 📦 Manual Installation

### Step 1: Install kubectl

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

# Update package list
sudo apt-get update

# Install kubectl
sudo apt-get install -y kubectl

# Verify installation
kubectl version --client
```

### Step 2: Install kubelet & kubeadm (Optional, untuk multi-node cluster)

```bash
# Install kubelet and kubeadm
sudo apt-get install -y kubelet kubeadm

# Hold packages to prevent auto-update
sudo apt-mark hold kubelet kubeadm kubectl
```

### Step 3: Install Minikube (Untuk single-node cluster)

```bash
# Download Minikube
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64

# Install Minikube
sudo install minikube-linux-amd64 /usr/local/bin/minikube

# Cleanup
rm minikube-linux-amd64

# Verify
minikube version
```

## 🎯 Setup Kubernetes Cluster

### Option A: Minikube (Recommended untuk development/testing)

**Start Minikube:**
```bash
# Start cluster dengan Docker driver
minikube start --driver=docker

# Verify cluster
kubectl cluster-info
kubectl get nodes
```

**Useful Minikube Commands:**
```bash
# Start cluster
minikube start

# Stop cluster
minikube stop

# Delete cluster
minikube delete

# Open dashboard
minikube dashboard

# Check status
minikube status

# Get cluster IP
minikube ip
```

### Option B: kubeadm (Untuk production multi-node cluster)

**Initialize Master Node:**
```bash
# Initialize cluster
sudo kubeadm init --pod-network-cidr=10.244.0.0/16

# Setup kubeconfig
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

# Install network plugin (Calico)
kubectl apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.26.1/manifests/calico.yaml

# Verify
kubectl get nodes
```

**Join Worker Nodes:**
```bash
# Di master node, generate join command
kubeadm token create --print-join-command

# Di worker node, run join command
sudo kubeadm join <master-ip>:6443 --token <token> --discovery-token-ca-cert-hash sha256:<hash>
```

## 🔧 Configure kubectl

### Setup kubeconfig

Jika menggunakan remote cluster atau kubeadm:

```bash
# Copy config dari master node
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

# Verify
kubectl cluster-info
kubectl get nodes
```

### Setup alias (Optional)

```bash
# Add to ~/.bashrc or ~/.zshrc
echo 'alias k=kubectl' >> ~/.bashrc
echo 'complete -F __start_kubectl k' >> ~/.bashrc
source ~/.bashrc
```

## ✅ Verify Installation

```bash
# Check kubectl version
kubectl version --client

# Check cluster connection
kubectl cluster-info

# List nodes
kubectl get nodes

# List all resources
kubectl get all --all-namespaces
```

## 🚀 Deploy UNSRI Backend ke Kubernetes

Setelah Kubernetes terinstall, deploy backend:

```bash
# Navigate to project
cd ~/sinergi/backend-unsri-mobile

# Create namespace
kubectl create namespace unsri-backend

# Create secrets
kubectl create secret generic unsri-secrets \
  --from-literal=jwt-secret=$(openssl rand -base64 32) \
  --from-literal=db-password=your-strong-password \
  --from-literal=redis-password=your-redis-password \
  --namespace=unsri-backend

# Deploy infrastructure
kubectl apply -f deployments/kubernetes/postgres.yaml
kubectl apply -f deployments/kubernetes/redis.yaml

# Deploy services
kubectl apply -f deployments/kubernetes/api-gateway.yaml
kubectl apply -f deployments/kubernetes/auth-service.yaml
# ... (other services)

# Check status
kubectl get pods -n unsri-backend
kubectl get svc -n unsri-backend
```

## 🔍 Troubleshooting

### kubectl: command not found

```bash
# Check if kubectl is installed
which kubectl

# If not found, reinstall
sudo apt-get install -y kubectl
```

### Cannot connect to cluster

```bash
# Check kubeconfig
kubectl config view

# Check cluster info
kubectl cluster-info

# For Minikube
minikube status
minikube start
```

### Minikube start fails

```bash
# Check Docker
docker ps

# Start with verbose logging
minikube start --driver=docker --v=7

# Delete and recreate
minikube delete
minikube start --driver=docker
```

### Permission denied

```bash
# Fix kubeconfig permissions
sudo chown $(id -u):$(id -g) $HOME/.kube/config
chmod 600 $HOME/.kube/config
```

## 📚 Additional Resources

- [Kubernetes Official Docs](https://kubernetes.io/docs/)
- [kubectl Cheat Sheet](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)
- [Minikube Documentation](https://minikube.sigs.k8s.io/docs/)
- [kubeadm Documentation](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/)

## 🎯 Next Steps

1. ✅ Install kubectl
2. ✅ Setup Kubernetes cluster (Minikube atau kubeadm)
3. ✅ Deploy UNSRI Backend services
4. ✅ Setup Ingress untuk external access
5. ✅ Configure monitoring dan logging

Lihat [Deployment Guide](./DEPLOYMENT.md) untuk panduan deployment lengkap.

