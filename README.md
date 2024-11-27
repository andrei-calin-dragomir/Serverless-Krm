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

