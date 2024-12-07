# To run
```
go run cmd/scheduler/main.go
```

POST pod object on -> localhost:8080/scheduler/pod

the binded pod is updated in etcd and then /pods endpoint of kubelet api is hit getting its address from status.addresses.internalIP field in selected node object.

KNATIVE service yaml file needs to have correct env variables

-------------------------------
ANY thing below this is dump

# To Cover
* etcd node fetch in memory caching
* no sending binding object to api-server to check for resource version field used in case of multiple schedulers and definetly applicable in our case


# Scheduler Node Filter Predicates
* PodFitsHostPorts
* PodFitsHost
* PodFitsResources
* PodMatchNodeSelector
* NoVolumeZoneConflict
* NoDiskConflict
* MaxCSIVolumeCount
* CheckNodeMemPressure
* CheckNodePIDPressure
* CheckNodeDiskPressure
* CheckNodeCondition
* PodToleratesNodeTaints
* CheckVolumeBinding


# Scheduler Node Score Criterias
* SelectorSpreadPriority
* InterPodAffinityPriority
* LeastRequestedPriority
* MostRequestedPriority
* RequestedToCapacityRatioPriority
* BalancedResourceAllocation
* NodePreferAvoidPodsPriority
* NodeAffinityPriority
* TaintTolerationPriority
* ImageLocalityPriority
* ServiceSpreadingPriority
* EqualPriority
* EvenPodsSpreadPriority