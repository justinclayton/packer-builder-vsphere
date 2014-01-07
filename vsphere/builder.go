package vsphere

import (
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/packer"
	"log"
	// "time"
)

const BuilderId = "justinclayton.vsphere"

type Builder struct {
	config *Config
	runner multistep.Runner
}

// Prepare processes the build configuration parameters.
func (b *Builder) Prepare(raws ...interface{}) ([]string, error) {
	c, warnings, errs := NewConfig(raws...)
	if errs != nil {
		return warnings, errs
	}
	b.config = c

	return warnings, nil
}

// Run executes a vsphere Packer build and returns a packer.Artifact
// representing the path to the newly created template VM.
func (b *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {

	// Login to vSphere.
	vim := NewVimSession(b.config.VsphereUsername, b.config.VspherePassword, b.config.VsphereHostUrl)

	// Set up the state.
	state := new(multistep.BasicStateBag)
	state.Put("config", b.config)
	state.Put("vim", vim)
	state.Put("hook", hook)
	state.Put("ui", ui)

	// Build the steps.
	steps := []multistep.Step{
		new(StepGetSourceVmInfo),
		new(StepDeployNewVm),
		// &common.StepConnectSSH{
		// 	SSHAddress:     sshAddress(),
		// 	SSHConfig:      sshConfig(),
		// 	SSHWaitTimeout: 5 * time.Minute,
		// },
		// new(common.StepProvision),
		new(StepMarkVmAsTemplate),
	}

	// Run the steps.
	if b.config.PackerDebug {
		b.runner = &multistep.DebugRunner{
			Steps:   steps,
			PauseFn: common.MultistepDebugFn(ui),
		}
	} else {
		b.runner = &multistep.BasicRunner{Steps: steps}
	}
	b.runner.Run(state)

	// Report any errors.
	//// XXX: "image_name" is a GCE thing. Figure out what it means for you
	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}
	if _, ok := state.GetOk("image_name"); !ok {
		log.Println("Failed to find image_name in state. Bug?")
		return nil, nil
	}

	artifact := &Artifact{
		templatePath: state.Get("template_path").(string),
		vim:          vim,
	}
	return artifact, nil
}

// Cancel.
func (b *Builder) Cancel() {
	if b.runner != nil {
		log.Println("Cancelling the step runner...")
		b.runner.Cancel()
	}
}