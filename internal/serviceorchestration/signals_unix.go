//go:build !windows

package serviceorchestration

import (
	"os"
	"syscall"
)

var shutdownSignals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
