package vsphere

import (
	// "fmt"
	"github.com/mitchellh/multistep"
	// "github.com/mitchellh/packer/packer"
	// "log"
)

type StepGetSourceVmInfo struct{}

func (s *StepGetSourceVmInfo) Run(state multistep.StateBag) multistep.StepAction {
	return multistep.ActionContinue
}

func (s *StepGetSourceVmInfo) Cleanup(state multistep.StateBag) {
}
