package vsphere

import (
	"fmt"
	"log"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

type StepMarkVmAsTemplate struct{}

func (s *StepMarkVmAsTemplate) Run(state multistep.StateBag) multistep.StepAction {
	log.Println("Begin MarkVmAsTemplate Step")

	ui := state.Get("ui").(packer.Ui)
	vim := state.Get("vim").(*VimClient)
	newVmId := state.Get("new_vm_id").(string)

	ui.Say("Converting new VM to template...")
	if err := vim.MarkAsTemplate(newVmId); err != nil {
		err = fmt.Errorf("Error converting VM to template: '%s'", err.Error())
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Message("Successfully converted VM to template.")

	log.Println("End MarkVmAsTemplate Step")
	return multistep.ActionContinue
}

func (s *StepMarkVmAsTemplate) Cleanup(state multistep.StateBag) {
}
