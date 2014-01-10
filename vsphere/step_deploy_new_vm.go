package vsphere

import (
	"github.com/mitchellh/multistep"
	"log"
)

type StepDeployNewVm struct{}

func (s *StepDeployNewVm) Run(state multistep.StateBag) multistep.StepAction {
	log.Println("Begin DeployNewVm Step")

	config := state.Get("config").(*Config)
	sourceVm := state.Get("source_vm").(*Vm)

	// log.Println("About to call DeployVM...")
	newVm := sourceVm.DeployVM(config.VmName)
	// log.Println("DeployVM returned!")

	state.Put("new_vm_ip", newVm.Ip)

	log.Println("End DeployNewVm Step")
	return multistep.ActionContinue
}

func (s *StepDeployNewVm) Cleanup(state multistep.StateBag) {
}
