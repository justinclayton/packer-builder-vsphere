package vsphere

import (
	"fmt"
	"log"
	"time"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

// StepGetVmSshInfo represents a Packer build step that retrieves
//  information about the newly created and powered on VM that
//  is needed by any Packer provisioners. Currently this just means IP.
type StepGetVmSshInfo struct{}

// Run executes the Packer build step that waits for
//  the IP to become available.
func (s *StepGetVmSshInfo) Run(state multistep.StateBag) multistep.StepAction {
	log.Println("Begin GetVmSshInfo Step")

	config := state.Get("config").(*Config)
	ui := state.Get("ui").(packer.Ui)
	vim := state.Get("vim").(*VimClient)
	newVmId := state.Get("new_vm_id").(string)

	resultCh := make(chan string, 1)
	errCh := make(chan error, 1)
	go vim.waitForIp(resultCh, errCh, newVmId)

	ui.Say("Waiting for IP to become available...")
	select {
	case ip := <-resultCh:
		// Things succeeded, store VM info for later tasks
		ui.Message(fmt.Sprintf("VM is up with IP '%s'", ip))
		state.Put("new_vm_ip", ip)
	case err := <-errCh:
		err = fmt.Errorf("Error retrieving IP: '%s'", err.Error())
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	case <-time.After(config.stateTimeout):
		err := fmt.Errorf("Error retrieving IP: Exceeded configured timeout value '%s'", config.RawStateTimeout)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	log.Println("End GetVmSshInfo Step")
	return multistep.ActionContinue
}

func (s *StepGetVmSshInfo) Cleanup(state multistep.StateBag) {
}
