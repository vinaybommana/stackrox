package data

import "time"

// NodeResourceInfo contains telemetry data about the resources belonging to a node in a Kubernetes cluster
type NodeResourceInfo struct {
	MilliCores   int
	MemoryBytes  uint64
	StorageBytes uint64
}

// NodeInfo contains telemetry data about a node in a Kubernetes cluster
type NodeInfo struct {
	ID string

	ProviderType                         string
	TotalResources, AllocatableResources *NodeResourceInfo
	Unschedulable                        bool
	HasTaints                            bool
	AdverseConditions                    []string

	KernelVersion           string
	OSImage                 string
	ContainerRuntimeVersion string
	KubeletVersion          string
	KubeProxyVersion        string
	OperatingSystem         string
	Architecture            string

	Collector  *CollectorInfo
	Compliance *RoxComponentInfo
}

// NamespaceInfo contains telemetry data about a namespace in a Kubernetes cluster
type NamespaceInfo struct {
	ID   string
	Name string `json:",omitempty"`

	NumPods        int
	NumDeployments int

	PodChurn        int
	DeploymentChurn int
}

// OrchestratorInfo contains information about an orchestrator
type OrchestratorInfo struct {
	Orchestrator        string
	OrchestratorVersion string
	CloudProvider       string
}

// SensorInfo contains information about a sensor and the cluster it is monitoring
type SensorInfo struct {
	*RoxComponentInfo

	ClusterID   string
	ClusterName string
	LastCheckIn *time.Time
}

// ClusterInfo contains telemetry data about a Kubernetes cluster
type ClusterInfo struct {
	Sensor       *SensorInfo
	Orchestrator *OrchestratorInfo
	Nodes        []*NodeInfo
	Namespaces   []*NamespaceInfo
}
