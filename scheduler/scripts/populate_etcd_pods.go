package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Pod struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"` // Added namespace field
		Labels    struct {
			App string `json:"app"`
		} `json:"labels"`
	} `json:"metadata"`
	Spec struct {
		NodeName   string `json:"nodeName"`
		Containers []struct {
			Name  string `json:"name"`
			Image string `json:"image"`
			Ports []struct {
				ContainerPort int `json:"containerPort"`
			} `json:"ports"`
			Resources struct {
				Requests struct {
					CPU    string `json:"cpu"`
					Memory string `json:"memory"`
				} `json:"requests"`
			} `json:"resources"`
		} `json:"containers"`
		Affinity struct {
			NodeAffinity struct {
				RequiredDuringSchedulingIgnoredDuringExecution struct {
					NodeSelectorTerms []struct {
						MatchExpressions []struct {
							Key      string   `json:"key"`
							Operator string   `json:"operator"`
							Values   []string `json:"values"`
						} `json:"matchExpressions"`
					} `json:"nodeSelectorTerms"`
				} `json:"requiredDuringSchedulingIgnoredDuringExecution"`
			} `json:"nodeAffinity"`
		} `json:"affinity"`
	} `json:"spec"`
}

func main() {
	// Define the pod object
	pod := Pod{
		APIVersion: "v1",
		Kind:       "Pod",
	}
	pod.Metadata.Name = "example-pod"
	pod.Metadata.Namespace = "default" // Set the namespace to default
	pod.Metadata.Labels.App = "example"
	pod.Spec.NodeName = ""
	pod.Spec.Containers = []struct {
		Name  string `json:"name"`
		Image string `json:"image"`
		Ports []struct {
			ContainerPort int `json:"containerPort"`
		} `json:"ports"`
		Resources struct {
			Requests struct {
				CPU    string `json:"cpu"`
				Memory string `json:"memory"`
			} `json:"requests"`
		} `json:"resources"`
	}{
		{
			Name:  "example-container",
			Image: "nginx:1.21.6",
			Ports: []struct {
				ContainerPort int `json:"containerPort"`
			}{
				{ContainerPort: 80},
			},
			Resources: struct {
				Requests struct {
					CPU    string `json:"cpu"`
					Memory string `json:"memory"`
				} `json:"requests"`
			}{
				Requests: struct {
					CPU    string `json:"cpu"`
					Memory string `json:"memory"`
				}{
					CPU:    "1",
					Memory: "256Mi",
				},
			},
		},
	}
	pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms = []struct {
		MatchExpressions []struct {
			Key      string   `json:"key"`
			Operator string   `json:"operator"`
			Values   []string `json:"values"`
		} `json:"matchExpressions"`
	}{
		{
			MatchExpressions: []struct {
				Key      string   `json:"key"`
				Operator string   `json:"operator"`
				Values   []string `json:"values"`
			}{
				{
					Key:      "zone",
					Operator: "In",
					Values:   []string{"us-west-1a"},
				},
			},
		},
	}

	// Marshal the pod into JSON
	podJSON, err := json.Marshal(pod)
	if err != nil {
		log.Fatalf("Failed to marshal pod: %v", err)
	}

	// Connect to the local etcd instance
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to connect to etcd: %v", err)
	}
	defer cli.Close()

	// Save the pod to etcd under a specific key
	key := "/kubernetes/pods/" + pod.Metadata.Namespace + "/" + pod.Metadata.Name
	_, err = cli.Put(context.Background(), key, string(podJSON))
	if err != nil {
		log.Fatalf("Failed to write pod to etcd: %v", err)
	}

	fmt.Println("Pod successfully written to etcd")
}
