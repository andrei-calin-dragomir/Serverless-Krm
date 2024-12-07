package domain

import (
	"sort"

	corev1 "k8s.io/api/core/v1"
)

// Struct to hold the pod and filtered nodes (before scoring)
type FilteredNodes struct {
	Pod           corev1.Pod    `json:"pod"`
	FilteredNodes []corev1.Node `json:"filtered_nodes"`
}

// Struct to hold the pod and scored nodes (after scoring)
type ScoredNodes struct {
	Pod         corev1.Pod   `json:"pod"`
	ScoredNodes []ScoredNode `json:"scored_nodes"`
}

// ScoredNode is a wrapper around a node with a score
type ScoredNode struct {
	Node  corev1.Node `json:"node"`
	Score float64     `json:"score"` // Example scoring metric
}

// SortByScore sorts the scored nodes by their score in descending order
func (sn *ScoredNodes) SortByScore() {
	sort.Slice(sn.ScoredNodes, func(i, j int) bool {
		return sn.ScoredNodes[i].Score > sn.ScoredNodes[j].Score
	})
}
