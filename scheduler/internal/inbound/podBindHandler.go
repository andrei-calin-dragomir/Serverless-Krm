package inbound

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"scheduler/internal/domain"
	"scheduler/internal/service"
)

type PodBindHandler struct {
	podBindService *service.PodBindService
}

func NewPodBindHandler(podBindService *service.PodBindService) *PodBindHandler {
	return &PodBindHandler{podBindService: podBindService}
}

func (nsh *PodBindHandler) HandlePodBind(w http.ResponseWriter, r *http.Request) {

	var scoredNodes domain.ScoredNodes
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&scoredNodes); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		slog.Error("Failed to decode filteredNodes JSON", slog.String("error", err.Error()))
		return
	}
	// Close the request body
	defer r.Body.Close()
	// Log received pod details
	slog.Info("Received pod JSON", "scored nodes object : ", fmt.Sprintf("%+v", scoredNodes))

	// Respond immediately with 200 OK and no body
	w.WriteHeader(http.StatusOK)

	// suggests remove nil check
	if scoredNodes.ScoredNodes == nil || len(scoredNodes.ScoredNodes) == 0 {
		slog.Info(" 0 scored Nodes")
	} else {
		slog.Info("scored nodes num" + fmt.Sprint(len(scoredNodes.ScoredNodes)))
		go nsh.podBindService.BindPod(scoredNodes)
	}
}
