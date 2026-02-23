package noderesources

const (
	storageUsedPercentAlert      = 95
	storageEphemeralPercentAlert = 95
)

type (
	NodeResource struct {
		Name                        string `json:"name" yaml:"name"`
		CPU                         int64  `json:"cpu" yaml:"cpu"`
		Memory                      int64  `json:"memory" yaml:"memory"`
		UsedCPU                     int64  `json:"used_cpu" yaml:"used_cpu"`
		UsedMemory                  int64  `json:"used_memory" yaml:"used_memory"`
		AllocatableCPU              int64  `json:"allocatable_cpu" yaml:"allocatable_cpu"`
		AllocatableMemory           int64  `json:"allocatable_memory" yaml:"allocatable_memory"`
		CPURequest                  int64  `json:"cpu_request" yaml:"cpu_request"`
		MemoryRequest               int64  `json:"memory_request" yaml:"memory_request"`
		CPULimit                    int64  `json:"cpu_limit" yaml:"cpu_limit"`
		MemoryLimit                 int64  `json:"memory_limit" yaml:"memory_limit"`
		AvailableCPU                int64  `json:"available_cpu" yaml:"available_cpu"`
		AvailableMemory             int64  `json:"available_memory" yaml:"available_memory"`
		FreeCPU                     int64  `json:"free_cpu" yaml:"free_cpu"`
		FreeMemory                  int64  `json:"free_memory" yaml:"free_memory"`
		Storage                     int64  `json:"storage" yaml:"storage"`
		AllocatableStorage          int64  `json:"allocatable_storage" yaml:"allocatable_storage"`
		UsedStorage                 int64  `json:"used_storage" yaml:"used_storage"`
		FreeStorage                 int64  `json:"free_storage" yaml:"free_storage"`
		StorageEphemeral            int64  `json:"storage_ephemeral" yaml:"storage_ephemeral"`
		AllocatableStorageEphemeral int64  `json:"allocatable_storage_ephemeral" yaml:"allocatable_storage_ephemeral"`
		UsedStorageEphemeral        int64  `json:"used_storage_ephemeral" yaml:"used_storage_ephemeral"`
		FreeStorageEphemeral        int64  `json:"free_storage_ephemeral" yaml:"free_storage_ephemeral"`
	}
	NodeResourceList         []NodeResource
	NodeResourceListEnvelope struct {
		Items NodeResourceList `json:"items,omitempty" yaml:"items,omitempty"`
	}
	// NodeResourceListEnvelop is kept as a compatibility alias for previous typoed name.
	NodeResourceListEnvelop = NodeResourceListEnvelope
	nodePredicate           func(n NodeResource) bool
)
