package vsphere

import (
	// "fmt"
	"github.com/mitchellh/multistep"
	// "github.com/mitchellh/packer/packer"
	// "log"
)

type StepMarkVmAsTemplate struct{}

func (s *StepMarkVmAsTemplate) Run(state multistep.StateBag) multistep.StepAction {

	dummyTemplatePath := "/MyDatacenter/vm/MyTemplatesFolder/MyNewTemplate"

	state.Put("template_path", dummyTemplatePath)
	return multistep.ActionContinue
}

func (s *StepMarkVmAsTemplate) Cleanup(state multistep.StateBag) {
}
