package vsphere

import (
	"fmt"
	"log"
)

// Artifact represents a vSphere path to the new template as the result of a Packer build.
type Artifact struct {
	templatePath string
	vim          *VimClient
}

// BuilderId returns the builder Id.
func (*Artifact) BuilderId() string {
	return BuilderId
}

// Destroy destroys the VM template represented by the artifact.
func (a *Artifact) Destroy() error {
	log.Printf("Destroying template: %s", a.templatePath)
	a.vim.DeleteVm(a.templatePath)
	return nil
}

// Files returns the files represented by the artifact.
func (*Artifact) Files() []string {
	return nil
}

// Id returns the path to the template VM in vSphere.
func (a *Artifact) Id() string {
	return a.templatePath
}

// String returns the string representation of the artifact.
func (a *Artifact) String() string {
	return fmt.Sprintf("A VM template was created: %v", a.templatePath)
}
