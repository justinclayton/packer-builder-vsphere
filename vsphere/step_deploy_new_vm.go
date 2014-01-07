package vsphere

import (
	// "fmt"
	"github.com/mitchellh/multistep"
	// "github.com/mitchellh/packer/packer"
	// "log"
)

type StepDeployNewVm struct{}

func (s *StepDeployNewVm) Run(state multistep.StateBag) multistep.StepAction {
	return multistep.ActionContinue
}

func (s *StepDeployNewVm) Cleanup(state multistep.StateBag) {
}
