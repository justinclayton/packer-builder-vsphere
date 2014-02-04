package vsphere

import (
	"fmt"
	"log"
	"time"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

type StepPowerDownVm struct{}

func (s *StepPowerDownVm) Run(state multistep.StateBag) multistep.StepAction {
	log.Println("Begin StepPowerDownVm")

	config := state.Get("config").(*Config)
	ui := state.Get("ui").(packer.Ui)
	vim := state.Get("vim").(*VimClient)
	newVmId := state.Get("new_vm_id").(string)

	ui.Say("Waiting for the new VM to power down...")
	if err := vim.shutdownGuest(newVmId); err != nil {
		err = fmt.Errorf("Error powering down VM: '%s'", err.Error())
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// keep polling VM until shutdown is complete
	errCh := make(chan error, 1)
	go vim.waitUntilVmShutdownComplete(errCh, newVmId)

	select {
	case err := <-errCh:
		if err != nil {
			err := fmt.Errorf("Error waiting for VM to shut down: '%s'", err.Error())
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
	case <-time.After(config.stateTimeout):
		err := fmt.Errorf("Error waiting for VM to shut down: Exceeded configured timeout value '%s'", config.RawStateTimeout)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Message("Power down complete.")

	log.Println("End StepPowerDownVm")
	return multistep.ActionContinue
}

func (s *StepPowerDownVm) Cleanup(state multistep.StateBag) {
}
