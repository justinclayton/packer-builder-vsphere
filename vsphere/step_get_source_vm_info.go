package vsphere

import (
	"fmt"
	"log"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

type StepGetSourceVmInfo struct{}

func (s *StepGetSourceVmInfo) Run(state multistep.StateBag) multistep.StepAction {
	log.Println("Begin GetSourceVmInfo Step")
	vim := state.Get("vim").(*VimSession)
	config := state.Get("config").(*Config)
	ui := state.Get("ui").(packer.Ui)

	ui.Say("Getting info on source VM...")
	sourceVm, err := vim.FindByInventoryPath(config.SourceVmPath)
	if err != nil {
		err := fmt.Errorf("Error retrieving info on source VM: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Message(fmt.Sprintf("Found VM: %s", sourceVm.Name))
	state.Put("source_vm", &sourceVm)
	log.Println("End GetSourceVmInfo Step")
	return multistep.ActionContinue
}

func (s *StepGetSourceVmInfo) Cleanup(state multistep.StateBag) {
}
