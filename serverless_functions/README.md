# Manual

This README file contains all the templates of the functions as well as the method to deploy them (from source -> to docker image -> to deployment on k8s)

## Notations

Our image repository is located at `andreicalindragomir/serverless-krm`
All images we store are stored under a `tag`. Each `tag` in our case is a serverless function.

## Deployment

### Namespace
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: knative-control-plane
```

### Creating and pushing an image

As an example, from your function's directory, you can run:

```bash
docker buildx build --platform linux/arm64,linux/amd64 -t {function_name}:v0.1 .

docker tag {function_name}:v0.1 andreicalindragomir/serverless-krm:{function_name}

docker push andreicalindragomir/serverless-krm:{function_name}
```

### Setting up the function .yaml configuration

In the deployment of each function:
```yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: <function_name>-service
  namespace: knative-control-plane
spec:
  template:
    metadata:
      annotations:
        autoscaling.knative.dev/minScale: "1"
        autoscaling.knative.dev/minScale: "10"
    spec:
      containers:
      - name: <function_name>-container
        image: andreicalindragomir/serverless-krm:<function_name>
        ports:
        - containerPort: 8080
        resources:
          limits:
            memory: "512Mi"
            cpu: "500m"
          requests:
            memory: "256Mi"
            cpu: "250m"
```


