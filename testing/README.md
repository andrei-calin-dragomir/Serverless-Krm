# **README: Testing Services on Your Kubernetes Cluster**

This guide explains various ways to test services deployed in your Kubernetes cluster, including services exposed through different types (e.g., `ClusterIP`, `NodePort`, `LoadBalancer`) and Knative services.

---

## **1. Testing ClusterIP Services**
ClusterIP services are only accessible within the Kubernetes cluster.

### **Method: Test from Within the Cluster**
1. **Deploy a Test Pod**:
   Use a temporary pod (e.g., `curl`) to access the service:
   ```bash
   kubectl run curlpod --image=radial/busyboxplus:curl -i --tty
   ```
2. **Access the Service**:
   Inside the test pod, curl the service's `ClusterIP`:
   ```bash
   curl <ClusterIP>:<PORT>
   ```
   Replace `<ClusterIP>` and `<PORT>` with the service's IP and port from:
   ```bash
   kubectl get service <service-name>
   ```

3. **Exit the Test Pod**:
   After testing, delete the pod:
   ```bash
   kubectl delete pod curlpod
   ```

---

## **2. Testing NodePort Services**
NodePort services expose a service on a port of every cluster node.

### **Method: Access Using Node IP and NodePort**
1. Get the NodePort assigned to your service:
   ```bash
   kubectl get service <service-name>
   ```
   Example output:
   ```
   NAME         TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)           AGE
   my-service   NodePort   10.96.0.1       <none>        80:32000/TCP      5m
   ```

2. Get the IP of a node:
   ```bash
   kubectl get nodes -o wide
   ```

3. Test the service using the node's IP and NodePort:
   ```bash
   curl http://<Node-IP>:<NodePort>
   ```

---

## **3. Testing LoadBalancer Services**
LoadBalancer services expose a service through an external IP, typically provisioned by a cloud provider or MetalLB in local setups.

### **Method: Access Using External IP**
1. Check if the service has an external IP assigned:
   ```bash
   kubectl get service <service-name>
   ```
   Example output:
   ```
   NAME         TYPE           CLUSTER-IP      EXTERNAL-IP     PORT(S)        AGE
   my-service   LoadBalancer   10.96.0.1       192.168.1.100   80:31245/TCP   5m
   ```

2. Test the service using the external IP:
   ```bash
   curl http://192.168.1.100
   ```

---

## **4. Testing Knative Services**
Knative services use host-based routing with custom domains.

### **Method 1: Using a Custom Domain**
1. Configure a wildcard DNS record for your domain (e.g., `*.knative.example.com`) pointing to the external IP of your Knative networking layer (e.g., Kourier):
   ```
   *.knative.example.com A <EXTERNAL-IP>
   ```

2. Deploy a test Knative service:
   ```yaml
   apiVersion: serving.knative.dev/v1
   kind: Service
   metadata:
     name: helloworld
     namespace: default
   spec:
     template:
       spec:
         containers:
         - image: gcr.io/knative-samples/helloworld-go
           env:
           - name: TARGET
             value: "Knative World"
   ```

   Apply it:
   ```bash
   kubectl apply -f helloworld.yaml
   ```

3. Access the service using its hostname:
   ```bash
   curl http://helloworld.default.knative.example.com
   ```

### **Method 2: Using External IP Without DNS** (Not Implemented)
1. Retrieve the external IP of your Knative networking layer:
   ```bash
   kubectl get service kourier -n kourier-system
   ```
   Example output:
   ```
   NAME      TYPE           CLUSTER-IP      EXTERNAL-IP     PORT(S)        AGE
   kourier   LoadBalancer   10.97.136.177   192.168.1.240   80:30921/TCP   10m
   ```

2. Replace the domain part of the Knative service URL with the external IP:
   ```bash
   curl -H "Host: helloworld.default.knative.example.com" http://192.168.1.240
   ```

---

## **5. Testing Services Using Port Forwarding**
Port forwarding allows you to access services without exposing them externally.

### **Method: Forward a Local Port to the Service**
1. Forward a local port to the service:
   ```bash
   kubectl port-forward service/<service-name> <LOCAL_PORT>:<SERVICE_PORT>
   ```
   Example:
   ```bash
   kubectl port-forward service/my-service 8080:80
   ```

2. Access the service locally:
   ```bash
   curl http://localhost:8080
   ```

---

## **6. Testing with Logs and Events**
If a service isn’t responding, inspect the logs and events.

### **Check Logs**
Inspect the logs of pods backing the service:
```bash
kubectl logs -l app=<app-label>
```

### **Check Events**
Review cluster events for errors:
```bash
kubectl get events -n <namespace>
```

---

## **7. Testing Autoscaling for Knative Services**
Knative services scale pods based on traffic.

1. Generate traffic:
   ```bash
   for i in {1..100}; do curl -s http://helloworld.default.knative.example.com; done
   ```

2. Monitor pod scaling:
   ```bash
   kubectl get pods -n default
   ```

---

## **8. Testing with Ingress or API Gateway**
If your cluster uses an ingress controller (e.g., NGINX, Istio), test through the ingress gateway.

1. Check the ingress configuration:
   ```bash
   kubectl get ingress
   ```

2. Access the ingress endpoint:
   ```bash
   curl http://<INGRESS-HOST>/<PATH>
   ```

---

## **9. Testing with External Tools**
- **Browser**: Open the service URL in a browser.
- **Postman**: Test APIs with custom headers or authentication.
- **cURL/Wget**: Automate requests and observe responses.

---

By following these methods, you can comprehensively test services deployed in your Kubernetes cluster. If you encounter issues, refer to logs, events, and configuration files to debug further.