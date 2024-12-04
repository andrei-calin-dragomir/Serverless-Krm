package service

import (
	"context"
	"encoding/json"
	"errors"
	"scheduler/internal/domain"
	"scheduler/internal/outbound"

	"log/slog"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	etcdNodesPrefix = "/kubernetes/nodes/"
	nodeScoreURL    = "http://localhost:8081/scheduler/nodeScore"
)

type NodeFilterService struct {
	etcdClient *outbound.ETCDClient
}

func NewNodeFilterService(etcdClient *outbound.ETCDClient) *NodeFilterService {
	return &NodeFilterService{
		etcdClient: etcdClient,
	}
}

func (s *NodeFilterService) FilterNodes(pod corev1.Pod) {
	// Fetch nodes from etcd
	slog.Info("Fetching nodes from etcd", slog.String("prefix", etcdNodesPrefix))
	nodesRes, err := s.etcdClient.GetKeysWithPrefix(context.Background(), etcdNodesPrefix)
	if err != nil {
		slog.Error("Error in getting nodes from etcd", slog.String("prefix", etcdNodesPrefix), slog.Any("error", err))
		return
	}

	// Convert nodesRes to a list of corev1.Node
	var nodes []corev1.Node
	for key, nodeData := range nodesRes {
		var node corev1.Node
		err := json.Unmarshal([]byte(nodeData), &node)
		if err != nil {
			slog.Error("Failed to unmarshal node data", slog.String("key", key), slog.Any("error", err))
			continue // Skip this node if unmarshaling fails
		}
		nodes = append(nodes, node)
		slog.Debug("Successfully unmarshaled node", slog.String("nodeName", node.Name), slog.String("key", key))
	}

	slog.Info("Retrieved nodes from etcd", slog.Int("nodeCount", len(nodes)))

	// Filter nodes based on predicates
	filteredNodes, err := filterNodes(pod, nodes)
	if err != nil {
		slog.Error("Error filtering nodes", slog.Any("error", err))
		return
	}

	slog.Info("Filtered nodes based on criteria", slog.Int("filteredNodeCount", len(filteredNodes)))

	// Send filtered nodes
	resp, err := outbound.PostJSON(nodeScoreURL, domain.FilteredNodes{
		Pod:           pod,
		FilteredNodes: filteredNodes,
	})
	if err != nil {
		slog.Error("Error sending Filtered Nodes to nodeScoreFunc", slog.Any("error", err))
		return
	}
	slog.Info("Received response from node score function", slog.Any("response", resp))
}

// filterNodes applies filtering criteria to nodes
func filterNodes(pod corev1.Pod, nodes []corev1.Node) ([]corev1.Node, error) {
	slog.Info("Starting node filtering", slog.Int("totalNodes", len(nodes)))

	var filteredNodes []corev1.Node

	for _, node := range nodes {
		if !isNodeSchedulable(node) {
			slog.Debug("Node is unschedulable", slog.String("nodeName", node.Name))
			continue
		}
		if !checkNodeAffinity(pod, node) {
			slog.Debug("Node does not match affinity rules", slog.String("nodeName", node.Name))
			continue
		}
		if !checkTaintsAndTolerations(pod, node) {
			slog.Debug("Node does not tolerate taints", slog.String("nodeName", node.Name))
			continue
		}
		if !checkResources(pod, node) {
			slog.Debug("Node does not have sufficient resources", slog.String("nodeName", node.Name))
			continue
		}
		filteredNodes = append(filteredNodes, node)
		slog.Debug("Node passed all filters", slog.String("nodeName", node.Name))
	}

	if len(filteredNodes) == 0 {
		slog.Warn("No nodes available after filtering")
		return nil, errors.New("no nodes available after filtering")
	}

	slog.Info("Node filtering complete", slog.Int("filteredNodes", len(filteredNodes)))
	return filteredNodes, nil
}

// isNodeSchedulable checks if a node is marked as schedulable
func isNodeSchedulable(node corev1.Node) bool {
	slog.Debug("Checking if node is schedulable", slog.String("nodeName", node.Name))
	return !node.Spec.Unschedulable
}

// checkNodeAffinity evaluates node affinity rules
func checkNodeAffinity(pod corev1.Pod, node corev1.Node) bool {
	if pod.Spec.Affinity == nil || pod.Spec.Affinity.NodeAffinity == nil {
		slog.Debug("No node affinity specified for the pod", slog.String("nodeName", node.Name))
		return true
	}

	nodeAffinity := pod.Spec.Affinity.NodeAffinity
	if nodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution == nil {
		slog.Debug("No required node affinity rules specified for the pod", slog.String("nodeName", node.Name))
		return true
	}

	terms := nodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms
	for _, term := range terms {
		if matchNodeSelectorTerm(term, node.Labels) {
			return true
		}
	}
	slog.Debug("Node affinity not satisfied", slog.String("nodeName", node.Name))
	return false
}

// matchNodeSelectorTerm checks if node labels match a selector term
func matchNodeSelectorTerm(term corev1.NodeSelectorTerm, labels map[string]string) bool {
	for _, expr := range term.MatchExpressions {
		switch expr.Operator {
		case corev1.NodeSelectorOpIn:
			if !contains(expr.Values, labels[expr.Key]) {
				return false
			}
		case corev1.NodeSelectorOpNotIn:
			if contains(expr.Values, labels[expr.Key]) {
				return false
			}
		case corev1.NodeSelectorOpExists:
			if _, exists := labels[expr.Key]; !exists {
				return false
			}
		case corev1.NodeSelectorOpDoesNotExist:
			if _, exists := labels[expr.Key]; exists {
				return false
			}
		}
	}
	return true
}

// contains checks if a value exists in a slice
func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// checkTaintsAndTolerations ensures the pod tolerates the node's taints
func checkTaintsAndTolerations(pod corev1.Pod, node corev1.Node) bool {
	for _, taint := range node.Spec.Taints {
		if !toleratesTaint(pod.Spec.Tolerations, &taint) {
			slog.Debug("Pod does not tolerate taint", slog.String("nodeName", node.Name), slog.Any("taint", taint))
			return false
		}
	}
	return true
}

// toleratesTaint checks if a pod's tolerations cover a taint
func toleratesTaint(tolerations []corev1.Toleration, taint *corev1.Taint) bool {
	for _, toleration := range tolerations {
		if toleration.Key == taint.Key &&
			(toleration.Effect == "" || toleration.Effect == taint.Effect) &&
			(toleration.Operator == corev1.TolerationOpExists || toleration.Value == taint.Value) {
			return true
		}
	}
	return false
}

// checkResources validates if the node has enough resources for the pod
func checkResources(pod corev1.Pod, node corev1.Node) bool {
	nodeCPU := node.Status.Allocatable[corev1.ResourceCPU]
	nodeMemory := node.Status.Allocatable[corev1.ResourceMemory]

	var podCPU, podMemory resource.Quantity
	for _, container := range pod.Spec.Containers {
		if container.Resources.Requests == nil {
			slog.Debug("Container resource requests are not set", slog.String("containerName", container.Name))
			continue
		}
		podCPU.Add(*container.Resources.Requests.Cpu())
		podMemory.Add(*container.Resources.Requests.Memory())
	}

	result := podCPU.Cmp(nodeCPU) <= 0 && podMemory.Cmp(nodeMemory) <= 0
	if !result {
		slog.Debug("Node does not meet resource requirements", slog.String("nodeName", node.Name),
			slog.String("podCPU", podCPU.String()), slog.String("nodeCPU", nodeCPU.String()),
			slog.String("podMemory", podMemory.String()), slog.String("nodeMemory", nodeMemory.String()))
	}
	return result
}
