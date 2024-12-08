package function

import (
	// "context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"time"

	// "go.etcd.io/etcd/client/v3"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/api/core/v1"
)

// var etcdClient *clientv3.Client
var etcd_timeout = 5 * time.Second

func init() {
	// var err error
	// // Initialize etcd client
	// etcdClient, err = clientv3.New(clientv3.Config{
	// 	//TODO: Add etcd address here
	// 	Endpoints:   []string{"http://127.0.0.1:2379"},
	// 	DialTimeout: 5,                                // Timeout for connection
	// })
	// if err != nil {
	// 	log.Fatalf("Failed to connect to etcd: %v", err)
	// }
}

func ParseJSON(r *http.Request, target interface{}) error {
	defer r.Body.Close()
	
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("invalid JSON: %v", err)
	}
	return nil
}

func CreatePod(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming request
	var pod v1.Pod
	if err := ParseJSON(r, &pod); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	slog.Info("Validating required fields")
	// Validate required fields
	if pod.ObjectMeta.Name == "" || pod.ObjectMeta.Namespace == "" {
		http.Error(w, "Missing required fields: name or namespace", http.StatusBadRequest)
		return
	}

	// // Serialize Pod into JSON
	// podData, err := json.Marshal(pod)
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("Failed to serialize Pod: %v", err), http.StatusInternalServerError)
	// 	return
	// }

	// Construct the etcd key for the Pod
	slog.Info("Constructing etcd key for the Pod")
	etcdKey := fmt.Sprintf("/registry/pods/%s/%s", pod.ObjectMeta.Namespace, pod.ObjectMeta.Name)

	// Store the serialized Pod in etcd
	// ctx, cancel := context.WithTimeout(context.Background(), etcd_timeout)
	// defer cancel()

	// _, err = etcdClient.Put(ctx, etcdKey, string(podData))
	// fmt.Sprintf("Pod %q created successfully in etcd with key %q", pod.ObjectMeta.Name, etcdKey)
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("Failed to store Pod in etcd: %v", err), http.StatusInternalServerError)
	// 	return
	// }

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Pod %q created successfully in etcd with key %q", pod.ObjectMeta.Name, etcdKey)
	slog.Info("Pod %q created successfully in etcd with key %q", pod.ObjectMeta.Name, etcdKey)
}

// DeletePod handles the deletion of a Pod.
func DeletePod(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming request
	var pod v1.Pod
	if err := ParseJSON(r, &pod); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if pod.Name == "" || pod.Namespace == "" {
		http.Error(w, "Missing required fields: name or namespace", http.StatusBadRequest)
		return
	}

	// Construct the etcd key for the Pod
	etcdKey := fmt.Sprintf("/registry/pods/%s/%s", pod.Namespace, pod.Name)

	// Delete the Pod from etcd
	// ctx, cancel := context.WithTimeout(context.Background(), etcdTimeout)
	// defer cancel()

	// _, err := etcdClient.Delete(ctx, etcdKey)
	// if err != nil {
		// http.Error(w, fmt.Sprintf("Failed to delete Pod from etcd: %v", err), http.StatusInternalServerError)
		// return
	// }
	slog.Info(fmt.Sprintf("Pod %q in namespace %q deleted successfully using key %s", pod.Name, pod.Namespace, etcdKey))
	// Respond with success
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Pod %q in namespace %q deleted successfully", pod.Name, pod.Namespace)
}

func Handle(w http.ResponseWriter, r *http.Request) {
	slog.Info("Received request")
	switch r.URL.Path {
	case "/create-pod":
		CreatePod(w, r)
	case "/delete-pod":
		DeletePod(w, r)
	default:
		http.Error(w, "404 Not Found", http.StatusNotFound)
	}

	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Received request")
	fmt.Printf("%q\n", dump)
	fmt.Fprintf(w, "%q", dump)
}
