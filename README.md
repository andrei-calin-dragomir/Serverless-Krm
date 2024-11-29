# Serverless-Krm

This document provides an overview of our current cluster setup, including SSH access, monitoring stack configuration, and the Kubernetes environment details.


## **Cluster Overview**
Our cluster consists of a **controller node** and **four worker nodes** configured to run a Kubernetes-based environment. The cluster utilizes **Kubernetes v1.28**, **Knative v1.16**, and **MetalLB** for load balancing.

The cluster is deployed using [Continuum](https://github.com/atlarge-research/continuum).

A monitoring stack using **Prometheus** and **Grafana** is available for performance monitoring and troubleshooting.


## **SSH Machine Administration**

### **Main Node**
The main administrative node provides access to manage the cluster:
```bash
ssh node1
```

## **SSH Access to Experimental Cluster**

### **Controller Node**
The controller node manages the cluster's operations, including scheduling and cluster control-plane tasks. Access it using:
```bash
ssh cloud_controller_dsmj3@192.168.222.2 -i /home/dsmj3/.ssh/id_rsa_continuum
```

### **Worker Nodes**
Worker nodes handle workloads and host the pods in the cluster. Each node is accessible via SSH:

1. **Worker Node 0**:
   ```bash
   ssh cloud0_dsmj3@192.168.222.3 -i /home/dsmj3/.ssh/id_rsa_continuum
   ```

2. **Worker Node 1**:
   ```bash
   ssh cloud1_dsmj3@192.168.222.4 -i /home/dsmj3/.ssh/id_rsa_continuum
   ```

3. **Worker Node 2**:
   ```bash
   ssh cloud2_dsmj3@192.168.222.5 -i /home/dsmj3/.ssh/id_rsa_continuum
   ```

4. **Worker Node 3**:
   ```bash
   ssh cloud3_dsmj3@192.168.222.6 -i /home/dsmj3/.ssh/id_rsa_continuum
   ```


## **Monitoring Stack**

### **Grafana**
Grafana is used for visualizing metrics collected by Prometheus. Access Grafana by port-forwarding:
```bash
ssh -L 3000:192.168.222.3:3000 cloud_controller_dsmj3@192.168.222.2 -i /home/dsmj3/.ssh/id_rsa_continuum
```
Once connected, open your browser and navigate to:
```
http://localhost:3000
```

### **Prometheus**
Prometheus is used for collecting and querying metrics from the cluster. Access Prometheus by port-forwarding:
```bash
ssh -L 9090:192.168.222.3:9090 cloud_controller_dsmj3@192.168.222.2 -i /home/dsmj3/.ssh/id_rsa_continuum
```
Once connected, open your browser and navigate to:
```
http://localhost:9090
```

## **Cluster Setup**

### **Kubernetes**
The cluster runs on **Kubernetes v1.28**, providing container orchestration for workloads. The controller and worker nodes are configured to maintain high availability and scalability.

### **Knative**
We use **Knative v1.16** to enable serverless workloads and event-driven architectures. Knative Serving is deployed with the **Kourier** networking layer to handle ingress traffic efficiently.

### **MetalLB**
**MetalLB** is configured to provide load balancing for `LoadBalancer` services in the local environment. The IP pool is defined as `192.168.1.240-192.168.1.250`, and the `L2Advertisement` ensures smooth IP allocation across the cluster.

## **Testing** TODO

Here's the updated documentation in **GitHub Markdown syntax**, with the payload provided as a YAML file for both `curl` and programmatic examples.

---

# **Interacting with the Kubernetes API**

This guide demonstrates how to interact with the Kubernetes API **from outside the node (node1 in our case)** where the API server resides, using both `curl` and a program.

---
## **Configuration**

- **CA Certificate Location**: `/home/dsmj3/certs/k8/ca.crt`
- **Knative Service Account Token**: Stored in `.env` file at `/home/dsmj3/.env`

### **.env File Content**
```plaintext
# /home/dsmj3/.env
K8S_API_TOKEN=<YOUR_TOKEN>
```
## **1. Using `curl`**

### **Prerequisites**
- **CA Certificate (`kubernetes-ca.crt`)**: Verifies the API server's identity.
- **Bearer Token (`TOKEN`)**: Authenticates the request.
   
- **Kubernetes API Server URL**: E.g., `https://192.168.222.2:6443`.
- **YAML Payload**: Save the payload in a file, e.g., `service.yaml`.

### **Example: YAML Payload for a Knative Service**

Save this YAML file as `service.yaml`:

```yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: hello-world
  namespace: default
spec:
  template:
    spec:
      containers:
      - image: gcr.io/knative-samples/helloworld-go
        env:
        - name: TARGET
          value: "Hello World from Knative!"
```

### **Test the Kubernetes API**

```bash
curl --cacert kubernetes-ca.crt \
     -H "Authorization: Bearer $TOKEN" \
     -X GET https://192.168.222.2:6443/api
```

### **Create the Knative Service**

Use the YAML file as input with `-d @service.yaml`:

```bash
curl --cacert kubernetes-ca.crt \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/yaml" \
     -X POST https://192.168.222.2:6443/apis/serving.knative.dev/v1/namespaces/default/services \
     --data-binary @service.yaml
```

---

## **2. Using a Program**

### **Prerequisites**
- **Install Required Libraries**:
  - Python: `pip install requests PyYAML`
  - Node.js: `npm install axios yaml dotenv`
- **CA Certificate (`ca.crt`)**: Stored locally.
- **Bearer Token (`TOKEN`)**: Provided via environment variables.
- **YAML Payload**: Save the file as `service.yaml` (same as above).

### **Python Example**

Save the token in an `.env` file:

```plaintext
# .env
K8S_API_TOKEN=<YOUR_TOKEN>
```

Python script to send the request:

```python
import requests
import yaml
from dotenv import load_dotenv
import os

# Load the .env file
load_dotenv()

# API server and token
api_server = "https://<API_SERVER_IP>:6443"
token = os.getenv("K8S_API_TOKEN")

# Read the YAML payload
with open("service.yaml", "r") as file:
    payload = yaml.safe_load(file)

# Set headers
headers = {
    "Authorization": f"Bearer {token}",
    "Content-Type": "application/yaml"
}

# Send the POST request
response = requests.post(
    f"{api_server}/apis/serving.knative.dev/v1/namespaces/default/services",
    headers=headers,
    json=payload,
    verify="kubernetes-ca.crt"
)

# Print the response
print("Status Code:", response.status_code)
print("Response:", response.json())
```

---