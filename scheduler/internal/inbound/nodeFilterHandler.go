package inbound

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"scheduler/internal/service"

	corev1 "k8s.io/api/core/v1"
)

type NodeFilterHandler struct {
	nodeFilterService *service.NodeFilterService
}

func NewNodeFilterHandler(nodeFilterService *service.NodeFilterService) *NodeFilterHandler {
	return &NodeFilterHandler{nodeFilterService: nodeFilterService}
}

func (nfh *NodeFilterHandler) HandleNodeFilter(w http.ResponseWriter, r *http.Request) {

	var pod corev1.Pod
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pod); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		slog.Error("Failed to decode pod JSON", slog.String("error", err.Error()))
		return
	}
	// Close the request body
	defer r.Body.Close()
	// Log received pod details
	slog.Info("Received pod JSON", "pod object : ", fmt.Sprintf("%+v", pod))

	// Respond immediately with 200 OK and no body
	w.WriteHeader(http.StatusOK)

	if pod.Spec.NodeName == "" {
		slog.Info("pod is unassigned")
		go nfh.nodeFilterService.FilterNodes(pod)
	} else {
		slog.Info("pod is assigned", "nodename", pod.Spec.NodeName)
	}

}
