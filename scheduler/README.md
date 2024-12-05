# To run
```
go run cmd/xxxxxxxx/main.go
```

* nodeFilter -> localhost:8080/scheduler/watchPod
* nodeScore -> localhost:8081/scheduler/nodeScore
* podBind -> localhost:8082/scheduler/podBind

* scripts contain some scipts to install etcd, populate with dummy nodes


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