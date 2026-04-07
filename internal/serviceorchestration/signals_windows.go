package serviceorchestration

import (
	"os"
)

var shutdownSignals = []os.Signal{os.Interrupt}
