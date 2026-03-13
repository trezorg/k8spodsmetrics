package noderesources

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	servicenoderesources "github.com/trezorg/k8spodsmetrics/internal/noderesources"
	"github.com/trezorg/k8spodsmetrics/internal/resources"
)

func TestCompactHeaderRow(t *testing.T) {
	row := compactHeaderRow(resources.Resources{resources.All})
	require.Equal(t, []any{"NAME", "CPU(alloc/used/free)", "CPU(req/lim)", "MEM(alloc/used/free)", "MEM(req/lim)", "STO(alloc/used/free)", "EPH(alloc/used/free)"}, []any(row))
}

func TestCompactNodeRow(t *testing.T) {
	resource := testCompactNodeResource()
	row := compactNodeRow(resource, resources.Resources{resources.CPU, resources.Memory})

	require.Equal(t, "node-a", row[0])
	require.Equal(t, "3900/1200/2700", row[1])
	require.Contains(t, row[2], "2200/")
	require.Contains(t, row[2], "6000")
	require.Equal(t, "8KiB/3KiB/5KiB", row[3])
	require.Contains(t, row[4], "8KiB/")
	require.Contains(t, row[4], "16KiB")
}

func TestPrintCompactToIncludesTotalFooter(t *testing.T) {
	var buf bytes.Buffer
	PrintCompactTo(&buf, servicenoderesources.NodeResourceList{testCompactNodeResource(), testSecondCompactNodeResource()}, resources.Resources{resources.CPU})

	output := buf.String()
	require.Contains(t, output, "CPU(ALLOC/USED/FREE)")
	require.Contains(t, output, "TOTAL")
	require.Contains(t, output, "11800/7300/4500")
	require.Contains(t, output, "7600/")
	require.Contains(t, output, "18000")
}

func TestPrintCompactToRespectsResources(t *testing.T) {
	var buf bytes.Buffer
	PrintCompactTo(&buf, servicenoderesources.NodeResourceList{testCompactNodeResource()}, resources.Resources{resources.Storage})

	output := buf.String()
	require.NotContains(t, output, "CPU(ALLOC/USED/FREE)")
	require.NotContains(t, output, "MEM(ALLOC/USED/FREE)")
	require.Contains(t, output, "STO(ALLOC/USED/FREE)")
	require.Contains(t, output, "EPH(ALLOC/USED/FREE)")
	require.NotContains(t, output, "TOTAL")
}

func TestPrintCompactToDoesNotTruncateNodeName(t *testing.T) {
	resource := testCompactNodeResource()
	resource.Name = "very-long-node-name-that-should-stay-complete-in-compact-view"

	var buf bytes.Buffer
	PrintCompactTo(&buf, servicenoderesources.NodeResourceList{resource}, resources.Resources{resources.CPU})

	output := buf.String()
	require.Contains(t, output, "very-long-node-name-that-should-stay-complete-in-compact-view")
	require.NotContains(t, output, "...")
}

func testCompactNodeResource() servicenoderesources.NodeResource {
	return servicenoderesources.NodeResource{
		Name:                        "node-a",
		CPU:                         4000,
		AllocatableCPU:              3900,
		UsedCPU:                     1200,
		CPURequest:                  2200,
		CPULimit:                    6000,
		FreeCPU:                     2700,
		Memory:                      12 * 1024,
		AllocatableMemory:           8 * 1024,
		UsedMemory:                  3 * 1024,
		MemoryRequest:               8 * 1024,
		MemoryLimit:                 16 * 1024,
		FreeMemory:                  5 * 1024,
		Storage:                     12 * 1024,
		AllocatableStorage:          10 * 1024,
		UsedStorage:                 4 * 1024,
		FreeStorage:                 6 * 1024,
		StorageEphemeral:            24 * 1024,
		AllocatableStorageEphemeral: 20 * 1024,
		UsedStorageEphemeral:        7 * 1024,
		FreeStorageEphemeral:        13 * 1024,
	}
}

func testSecondCompactNodeResource() servicenoderesources.NodeResource {
	return servicenoderesources.NodeResource{
		Name:                        "node-b",
		CPU:                         8000,
		AllocatableCPU:              7900,
		UsedCPU:                     6100,
		CPURequest:                  5400,
		CPULimit:                    12000,
		FreeCPU:                     1800,
		Memory:                      40 * 1024,
		AllocatableMemory:           31 * 1024,
		UsedMemory:                  26 * 1024,
		MemoryRequest:               28 * 1024,
		MemoryLimit:                 42 * 1024,
		FreeMemory:                  5 * 1024,
		Storage:                     600 * 1024,
		AllocatableStorage:          500 * 1024,
		UsedStorage:                 410 * 1024,
		FreeStorage:                 90 * 1024,
		StorageEphemeral:            200 * 1024,
		AllocatableStorageEphemeral: 150 * 1024,
		UsedStorageEphemeral:        60 * 1024,
		FreeStorageEphemeral:        90 * 1024,
	}
}
