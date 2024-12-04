package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"scheduler/internal/domain"
	"scheduler/internal/outbound"
)

const (
	etcdPodsPrefix = "/kubernetes/pods/"
)

type PodBindService struct {
	etcdClient *outbound.ETCDClient
}

func NewPodBindService(etcdClient *outbound.ETCDClient) *PodBindService {
	return &PodBindService{
		etcdClient: etcdClient,
	}
}

func (s *PodBindService) BindPod(scoredNodes domain.ScoredNodes) {
	bindedPod := scoredNodes.Pod
	bindedPod.Spec.NodeName = scoredNodes.ScoredNodes[0].Node.Name

	bindedPodJson, err := json.Marshal(bindedPod)
	if err != nil {
		slog.Error("Error in marshalling bound pod", "podName", bindedPod.Name, "err", err)
		return
	}

	err = s.etcdClient.SaveKey(context.TODO(), etcdPodsPrefix+bindedPod.Namespace+"/"+bindedPod.Name, string(bindedPodJson))
	if err != nil {
		slog.Error("Error saving pod binding to etcd",
			"podName", bindedPod.Name,
			"namespace", bindedPod.Namespace,
			"error", err)
	} else {
		slog.Info("pod Bindeed in etcd")
	}
}
