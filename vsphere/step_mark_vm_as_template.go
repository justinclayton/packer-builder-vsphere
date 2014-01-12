package vsphere

import (
	// "fmt"
	"github.com/mitchellh/multistep"
	// "github.com/mitchellh/packer/packer"
	// "log"
)

type StepMarkVmAsTemplate struct{}

func (s *StepMarkVmAsTemplate) Run(state multistep.StateBag) multistep.StepAction {
	newVm := state.Get("new_vm").(*Vm)
	newVm.MarkAsTemplate()
	return multistep.ActionContinue
}

func (s *StepMarkVmAsTemplate) Cleanup(state multistep.StateBag) {
}
