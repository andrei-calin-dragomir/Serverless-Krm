package service

import (
	"fmt"
	"log/slog"
	"math"
	"scheduler/internal/domain"
	"scheduler/internal/outbound"

	corev1 "k8s.io/api/core/v1"
)

const (
	podBindURL = "http://localhost:8082/scheduler/podBind"
)

type NodeScoreService struct {
	etcdClient *outbound.ETCDClient
}

func NewNodeScoreService(etcdClient *outbound.ETCDClient) *NodeScoreService {
	return &NodeScoreService{
		etcdClient: etcdClient,
	}
}

func (s *NodeScoreService) ScoreNodes(filteredNodes domain.FilteredNodes) {
	var scoredNodes domain.ScoredNodes
	scoredNodes.Pod = filteredNodes.Pod

	// Score each filtered node
	for _, node := range filteredNodes.FilteredNodes {
		score := calculateNodeScore(filteredNodes.Pod, node)
		scoredNodes.ScoredNodes = append(scoredNodes.ScoredNodes, domain.ScoredNode{
			Node:  node,
			Score: score,
		})
	}

	scoredNodes.SortByScore()

	// Send scored nodes to pod binding function
	if len(scoredNodes.ScoredNodes) > 0 {
		fmt.Printf("%+v", scoredNodes)
		fmt.Println((scoredNodes.ScoredNodes[0].Score))
		fmt.Println((scoredNodes.ScoredNodes[1].Score))
		fmt.Println((scoredNodes.ScoredNodes[2].Score))
		resp, err := outbound.PostJSON(podBindURL, scoredNodes)
		if err != nil {
			slog.Error("Error sending scored nodes to pod binding function",
				slog.Any("error", err),
				slog.String("url", podBindURL),
			)
			return
		}
		slog.Info("Received response from pod binding function",
			slog.Any("response", resp),
		)
	}
}

func calculateNodeScore(pod corev1.Pod, node corev1.Node) float64 {
	// Calculate individual scores
	resourceScore := calculateResourceScore(pod, node)
	utilizationScore := calculateNodeUtilizationScore(node)
	affinityScore := calculateNodeAffinityScore(pod, node)
	stabilityScore := calculateNodeStabilityScore(node)

	// Combine scores equally (average)
	totalScore := (resourceScore + utilizationScore + affinityScore + stabilityScore) / 4.0

	return totalScore
}

func calculateResourceScore(pod corev1.Pod, node corev1.Node) float64 {
	var totalRequestScore float64
	var totalAllocatableScore float64

	// Check CPU requirements
	for _, container := range pod.Spec.Containers {
		requestCPU := container.Resources.Requests.Cpu()
		allocatableCPU := node.Status.Allocatable.Cpu()

		if requestCPU.Cmp(*allocatableCPU) <= 0 {
			totalRequestScore += 1.0
		} else {
			// Penalize if request exceeds allocatable
			totalRequestScore -= 0.5
		}

		// Calculate utilization score
		totalAllocatableScore += float64(allocatableCPU.MilliValue())
	}

	// Normalize the score
	if totalAllocatableScore > 0 {
		return math.Min(1.0, math.Max(0, totalRequestScore/float64(len(pod.Spec.Containers))))
	}
	return 0
}

func calculateNodeUtilizationScore(node corev1.Node) float64 {
	allocatable := node.Status.Allocatable
	used := node.Status.Capacity

	cpuAllocatable := allocatable.Cpu()
	cpuUsed := used.Cpu()

	memAllocatable := allocatable.Memory()
	memUsed := used.Memory()

	cpuUtilization := float64(cpuUsed.MilliValue()) / float64(cpuAllocatable.MilliValue())
	memUtilization := float64(memUsed.Value()) / float64(memAllocatable.Value())

	// Lower utilization is better, so invert the score
	return 1.0 - math.Min(1.0, (cpuUtilization+memUtilization)/2)
}

func calculateNodeAffinityScore(pod corev1.Pod, node corev1.Node) float64 {
	// Check node selector requirements
	if pod.Spec.NodeSelector != nil {
		for key, value := range pod.Spec.NodeSelector {
			if nodeValue, exists := node.Labels[key]; !exists || nodeValue != value {
				return 0
			}
		}
	}

	return 1.0
}

func calculateNodeStabilityScore(node corev1.Node) float64 {
	var stabilityScore float64 = 1.0

	for _, condition := range node.Status.Conditions {
		switch condition.Type {
		case corev1.NodeReady:
			if condition.Status != corev1.ConditionTrue {
				stabilityScore -= 0.5
			}
		case corev1.NodeDiskPressure:
			if condition.Status == corev1.ConditionTrue {
				stabilityScore -= 0.3
			}
		case corev1.NodeMemoryPressure:
			if condition.Status == corev1.ConditionTrue {
				stabilityScore -= 0.2
			}
		}
	}

	return math.Max(0, stabilityScore)
}
