package main

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"time"

// 	clientv3 "go.etcd.io/etcd/client/v3"
// 	corev1 "k8s.io/api/core/v1"
// 	"k8s.io/apimachinery/pkg/api/resource"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// )

// func main() {
// 	// Set up etcd client with a custom dial timeout (5 seconds)
// 	etcdClient, err := clientv3.New(clientv3.Config{
// 		Endpoints:   []string{"http://127.0.0.1:2379"}, // Adjust to your etcd endpoint
// 		DialTimeout: 5 * time.Second,                   // Custom dial timeout of 5 seconds
// 	})
// 	if err != nil {
// 		log.Fatalf("Failed to create etcd client: %v", err)
// 	}
// 	defer etcdClient.Close()

// 	// Create sample node objects with different characteristics
// 	nodes := []corev1.Node{
// 		{
// 			ObjectMeta: metav1.ObjectMeta{
// 				Name:   "node-unschedulable",
// 				Labels: map[string]string{"role": "worker", "zone": "us-west-1a"},
// 			},
// 			Spec: corev1.NodeSpec{
// 				Unschedulable: true,
// 				ProviderID:    "provider-id-unschedulable",
// 			},
// 		},
// 		{
// 			ObjectMeta: metav1.ObjectMeta{
// 				Name:   "node-tainted",
// 				Labels: map[string]string{"role": "worker", "zone": "us-west-1b"},
// 			},
// 			Spec: corev1.NodeSpec{
// 				Taints: []corev1.Taint{
// 					{
// 						Key:    "dedicated",
// 						Value:  "special",
// 						Effect: corev1.TaintEffectNoSchedule,
// 					},
// 				},
// 				ProviderID: "provider-id-tainted",
// 			},
// 		},
// 		{
// 			ObjectMeta: metav1.ObjectMeta{
// 				Name:   "node-no-affinity",
// 				Labels: map[string]string{"role": "worker", "zone": "us-west-1c"},
// 			},
// 			Spec: corev1.NodeSpec{
// 				ProviderID: "provider-id-no-affinity",
// 			},
// 		},
// 		{
// 			ObjectMeta: metav1.ObjectMeta{
// 				Name:   "node-low-resources",
// 				Labels: map[string]string{"role": "worker", "zone": "us-west-1d"},
// 			},
// 			Status: corev1.NodeStatus{
// 				Allocatable: corev1.ResourceList{
// 					corev1.ResourceCPU:    resource.MustParse("500m"),  // 0.5 CPU
// 					corev1.ResourceMemory: resource.MustParse("128Mi"), // 128 MiB memory
// 				},
// 			},
// 			Spec: corev1.NodeSpec{
// 				ProviderID: "provider-id-low-resources",
// 			},
// 		},
// 		{
// 			ObjectMeta: metav1.ObjectMeta{
// 				Name:   "node-valid",
// 				Labels: map[string]string{"role": "worker", "zone": "us-west-1a"},
// 			},
// 			Status: corev1.NodeStatus{
// 				Allocatable: corev1.ResourceList{
// 					corev1.ResourceCPU:    resource.MustParse("2"),   // 2 CPUs
// 					corev1.ResourceMemory: resource.MustParse("4Gi"), // 4 GiB memory
// 				},
// 			},
// 			Spec: corev1.NodeSpec{
// 				ProviderID: "provider-id-valid",
// 			},
// 		},
// 		{
// 			ObjectMeta: metav1.ObjectMeta{
// 				Name:   "node-valid2",
// 				Labels: map[string]string{"role": "worker", "zone": "us-west-1a"},
// 			},
// 			Status: corev1.NodeStatus{
// 				Allocatable: corev1.ResourceList{
// 					corev1.ResourceCPU:    resource.MustParse("2"),   // 2 CPUs
// 					corev1.ResourceMemory: resource.MustParse("3Gi"), // 4 GiB memory
// 				},
// 			},
// 			Spec: corev1.NodeSpec{
// 				ProviderID: "provider-id-valid",
// 			},
// 		},
// 		{
// 			ObjectMeta: metav1.ObjectMeta{
// 				Name:   "node-valid3",
// 				Labels: map[string]string{"role": "worker", "zone": "us-west-1a"},
// 			},
// 			Status: corev1.NodeStatus{
// 				Allocatable: corev1.ResourceList{
// 					corev1.ResourceCPU:    resource.MustParse("3"),   // 2 CPUs
// 					corev1.ResourceMemory: resource.MustParse("5Gi"), // 4 GiB memory
// 				},
// 			},
// 			Spec: corev1.NodeSpec{
// 				ProviderID: "provider-id-valid",
// 			},
// 		},
// 	}

// 	// Insert nodes into etcd with unique keys
// 	for _, node := range nodes {
// 		// Marshal the node into JSON format
// 		nodeJson, err := json.Marshal(node)
// 		if err != nil {
// 			log.Printf("Error marshalling node %s: %v", node.Name, err)
// 			continue
// 		}

// 		// Define the etcd key for the node (e.g., /kubernetes/nodes/node-1)
// 		nodeKey := fmt.Sprintf("/kubernetes/nodes/%s", node.Name)

// 		// Insert the node JSON into etcd
// 		_, err = etcdClient.Put(context.Background(), nodeKey, string(nodeJson))
// 		if err != nil {
// 			log.Printf("Failed to insert node %s into etcd: %v", node.Name, err)
// 		} else {
// 			fmt.Printf("Successfully added node '%s' to etcd\n", node.Name)
// 		}
// 	}
// }
