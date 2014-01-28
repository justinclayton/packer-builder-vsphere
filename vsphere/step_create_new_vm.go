package vsphere

import (
	"fmt"
	"log"
	"time"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

// StepCreateNewVm represents a Packer build step that clones
//  a new VM from a source VM.
type StepCreateNewVm struct{}

// Run executes the Packer build step that clones a new VM from
// a source VM.
func (s *StepCreateNewVm) Run(state multistep.StateBag) multistep.StepAction {
	log.Println("Begin CreateNewVm Step")

	config := state.Get("config").(*Config)
	ui := state.Get("ui").(packer.Ui)
	vim := state.Get("vim").(*VimClient)
	sourceVmId := state.Get("source_vm_id").(string)

	ui.Say("Cloning new VM...")

	log.Println("Getting Info from Source VM")
	sourceVmInfo, err := vim.getVmBasicInfo(sourceVmId)
	if err != nil {
		err = fmt.Errorf("Error getting additional info for source VM: '%s'", err.Error())
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	destFolder := sourceVmInfo["parent"]
	newVmName := config.VmName
	taskId, err := vim.cloneVmTask(sourceVmId, destFolder, newVmName)
	if err != nil {
		err = fmt.Errorf("Error invoking CloneVM_Task on vSphere: '%s'", err.Error())
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	t := vim.NewTask(taskId)

	resultCh := make(chan string, 1)
	errCh := make(chan error, 1)
	go t.WaitForCompletion(resultCh, errCh)

	select {
	case newVmId := <-resultCh:
		// Things succeeded, store VM info for later tasks
		ui.Message("New VM created.")
		state.Put("new_vm_id", newVmId)
	case err := <-errCh:
		err = fmt.Errorf("Error waiting for task to complete: '%s'", err.Error())
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	case <-time.After(config.stateTimeout):
		err := fmt.Errorf("Error waiting for task to complete: Exceeded configured timeout value '%s'", config.RawStateTimeout)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	log.Println("End CreateNewVm Step")
	return multistep.ActionContinue
}

func (s *StepCreateNewVm) Cleanup(state multistep.StateBag) {
}
