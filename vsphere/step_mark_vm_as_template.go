package vsphere

import (
	"fmt"
	"log"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

type StepMarkVmAsTemplate struct{}

func (s *StepMarkVmAsTemplate) Run(state multistep.StateBag) multistep.StepAction {
	log.Println("Begin StepMarkVmAsTemplate")

	ui := state.Get("ui").(packer.Ui)
	vim := state.Get("vim").(*VimClient)
	newVmId := state.Get("new_vm_id").(string)

	ui.Say("Waiting for the new VM to power down...")
	if err := vim.markAsTemplate(newVmId); err != nil {
		err = fmt.Errorf("Error marking VM as template: '%s'", err.Error())
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Message("Successfully converted VM to template.")

	log.Println("End StepMarkVmAsTemplate")
	return multistep.ActionContinue
}

func (s *StepMarkVmAsTemplate) Cleanup(state multistep.StateBag) {
}
