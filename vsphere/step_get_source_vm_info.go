package vsphere

import (
	"github.com/mitchellh/multistep"
	"log"
	// "github.com/mitchellh/packer/packer"
)

type StepGetSourceVmInfo struct{}

func (s *StepGetSourceVmInfo) Run(state multistep.StateBag) multistep.StepAction {
	log.Println("Begin GetSourceVmInfo Step")
	vim := state.Get("vim").(*VimSession)
	config := state.Get("config").(*Config)

	sourceVm := vim.GetVmTemplate(config.SourceVmPath)
	state.Put("source_vm", &sourceVm)

	log.Println("End GetSourceVmInfo Step")
	return multistep.ActionContinue
}

func (s *StepGetSourceVmInfo) Cleanup(state multistep.StateBag) {
}
