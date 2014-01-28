package vsphere

import (
	"fmt"
	"log"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

type StepGetTemplatePath struct{}

func (s *StepGetTemplatePath) Run(state multistep.StateBag) multistep.StepAction {
	log.Println("Begin GetTemplatePath Step")

	ui := state.Get("ui").(packer.Ui)
	vim := state.Get("vim").(*VimClient)
	newVmId := state.Get("new_vm_id").(string)

	vmPath, err := vim.getVmPath(newVmId)
	if err != nil {
		err = fmt.Errorf("Error getting path for VM: '%s'", err.Error())
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Message(fmt.Sprintf("New VM template's path is '%s'", vmPath))
	state.Put("template_path", vmPath)

	log.Println("End GetTemplatePath Step")
	return multistep.ActionContinue
}

func (s *StepGetTemplatePath) Cleanup(state multistep.StateBag) {
}
