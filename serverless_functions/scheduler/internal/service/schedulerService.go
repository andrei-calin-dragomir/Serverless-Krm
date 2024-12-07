package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"scheduler/internal/domain"
	"scheduler/internal/outbound"

	corev1 "k8s.io/api/core/v1"
)

const (
	kubeletPodsPrefix = "/pods"
	// etcdNodesPrefix = "/kubernetes/nodes/"
	etcdPodsPrefix = "/kubernetes/pods/"
)

type SchedulerService struct {
	etcdClient *outbound.ETCDClient
}

func NewSchedulerService(etcdClient *outbound.ETCDClient) *SchedulerService {
	return &SchedulerService{
		etcdClient: etcdClient,
	}
}

func (ss *SchedulerService) SchedulePod(pod corev1.Pod) error {
	// Fetch nodes from etcd
	slog.Info("Fetching nodes from etcd", slog.String("prefix", etcdNodesPrefix))
	nodesRes, err := ss.etcdClient.GetKeysWithPrefix(context.Background(), etcdNodesPrefix)
	if err != nil {
		slog.Error("Error in getting nodes from etcd", slog.String("prefix", etcdNodesPrefix), slog.Any("error", err))
		return err
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
	filteredNodesList, err := filterNodes(pod, nodes)
	if err != nil {
		slog.Error("Error filtering nodes", slog.Any("error", err))
		return err
	}

	filteredNodes := domain.FilteredNodes{
		Pod:           pod,
		FilteredNodes: filteredNodesList,
	}

	slog.Info("Filtered nodes based on criteria", slog.Int("filteredNodeCount", len(filteredNodesList)))

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

	bindedPod := scoredNodes.Pod
	bindedPod.Spec.NodeName = scoredNodes.ScoredNodes[0].Node.Name

	bindedPodJson, err := json.Marshal(bindedPod)
	if err != nil {
		slog.Error("Error in marshalling bound pod", "podName", bindedPod.Name, "err", err)
		return err
	}

	err = ss.etcdClient.SaveKey(context.TODO(), etcdPodsPrefix+bindedPod.Namespace+"/"+bindedPod.Name, string(bindedPodJson))
	if err != nil {
		slog.Error("Error saving pod binding to etcd",
			"podName", bindedPod.Name,
			"namespace", bindedPod.Namespace,
			"error", err)
		return err
	} else {
		slog.Info("pod Binded in etcd")
	}

	var nodeIP string
	for _, address := range scoredNodes.ScoredNodes[0].Node.Status.Addresses {
		if address.Type == corev1.NodeInternalIP {
			nodeIP = address.Address
			break
		}
	}

	if nodeIP == "" {
		slog.Error("node has no internal IP")
		return errors.New("node has no internal IP")
	}

	kubeletURL := fmt.Sprintf("http://%s:10250", nodeIP)
	resp, err := outbound.PostJSON(kubeletURL+kubeletPodsPrefix, bindedPodJson)
	if err != nil {
		slog.Error("Error sending binded pod to kubelet",
			slog.Any("error", err),
			slog.String("url", kubeletURL),
		)
		return err
	}
	slog.Info("Received response from kubelet",
		slog.Any("response", resp),
	)

	return nil
}
