package inbound

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"scheduler/internal/domain"
	"scheduler/internal/service"
)

type NodeScoreHandler struct {
	nodeScoreService *service.NodeScoreService
}

func NewNodeScoreHandler(nodeScoreService *service.NodeScoreService) *NodeScoreHandler {
	return &NodeScoreHandler{nodeScoreService: nodeScoreService}
}

func (nsh *NodeScoreHandler) HandleNodeScore(w http.ResponseWriter, r *http.Request) {

	var filteredNodes domain.FilteredNodes
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&filteredNodes); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		slog.Error("Failed to decode filteredNodes JSON", slog.String("error", err.Error()))
		return
	}
	// Close the request body
	defer r.Body.Close()
	// Log received pod details
	slog.Info("Received pod JSON", "pod object : ", fmt.Sprintf("%+v", filteredNodes))

	// Respond immediately with 200 OK and no body
	w.WriteHeader(http.StatusOK)

	// suggests remove nil check
	if filteredNodes.FilteredNodes == nil || len(filteredNodes.FilteredNodes) == 0 {
		slog.Info(" 0 filtered Nodes")
	} else {
		slog.Info("filtered nodes in" + fmt.Sprint(len(filteredNodes.FilteredNodes)))
		go nsh.nodeScoreService.ScoreNodes(filteredNodes)
	}
}
