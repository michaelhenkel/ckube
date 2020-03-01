// +build !windows

package hyperv

// Fallback implementation

import (
	"log"
)

func HypervStartConsole(vmName string) error {
	log.Fatalf("This function should not be called")
	return nil
}

func HypervRestoreConsole() {
}
