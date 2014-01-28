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
	vim := state.Get("vim").(*VimClient)
	config := state.Get("config").(*Config)
	ui := state.Get("ui").(packer.Ui)

	ui.Say("Locating source VM...")
	sourceVmId, err := vim.FindByInventoryPath(config.SourceVmPath)
	log.Printf("FindByInventoryPath returned value '%s' for sourceVmId", sourceVmId)
	if err != nil {
		err = fmt.Errorf("Error locating source VM: %s", err.Error())
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Message("Located source VM.")
	state.Put("source_vm_id", sourceVmId)
	log.Println("End GetSourceVmInfo Step")
	return multistep.ActionContinue
}

func (s *StepGetSourceVmInfo) Cleanup(state multistep.StateBag) {
}
