package vsphere

import (
	"fmt"
	"log"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

// StepDeployNewVm represents a Packer build step that clones
//  a new VM from a source VM.
type StepDeployNewVm struct{}

// Run executes the Packer build step that clones a new VM from
// a source VM.
func (s *StepDeployNewVm) Run(state multistep.StateBag) multistep.StepAction {
	log.Println("Begin DeployNewVm Step")

	config := state.Get("config").(*Config)
	sourceVm := state.Get("source_vm").(*Vm)
	ui := state.Get("ui").(packer.Ui)

	ui.Say("Cloning new VM...")
	newVm, err := sourceVm.DeployVM(config.VmName)
	if err != nil {
		err := fmt.Errorf("Error cloning VM: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Say(fmt.Sprintf("New VM created: %s", newVm.Name))
	state.Put("new_vm", &newVm)

	log.Println("End DeployNewVm Step")
	return multistep.ActionContinue
}

func (s *StepDeployNewVm) Cleanup(state multistep.StateBag) {
}
