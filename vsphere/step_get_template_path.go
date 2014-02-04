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
	config := state.Get("config").(*Config)

	// TODO: Stop making assumptions based on config strings (jclayton)
	templatePath := getAssumedTemplatePath(config.SourceVmPath, config.VmName)

	ui.Message(fmt.Sprintf("New VM template's path is '%s'", templatePath))
	state.Put("template_path", templatePath)

	log.Println("End GetTemplatePath Step")
	return multistep.ActionContinue
}

func (s *StepGetTemplatePath) Cleanup(state multistep.StateBag) {
}
